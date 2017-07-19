package currencybackend

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"time"
	"strings"
)



type Provider interface {
	PullUpdate() (map[string]float64, time.Time)
}


type FixerData struct {
	Base string
	Date string
	Rates map[string]float64
}

const (
        FixerEndpoint = "https://api.fixer.io"
        FixerQuerystring = "latest?base="
	FixerBase = "USD"

	FixerBackendError = "Error pulling from Fixer: %s\n"
)

type Fixer struct {}

func (f *Fixer) ErrorCheck(err error) {
	if err != nil {
		log.Printf(FixerBackendError, err)
		panic(err)
	}
}

func (f *Fixer) PullUpdate() (map[string]float64, time.Time) {
	query := FixerEndpoint + "/" + FixerQuerystring + FixerBase
	resp, err := http.Get(query)
	f.ErrorCheck(err)

	// Reading messages at once is typically bad practice, sets you up for out of memory issues 
	var jsonUnparsed []byte
	jsonUnparsed, err = ioutil.ReadAll(resp.Body)
	f.ErrorCheck(err)
	
	ret := FixerData{}
	err = json.Unmarshal(jsonUnparsed, &ret)
	f.ErrorCheck(err)

	upperCaseMap := make(map[string]float64, 0)
	for k, v := range ret.Rates {
		upperCaseMap[strings.ToUpper(k)] = v
	}
	upperCaseMap[FixerBase] = 1.0

	datetime := time.Now().UTC()
	return upperCaseMap, datetime
}

