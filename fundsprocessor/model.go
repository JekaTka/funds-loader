package fundsprocessor

import (
	"strconv"
	"strings"
	"time"
)

type customerID string

type loadfund struct {
	ID string `json:"id"`
	CustomerID customerID `json:"customer_id"`
	LoadAmount string `json:"load_amount"`
	Amount float64
	Time time.Time `json:"time"`
}

func (lf *loadfund) Parse() error {
	n := strings.ReplaceAll(lf.LoadAmount, "$", "")
	f, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return err
	}

	lf.Amount = f
	return nil
}

type processedfund struct {
	ID string `json:"id"`
	CustomerID customerID `json:"customer_id"`
	Accepted bool `json:"accepted"`
}
