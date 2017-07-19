package main

import (
	"net/http"
	"log"
	"net/url"
	"strings"
)

const (
	port = ":9000"

	InvalidRequestErrorCode = 400
	QueryStringParseErrorMsg = "Invalid querystring"
	QueryStringUnrecognizedParameterMsg = "Query parameter %s not recognized"
	CurrencyUnrecognizedMsg = "Currency %s is not recognized"

)

func queryParameterKnown(param string) bool {
	AcceptedQueryParameters := []string{"base", "target", "timestamp"}
	for _, v := range AcceptedQueryParameters {
		if v == param {
			return true
		}
	}
	return false
}

func currencyRecognized(currency string, w http.ResponseWriter) bool {
	if !source.CurrencySupported(currency) {
		errorString := fmt.Sprintf(CurrencyUnrecognizedMsg, currency)
		log.Println(errorString)
		http.Error(w, errorString, InvalidRequestErrorCode)
		return false
	}
	return true
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	querystring := r.URL.RawQuery
	params, err := url.ParseQuery(querystring)
        if err != nil {
                log.Println(QuerystringParseErrorMsg + ":\t" + err.Error())
		http.Error(w, QueryStringParseErrorMsg, InvalidRequestErrorCode)
		return
        }

	for k, _ := range params {
		recognized := queryParameterKnown(k)
		if !recognized {
			errorString := fmt.Sprintf(QueryStringUnrecognizedParameterMsg, k)
			log.Println(errorString)
			http.Error(w, errorString, InvalidRequestErrorCode)
			return
		}
		// Yes we're upper-casing the date, doesn't matter given RFC 3339 format
		params[k] = strings.UpperCase(params[k])
	}

	currencyFields := []string{"base", "target"}
	for _, field := range currencyFields {
		if currency, ok := params[f]; ok {
			if currencyRecognized(currency, w) {
				return
			}
		}
	}

	targets := make([]string, 0)
	if val, ok := params["target"]; !ok {
		targets = source.CurrencyList
	} else if {
		targets = append(targets, val)
	}

	rate = CurrencyData.GetRate(
		params["base"], 
		targets,
	)

	

}


var source CurrencyData

func main() {
	source = NewCurrencyData()
	source.Update()
	log.Println("Listening to " + port)
	http.HandleFunc("/rates", requestHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
