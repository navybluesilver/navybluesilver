package futures

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	future1 = "U18"
	future2 = "Z18"
)


type BitmexLastPriceResponse []struct {
	LastPrice                      float64         `json:"lastPrice"`
}

type Tickers []struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ID           int     `json:"id"`
	CurrentValue float64 `json:"currentValue"`
}

type Products struct {
	Tickers            []Product
	Future1Description string
	Future2Description string
}

type Product struct {
	Underlying string  `json:"underlying"`
	SPOT       float64 `json:"SPOT"`
	Future1    float64 `json:"Future1"`
	Future2    float64 `json:"Future2"`
}

func GetAllProducts() (p Products) {
	t, err := getTickers()
	handleError(err)

	p.Tickers = append(p.Tickers, getProduct(t, "USD"))
	p.Tickers = append(p.Tickers, getProduct(t, "ADA"))
	p.Tickers = append(p.Tickers, getProduct(t, "BCH"))
	p.Tickers = append(p.Tickers, getProduct(t, "EOS"))
	p.Tickers = append(p.Tickers, getProduct(t, "ETH"))
	p.Tickers = append(p.Tickers, getProduct(t, "LTC"))
	p.Tickers = append(p.Tickers, getProduct(t, "TRX"))
	p.Tickers = append(p.Tickers, getProduct(t, "XRP"))
	p.Tickers = append(p.Tickers, getProduct(t, "XAU"))
	p.Future1Description = future1
	p.Future2Description = future2
	fmt.Printf("%v\n", p)
	return p
}

func getProduct(t Tickers, u string) (p Product) {
	p.Underlying = u
	p.SPOT = getTicker(t, u)
	p.Future1 = getFuture(t, u, future1)
	p.Future2 = getFuture(t, u, future2)
	return p
}

func getFuture(all Tickers, u string, future string) (value float64) {
	name := fmt.Sprintf("%s%s", u, future)
	return getTicker(all, name)
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getTickers() (t Tickers, err error) {
	req, err := http.NewRequest("GET", "http://localhost:3000/api/datasources", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}

func getTicker(all Tickers, u string) (value float64) {
	for _, t := range all {
		if t.Name == u {
			return t.CurrentValue
		}
	}
	return 0
}
