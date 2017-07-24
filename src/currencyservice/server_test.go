package main 

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"strings"
	"encoding/json"
	"regexp"
)

func CheckCodeForOK(rr *httptest.ResponseRecorder, t *testing.T) {
	CheckCode(rr, t, http.StatusOK)
}

func CheckCode(rr *httptest.ResponseRecorder, t *testing.T, expectedCode int) {
	if rr.Code != expectedCode {
		t.Errorf(
			`Wrong code, expected "%v" got "%v"`,
			rr.Code,
			expectedCode,
		)
	}
}

func BodyCheck(rr *httptest.ResponseRecorder, t *testing.T, expectedBody string) {
	trimmedBody := strings.TrimSpace(rr.Body.String())
	trimmedExpected := strings.TrimSpace(expectedBody)
	if trimmedBody != trimmedExpected {
		t.Errorf(
			`Expected: """%s"""`,
			trimmedExpected,
		)
		t.Errorf(
			`Got:      """%s"""`,
			trimmedBody,
		)
	}
}

func ErrorBodyCheck(rr *httptest.ResponseRecorder, t *testing.T, expectedErrorMsg string) {
	errObj := &JSONError{Error: expectedErrorMsg}
	expectedBody, _ := json.Marshal(errObj)
	BodyCheck(rr, t, string(expectedBody))
}

func RequestCheck(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func ExpectedBody20160429(target []string, base string) (string) {
	var expectedBody string
	if len(target) == 1 && base == "USD" {
		if target[0] == "CAD" {
			expectedBody = 
				`{"base":"USD","date":"2016-04-29T14:34:46Z"` +
				`,"rates":{"CAD":1.2528}}`
		}
		if target[0] == "" {
			expectedBody =
				`{"base":"USD",` +
				`"date":"2016-04-29T14:34:46Z",` +
				`"rates":{"AUD":1.3109,"BGN":1.7152,"BRL":3.4849,"CAD":1.2528,` +
				`"CHF":0.96326,"CNY":6.4845,"CZK":23.711,"DKK":6.5281,"EUR":0.87696,"GBP":0.68425,` +
				`"HKD":7.7581,"HRK":6.5869,"HUF":273.81,"IDR":13185,"ILS":3.7416,` +
				`"INR":66.384,"JPY":107.29,"KRW":1141.4,"MXN":17.151,"MYR":3.9065,` +
				`"NOK":8.0812,"NZD":1.4344,"PHP":46.921,"PLN":3.8556,"RON":3.9262,` +
				`"RUB":64.219,"SEK":8.0408,"SGD":1.3427,"THB":34.92,"TRY":2.8005,` +
				`"USD":1,"ZAR":14.169}}`
			
		}
		if target[0] == "USD"{
			expectedBody = 
                                `{"base":"USD",` +
                                `"date":"2016-04-29T14:34:46Z",` +
				`"rates":{"USD":1}}`
		}
	} else if len(target) == 2 && base == "USD" {
		if target[0] == "CAD" && target[1] == "INR" {
                        expectedBody =
                                `{"base":"USD","date":"2016-04-29T14:34:46Z",` +
                                `"rates":{"CAD":1.2528,` +
				`"INR":66.384}}`
		}
	}
	return expectedBody
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	RequestCheck(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(rr, req)
	CheckCodeForOK(rr, t)

	expectedBody := `{"alive": true}`
	BodyCheck(rr, t, expectedBody)
}

func TestWithBaseOneTarget(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?base=USD&target=CAD&timestamp=2016-04-29T14:34:46Z", nil)
	RequestCheck(t, err)
	
	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCodeForOK(rr, t)

	expectedBody := ExpectedBody20160429([]string{"CAD"}, "USD")
	BodyCheck(rr, t, expectedBody)
}

func TestWithMultipleBase(t *testing.T) {
        req, err := http.NewRequest("GET", "/rates?base=USD&base=CAD&target=CAD&timestamp=2016-04-29T14:34:46Z", nil)
        RequestCheck(t, err)

        server := NewCurrencyServer()
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(server.RequestHandler)
        handler.ServeHTTP(rr, req)

        CheckCode(rr, t, InvalidRequestErrorCode)
        ErrorBodyCheck(rr, t, MultipleBasesSpecifiedMsg)
}


// Assumes Fixer.io in list of returned currencies
func TestNoTarget(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?base=USD&timestamp=2016-04-29T14:34:46Z", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCodeForOK(rr, t)

        expectedBody := ExpectedBody20160429([]string{""}, "USD")
	BodyCheck(rr, t, expectedBody)
}

func TestNoBase(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?timestamp=2016-04-29T14:34:46Z", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCodeForOK(rr, t)

	expectedBody := ExpectedBody20160429([]string{""}, "USD")
	BodyCheck(rr, t, expectedBody)
}

func TestMultipleTargets(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?base=USD&target=CAD&target=INR&timestamp=2016-04-29T14:34:46Z", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCodeForOK(rr, t)

        expectedBody := ExpectedBody20160429([]string{"CAD", "INR"}, "USD")
	BodyCheck(rr, t, expectedBody)
}

func TestBaseUnrecognized(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?base=ABC", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)
	
	CheckCode(rr, t, InvalidRequestErrorCode)

	expectedBody := fmt.Sprintf(CurrencyUnrecognizedMsg,  "ABC")
	ErrorBodyCheck(rr, t, expectedBody)
}

func TestTargetUnrecognized(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?target=ABC", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCode(rr, t, InvalidRequestErrorCode)

	expectedBody := fmt.Sprintf(CurrencyUnrecognizedMsg, "ABC")
	ErrorBodyCheck(rr, t, expectedBody)
}

func TestMultipleTargetOneUnrecognized(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?target=CAD&target=ABC", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCode(rr, t, InvalidRequestErrorCode)

	expectedBody := fmt.Sprintf(CurrencyUnrecognizedMsg, "ABC")
	ErrorBodyCheck(rr, t, expectedBody)
}



func TestSameDaySameResult(t *testing.T) {
	req1, err1 := http.NewRequest("GET", "/rates?timestamp=2016-04-29T00:00:01Z", nil)
	RequestCheck(t, err1)

	server := NewCurrencyServer()
	rr1 := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr1, req1)

	CheckCodeForOK(rr1, t)
	expectedBody := rr1.Body.String()
	
	req2, err2 := http.NewRequest("GET", "/rates?timestamp=2016-04-29T23:59:59Z", nil)
	RequestCheck(t, err2)

	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	CheckCodeForOK(rr2, t)
	curBody := rr2.Body.String()

	matcher := regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	expectedBody = matcher.ReplaceAllLiteralString(expectedBody, "")
        curBody = matcher.ReplaceAllLiteralString(curBody, "")
	
	if curBody != expectedBody {
		t.Errorf("Same day results not the same")
		t.Errorf(expectedBody)
		t.Errorf(curBody)
	}
}

func TestFutureDay(t *testing.T) {
	req, err := http.NewRequest("GET", "/rates?timestamp=2018-04-29T00:00:01Z", nil)
	RequestCheck(t, err)

	server := NewCurrencyServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RequestHandler)
	handler.ServeHTTP(rr, req)

	CheckCode(rr, t, InvalidRequestErrorCode)

	expectedBody := fmt.Sprintf(TimestampFutureMsg, "2018-04-29T00:00:01Z")
	ErrorBodyCheck(rr, t, expectedBody)
}

