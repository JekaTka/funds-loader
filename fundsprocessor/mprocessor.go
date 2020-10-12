package fundsprocessor

import (
	"time"
)

type isoWeekKey struct {
	Year int
	Week int
}

/* isoWeeks is a map with key in ISOWeek format (2020-48).
The ISO 8601 year and week number in which time occurs.
Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1,
and Dec 29 to Dec 31 might belong to week 1 of year n+1.
 */
type isoWeeks map[isoWeekKey][]*loadfund

type mProcessor map[customerID]isoWeeks

type transaction struct {
	ID string
	CustomerID customerID
}

var transactionIDs = make(map[transaction]struct{}) // empty structs doesn't allocate memory

/*
limits:
- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount
 */
func (m mProcessor) Process(lf *loadfund) bool {
	if lf.Amount > 5000 {
		return false
	}
	// check if this transaction ID was successful before
	if _, ok := transactionIDs[transaction{
		ID:         lf.ID,
		CustomerID: lf.CustomerID,
	}]; ok {
		return false
	}

	year, week := lf.Time.ISOWeek()
	key := isoWeekKey{
		Year: year,
		Week: week,
	}

	if _, exists := m[lf.CustomerID]; !exists {
		m.create(key, lf)
		transactionIDs[transaction{lf.ID, lf.CustomerID}] = struct{}{}
		return true
	}

	transactions, exists := m[lf.CustomerID][key] // transactions for current week
	// check if year and week already exists
	if !exists {
		// we need to cleanup previous weeks (overwrite by recreating)
		m.create(key, lf)
		transactionIDs[transaction{lf.ID, lf.CustomerID}] = struct{}{}
		return true
	}

	var (
		totalPerDay float64
		totalPerWeek float64
		countPerDay int
	)
	for _, tr := range transactions {
		totalPerWeek += tr.Amount

		// calculate current day numbers
		if sameDay(lf.Time, tr.Time) {
			totalPerDay += tr.Amount
			countPerDay += 1
		} else {
			totalPerDay, countPerDay = 0, 0
		}
	}

	// checking limits
	if totalPerDay + lf.Amount > 5000 ||
		totalPerWeek + lf.Amount > 20000 ||
		countPerDay >= 3 {
		return false
	}

	transactions = append(transactions, lf)
	m[lf.CustomerID][key] = transactions
	transactionIDs[transaction{lf.ID, lf.CustomerID}] = struct{}{}

	return true
}

// create or initiate map for customerID & isoWeek
func (m mProcessor) create(w isoWeekKey, lf *loadfund) {
	transactions := make([]*loadfund, 0, 3*7) // because of 3 loads per day, and 7 days in a week
	transactions = append(transactions, lf)
	m[lf.CustomerID] = make(isoWeeks)
	m[lf.CustomerID][w] = transactions
}

func sameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
