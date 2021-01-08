package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/namnd/stockvn-graphql/graph/model"
)

type Trade struct {
	ClosePrice float64 `json:"ClosePrice"`
	Volume     float64 `json:"TotalShare"`
	Date       string  `json:"ReportDate"`
}

func GetTrades(searchParams *model.TradeSearchParams) ([]*model.Trade, error) {
	var trades []*model.Trade

	layout := "2006-02-01"
	format := "01.02.2006"
	from, err := time.Parse(layout, searchParams.From)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse(layout, searchParams.To)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://www.hsx.vn/Modules/Rsde/Report/GetTradingInfo?symbol=%s&dateFrom=%s&dateTo=%s", searchParams.Code, from.Format(format), to.Format(format))

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	var tradesData []Trade
	json.Unmarshal(responseData, &tradesData)

	for _, trade := range tradesData {
		ds := trade.Date
		ts, _ := strconv.Atoi(ds[6 : len(ds)-2])
		trades = append(trades, &model.Trade{
			ClosePrice: int(trade.ClosePrice),
			Volume:     int(trade.Volume),
			Timestamp:  ts,
		})
	}
	return trades, nil
}
