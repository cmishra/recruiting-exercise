package main 

import (
	"currencybackend"
	"time"
)


type CurrencyServer struct {
	Rates map[string]float64
	CurrencyUpdateTime time.Time
	Source currencybackend.Provider 
	CurrencyList Set
}


func NewCurrencyServer() CurrencyServer {
	ret := CurrencyServer{}
	ret.Source = &currencybackend.Fixer{}
	ret.Update()
	return ret
}


func (f *CurrencyServer) CurrencySupported(currency string) bool {
	_, ok := f.Rates[currency]
	return ok
}


func (f *CurrencyServer) GetRates(base string, targets Set, requestTime time.Time) (map[string]float64) {
	rates := make(map[string]float64)
	if requestTime.Year() == f.CurrencyUpdateTime.Year() && 
		requestTime.YearDay() == f.CurrencyUpdateTime.YearDay() {
		for k, _ := range targets {
			rates[k] = f.Rates[k]/f.Rates[base]
		}
	} else {
		ratesReturn := f.Source.CustomRequest(base, requestTime)
		for k, _ := range targets {
			rates[k] = ratesReturn[k]
		}
	}
	
	return rates
}

func (f *CurrencyServer) Update() {
	rates, curtime := f.Source.PullUpdate()
	f.Rates = rates
	f.CurrencyUpdateTime = curtime
	
	currencyList := make(Set)
	for k, _ := range rates {
		currencyList[k] = true
	}
	f.CurrencyList = currencyList
}

