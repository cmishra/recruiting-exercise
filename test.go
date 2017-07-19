package main

import (
	"currencybackend"
	"fmt"
)

func main() {
	source := currencybackend.Fixer{}
	fmt.Println(source.PullUpdate())
}
