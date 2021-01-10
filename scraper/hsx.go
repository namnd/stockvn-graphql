package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/namnd/stockvn-graphql/graph/model"
)

type Trade struct {
	ClosePrice float64 `json:"ClosePrice"`
	Volume     float64 `json:"TotalShare"`
	Date       string  `json:"ReportDate"`
}

type Response struct {
	Rows []Row `json:"rows"`
}
type Row struct {
	Cell []string `json:"cell"`
}

func GetTrades(code string, from time.Time, to time.Time) ([]*model.Trade, error) {
	format := "02.01.2006"
	url := fmt.Sprintf("https://www.hsx.vn/Modules/Rsde/Report/GetTradingInfo?symbol=%s&dateFrom=%s&dateTo=%s", code, from.Format(format), to.Format(format))

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

	var trades []*model.Trade
	for _, trade := range tradesData {
		ds := trade.Date
		ts, _ := strconv.Atoi(ds[6 : len(ds)-2])
		ts = (ts / (1000 * 24 * 60 * 60)) * 24 * 60 * 60
		date := time.Unix(int64(ts), 0)
		trades = append(trades, &model.Trade{
			Code:       code,
			ClosePrice: int(trade.ClosePrice),
			Volume:     int(trade.Volume),
			Date:       date,
		})
	}
	return trades, nil
}

func GetTradeMatchingResults(code string) ([]*model.Trade, error) {
	layout := "02/01/2006"
	url := fmt.Sprintf("https://www.hsx.vn/Modules/Rsde/Report/GetMatchingResult?symbol=%s&_search=false&nd=1610168059419&rows=2147483647&page=1&sidx=id&sord=desc", code)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	var trades []*model.Trade
	for _, row := range responseObject.Rows {
		date, _ := time.Parse(layout, row.Cell[0])
		trades = append(trades, &model.Trade{
			Code:          code,
			Date:          date,
			OpenPrice:     convertPrice(row.Cell[1]) * 10,
			ClosePrice:    convertPrice(row.Cell[2]) * 10,
			HighPrice:     convertPrice(row.Cell[3]) * 10,
			LowPrice:      convertPrice(row.Cell[4]) * 10,
			AvgPrice:      convertPrice(row.Cell[5]) * 10,
			BuyOrder:      convertPrice(row.Cell[6]) / 100,
			BuyVolume:     convertPrice(row.Cell[7]) / 10,
			SellOrder:     convertPrice(row.Cell[8]) / 100,
			SellVolume:    convertPrice(row.Cell[9]) / 10,
			MatchedVolume: convertPrice(row.Cell[10]),
			MatchedValue:  convertPrice(row.Cell[11]) * 10000,
		})
	}
	return trades, nil
}

func convertPrice(s string) int {
	s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, ".", "", -1)
	i, _ := strconv.Atoi(s)
	return i
}

func GetPutThroughResults(code string) ([]*model.Trade, error) {
	layout := "02/01/2006"
	url := fmt.Sprintf("https://www.hsx.vn/Modules/Rsde/Report/GetPutThroughResult?symbol=%s&_search=false&nd=1610229346028&rows=2147483647&page=1&sidx=id&sord=desc", code)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	var trades []*model.Trade
	for _, row := range responseObject.Rows {
		date, _ := time.Parse(layout, row.Cell[0])
		trades = append(trades, &model.Trade{
			Code:             code,
			Date:             date,
			PutThroughOrder:  convertPrice(row.Cell[1]) / 100,
			PutThroughVolume: convertPrice(row.Cell[2]),
			PutThroughValue:  convertPrice(row.Cell[3]) * 10000,
		})
	}
	return trades, nil
}

func GetForeignResults(code string) ([]*model.Trade, error) {
	layout := "02/01/2006"
	url := fmt.Sprintf("https://www.hsx.vn/Modules/Rsde/Report/GetForeignTradingResult?symbol=%s&_search=false&nd=1610236356079&rows=2147483647&page=1&sidx=id&sord=desc", code)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	var trades []*model.Trade
	for _, row := range responseObject.Rows {
		date, _ := time.Parse(layout, row.Cell[0])
		trades = append(trades, &model.Trade{
			Code:                    code,
			Date:                    date,
			ForeignRemainVolume:     convertPrice(row.Cell[2]) / 100,
			ForeignBuyVolume:        convertPrice(row.Cell[3]),
			ForeignBuyValue:         convertPrice(row.Cell[4]) * 10000,
			ForeignSellVolume:       convertPrice(row.Cell[5]),
			ForeignSellValue:        convertPrice(row.Cell[6]) * 10000,
			ForeignPutThroughVolume: convertPrice(row.Cell[7]),
			ForeignPutThroughValue:  convertPrice(row.Cell[8]) * 10000,
		})
	}
	return trades, nil
}
