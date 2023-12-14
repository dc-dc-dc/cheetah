package main

import (
	"context"
	"fmt"

	"github.com/dc-dc-dc/cheetah/market/research"
)

func main() {
	listClient := research.NewStockListResearch()
	res, err := listClient.GetList(context.Background(), research.StockListDowJones)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d\n", len(res))
}
