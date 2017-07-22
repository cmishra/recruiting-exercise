package main 

import (
	"currencybackend"
	"time"
)


type CurrencyServer struct {
	Rates map[string]float64
	LastUpdate time.Time
	Source currencybackend.Provider 
	CurrencyList []string
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


func (f *CurrencyServer) GetRate(base string, target string) (float64) {
	ratio := f.Rates[target]/f.Rates[base]
	return ratio
}

func (f *CurrencyServer) Update() {
	rates, curtime := f.Source.PullUpdate()
	f.Rates = rates
	f.LastUpdate = curtime
	
	currencyList := make([]string, 0)
	for k, _ := range rates {
		currencyList = append(currencyList, k)
	}
	f.CurrencyList = currencyList
}


