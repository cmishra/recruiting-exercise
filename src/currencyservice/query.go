package main 

import (
	"currencybackend"
	"time"
)


type CurrencyData struct {
	Rates map[string]float64
	LastUpdate time.Time
	Source Provider 
	CurrencyList []string
}


func NewCurrencyData() CurrencyData {
	ret := CurrencyData{}
	ret.Provider = currencybackend.FixerData{}
	return ret
}


func (f *CurrencyData) CurrencySupported(currency string) bool {
	_, ok := self.Rates[currency]
	return ok
}


func (f *CurrencyData) GetRate(base string, targets []string) (string) {
	ratio := self.Rates[target]/self.Rates[base]
	moneyFormatting := fmt.Sprintf("%.2f", ratio)
	return moneyFormatting
}

func (f *CurrencyData) Update() {
	rates, curtime := self.Provider.PullUpdate()
	f.Rates = rates
	f.LastUpdate = curtime
	
	currencyList := make([]string, 0)
	for k, _ := range rates {
		currencyList = append(currencyList, k)
	}
	self.CurrencyList = currencyList
}


