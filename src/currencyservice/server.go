package main

import (
	"net/http"
	"log"
	"net/url"
	"strings"
	"fmt"
	"encoding/json"
	"io"
)

const (
	port = ":9000"

	InvalidRequestErrorCode = 400
	QueryStringParseErrorMsg = "Invalid querystring"
	QueryStringUnrecognizedParameterMsg = "Query parameter %s not recognized"
	CurrencyUnrecognizedMsg = "Currency %s is not recognized"
	TimestampFutureMsg = "Timestamp is in the future: %s"
	MultipleBasesSpecifiedMsg = "Multiple base currencies specified"

)

type JSONError struct {
    Error string         `json:"error"`
}

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
	

	targets := make([]string, 0)
	if _, ok := params["target"]; !ok {
		targets = server.CurrencyList
	} else {
		targets = params["target"]
	}

	rates := make(map[string]float64)
	for _, val := range targets {
		rates[val] = server.GetRate(
			params["base"][0],
			val,
		)
	}

	retObj := ResponseData{
		Date: server.LastUpdate.Format("2017-01-01"),
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
