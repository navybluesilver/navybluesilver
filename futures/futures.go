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

type Tickers []struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ID           int     `json:"id"`
	CurrentValue float64 `json:"currentValue"`
}

type Futures struct {
	Futures            []Future
	Future1Description string
	Future2Description string
}

type Future struct {
	Underlying string  `json:"underlying"`
	SPOT       float64 `json:"SPOT"`
	Future1    float64 `json:"Future1"`
	Future2    float64 `json:"Future2"`
}

func LoadFutures() (f Futures) {
	t, err := getTickers()
	handleError(err)

	f.Futures = append(f.Futures, load(t, "USD"))
	f.Futures = append(f.Futures, load(t, "ADA"))
	f.Futures = append(f.Futures, load(t, "BCH"))
	f.Futures = append(f.Futures, load(t, "EOS"))
	f.Futures = append(f.Futures, load(t, "ETH"))
	f.Futures = append(f.Futures, load(t, "LTC"))
	f.Futures = append(f.Futures, load(t, "TRX"))
	f.Futures = append(f.Futures, load(t, "XRP"))
	f.Future1Description = future1
	f.Future2Description = future2
	fmt.Printf("%v\n", f)
	return f
}

func load(t Tickers, u string) (f Future) {
	f.Underlying = u
	f.SPOT = loadTicker(t, u)
	f.Future1 = loadFuture(t, u, future1)
	f.Future2 = loadFuture(t, u, future2)
	return f
}

func loadTicker(all Tickers, u string) (value float64) {
	for _, t := range all {
		if t.Name == u {
			return t.CurrentValue
		}
	}
	return 0
}

func loadFuture(all Tickers, u string, future string) (value float64) {
	name := fmt.Sprintf("%s%s", u, future)
	return loadTicker(all, name)
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
