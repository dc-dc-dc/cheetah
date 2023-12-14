package research

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

type StockListOption string

var (
	StockListNasdaqExchange StockListOption = "exchange-is-NASDAQ"
	StockListNyseExchange   StockListOption = "exchange-is-NYSE"
	StockListNasdaqIndex    StockListOption = "inIndex-includes-NASDAQ"
	StockListDowJones       StockListOption = "inIndex-includes-DOW30"
	StockListSp500          StockListOption = "inIndex-includes-SP500"
	StockListActiveReits    StockListOption = "industry-contains-REIT"
)

// sector   is       tag
// industry contains tag

type stockListResearch struct {
	client *http.Client
}

type StockList struct {
	Symbol    string
	Name      string
	MarketCap decimal.Decimal
	Volume    decimal.Decimal
	Industry  string
	Sector    string
}

type stockListResponse struct {
	Status int `json:"status"`
	Data   struct {
		Data []struct {
			Num       int             `json:"no"`
			Symbol    string          `json:"s"`
			Name      string          `json:"n"`
			MarketCap decimal.Decimal `json:"marketCap"`
			Volume    decimal.Decimal `json:"volume"`
			Industry  string          `json:"industry"`
			Sector    string          `json:"sector"`
		} `json:"data"`
		Result int `json:"resultsCount"`
	} `json:"data"`
}

func NewStockListResearch() *stockListResearch {
	return &stockListResearch{
		client: http.DefaultClient,
	}
}

func (sl *stockListResearch) getStockList(ctx context.Context, f StockListOption) ([]StockList, error) {
	url := fmt.Sprintf(util.GetEnv("RESEARCH_API_ENDPOINT", "TODO"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code %d", res.StatusCode)
	}
	var data stockListResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	resData := make([]StockList, data.Data.Result)
	for i, raw := range data.Data.Data {
		resData[i] = StockList{
			Symbol:    raw.Symbol,
			Name:      raw.Name,
			MarketCap: raw.MarketCap,
			Volume:    raw.Volume,
			Industry:  raw.Industry,
			Sector:    raw.Sector,
		}
	}
	return resData, nil
}

func GetSectorOption(name string) StockListOption {
	return StockListOption(fmt.Sprintf("sector-is-%s", name))
}

func GetIndustryOption(name string) StockListOption {
	return StockListOption(fmt.Sprintf("industry-contains-%s", name))
}

func (sl *stockListResearch) GetList(ctx context.Context, listOption StockListOption) ([]StockList, error) {
	return sl.getStockList(ctx, listOption)
}
