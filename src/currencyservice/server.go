package main

import (
	"net/http"
	"log"
	"net/url"
	"strings"
	"fmt"
	"encoding/json"
	"io"
	"time"
)

const (
	port = ":9000"

	InvalidRequestErrorCode = 400
	QueryStringParseErrorMsg = "Invalid querystring"
	QueryStringUnrecognizedParameterMsg = "Query parameter %s not recognized"
	CurrencyUnrecognizedMsg = "Currency %s is not recognized"
	TimestampFormatInvalidMsg = "Timestamp could not be parsed, please submit requests as RFC 3339"
	TimestampFutureMsg = "Timestamp is in the future: %s"
	MultipleBasesSpecifiedMsg = "Multiple base currencies specified"

	TimestampFormat = time.RFC3339

)

type JSONError struct {
    Error string         `json:"error"`
}

type Set map[string]bool

func ErrorResponseJSON(w http.ResponseWriter, errMsg string, ErrorCode int) {
	err := JSONError{Error: errMsg}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(ErrorCode)
        encoder := json.NewEncoder(w)
        encoder.Encode(err)
}

func queryParameterKnown(param string) bool {
	AcceptedQueryParameters := []string{"base", "target", "timestamp"}
	for _, v := range AcceptedQueryParameters {
		if v == param {
			return true
		}
	}
	return false
}

func (server *CurrencyServer) currencyRecognized(currency string, w http.ResponseWriter) bool {
	if !server.CurrencySupported(currency) {
		errorString := fmt.Sprintf(CurrencyUnrecognizedMsg, currency)
		ErrorResponseJSON(w, errorString, InvalidRequestErrorCode)
		return false
	}
	return true
}

type ResponseData struct {
	Base string                  `json:"base"`
	Date string                  `json:"date"`
	Rates map[string]float64     `json:"rates"`
}


func (server *CurrencyServer) RequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	querystring := r.URL.RawQuery
	params, err := url.ParseQuery(querystring)
        if err != nil {
                log.Println(QueryStringParseErrorMsg + ":\t" + err.Error())
		ErrorResponseJSON(w, QueryStringParseErrorMsg, InvalidRequestErrorCode)
		return
        }

	for k, _ := range params {
		recognized := queryParameterKnown(k)
		if !recognized {
			errorString := fmt.Sprintf(QueryStringUnrecognizedParameterMsg, k)
			log.Println(errorString)
			ErrorResponseJSON(w, errorString, InvalidRequestErrorCode)
			return
		}
		// Yes we're upper-casing the date, doesn't matter given RFC 3339 format
		upperCaseValues := make([]string, 0)
		for _, v := range params[k] {
			upperCaseValues = append(upperCaseValues, strings.ToUpper(v))
		}
		params[k] = upperCaseValues
	}

	currencyFields := []string{"base", "target"}
	for _, field := range currencyFields {
		if currencies, ok := params[field]; ok {
			for _, currency := range currencies {
				if !server.currencyRecognized(currency, w) {
					return
				}
			}
		}
	}

	if len(params["base"]) > 1 {
		log.Println(MultipleBasesSpecifiedMsg)
		ErrorResponseJSON(w, MultipleBasesSpecifiedMsg, InvalidRequestErrorCode)
		return
	} else if len(params["base"]) == 0 {
		log.Println("No base currency request, assuming USD")
		params["base"] = []string{"USD"}
	}

	if len(params["timestamp"]) == 0 {
		curtime := time.Now().UTC().Format(TimestampFormat)
		params["timestamp"] = []string{curtime}
		log.Printf("No timestamp requested, assuming %s\n", curtime)
	}

	requestTime, err := time.Parse(TimestampFormat, params["timestamp"][0])
	if err != nil {
		log.Println(TimestampFormatInvalidMsg)
		ErrorResponseJSON(w, TimestampFormatInvalidMsg, InvalidRequestErrorCode)
		return
	} else if requestTime.After(time.Now()) {
		errorMsg := fmt.Sprintf(TimestampFutureMsg, requestTime.Format(TimestampFormat))
		log.Println(TimestampFutureMsg, requestTime.Format(TimestampFormat))
		ErrorResponseJSON(w, errorMsg, InvalidRequestErrorCode)
		return
	}

	var targets Set
	if _, ok := params["target"]; !ok {
		targets = server.CurrencyList
	} else {
		targets = make(Set)
		for _, t := range params["target"] {
			targets[t] = true
		}
	}

	rates := server.GetRates(
		params["base"][0],
		targets,
		requestTime,
	)

	retObj := ResponseData{
		Date: requestTime.Format(TimestampFormat),
		Base: params["base"][0],
		Rates: rates,
	}
        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
        encoder := json.NewEncoder(w)
        encoder.Encode(retObj)
}


func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}





func main() {
	server := NewCurrencyServer()
	log.Println("Updating internal data")
	log.Println("Listening to " + port)
	http.HandleFunc("/rates", server.RequestHandler)
	http.HandleFunc("/health-check", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
