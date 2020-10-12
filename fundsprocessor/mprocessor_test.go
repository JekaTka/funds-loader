package fundsprocessor

import (
	"testing"
	"time"
)

func Test_mProcessor_Process(t *testing.T) {
	type args struct {
		lf *loadfund
	}
	toTime := func (t string) time.Time {
		lfTime, _ := time.Parse(time.RFC3339, t)
		return lfTime
	}
	m := make(mProcessor)

	tests := []struct {
		name string
		m    mProcessor
		args args
		want bool
	}{
		{
			name: "should return false if load more then 5000",
			m: m,
			args: args{
				&loadfund{
					ID: "1",
					CustomerID: "1",
					LoadAmount: "$5000.5",
					Amount: 5000.5,
					Time: toTime("2020-01-01T01:01:01Z"),
				},
			},
			want: false,
		},
		{
			name: "should return true if customer id is not in mProcessor",
			m: m,
			args: args{
				&loadfund{
					ID: "2",
					CustomerID: "1",
					LoadAmount: "$25",
					Amount: 25,
					Time: toTime("2020-01-01T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return true. transaction #2 (2020-01-01)",
			m: m,
			args: args{
				&loadfund{
					ID: "3",
					CustomerID: "1",
					LoadAmount: "$35",
					Amount: 35,
					Time: toTime("2020-01-01T01:02:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return false. Because transaction id and customer id already processed (2020-01-01)",
			m: m,
			args: args{
				&loadfund{
					ID: "3",
					CustomerID: "1",
					LoadAmount: "$35",
					Amount: 35,
					Time: toTime("2020-01-01T01:02:08Z"),
				},
			},
			want: false,
		},
		{
			name: "should return true. transaction #3 (2020-01-01)",
			m: m,
			args: args{
				&loadfund{
					ID: "4",
					CustomerID: "1",
					LoadAmount: "$45",
					Amount: 45,
					Time: toTime("2020-01-01T01:03:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return false. Because transaction #4 (2020-01-01)",
			m: m,
			args: args{
				&loadfund{
					ID: "5",
					CustomerID: "1",
					LoadAmount: "$55",
					Amount: 55,
					Time: toTime("2020-01-01T01:04:01Z"),
				},
			},
			want: false,
		},
		{
			name: "should return true if isoWeek is not in mProcessor[customerID] (2020-01-13)",
			m: m,
			args: args{
				&loadfund{
					ID: "6",
					CustomerID: "1",
					LoadAmount: "$25",
					Amount: 25,
					Time: toTime("2020-01-13T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return false. Total limit per day 5000 (2020-01-13)",
			m: m,
			args: args{
				&loadfund{
					ID: "7",
					CustomerID: "1",
					LoadAmount: "$4990",
					Amount: 4990,
					Time: toTime("2020-01-13T01:01:01Z"),
				},
			},
			want: false,
		},
		{
			name: "should return true (2020-01-14)",
			m: m,
			args: args{
				&loadfund{
					ID: "8",
					CustomerID: "1",
					LoadAmount: "$4000",
					Amount: 4000,
					Time: toTime("2020-01-14T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return true (2020-01-15)",
			m: m,
			args: args{
				&loadfund{
					ID: "9",
					CustomerID: "1",
					LoadAmount: "$4000",
					Amount: 4000,
					Time: toTime("2020-01-15T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return true (2020-01-16)",
			m: m,
			args: args{
				&loadfund{
					ID: "10",
					CustomerID: "1",
					LoadAmount: "$4000",
					Amount: 4000,
					Time: toTime("2020-01-16T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return true (2020-01-17)",
			m: m,
			args: args{
				&loadfund{
					ID: "11",
					CustomerID: "1",
					LoadAmount: "$4000",
					Amount: 4000,
					Time: toTime("2020-01-17T01:01:01Z"),
				},
			},
			want: true,
		},
		{
			name: "should return false. Maximum 20,000 per week (2020-01-18)",
			m: m,
			args: args{
				&loadfund{
					ID: "12",
					CustomerID: "1",
					LoadAmount: "$4500",
					Amount: 4500,
					Time: toTime("2020-01-18T01:01:01Z"),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Process(tt.args.lf); got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mProcessor_create(t *testing.T) {
	type args struct {
		w  isoWeekKey
		lf *loadfund
	}
	lfTime, _ := time.Parse(time.RFC3339, "2020-01-01T01:01:01Z")

	tests := []struct {
		name string
		m    mProcessor
		args args
	}{
		{
			name: "create or initiate map",
			m: make(mProcessor),
			args: args{
				w:  isoWeekKey{
					Year: 2020,
					Week: 1,
				},
				lf: &loadfund{
					ID:         "1",
					CustomerID: "2",
					LoadAmount: "$3.5",
					Time:       lfTime,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.create(tt.args.w, tt.args.lf)
			isoWeeksM, exists := tt.m[tt.args.lf.CustomerID]
			if !exists {
				t.Errorf("create() didn't create customerid map")
			}

			if _, exists := isoWeeksM[tt.args.w]; !exists {
				t.Errorf("create() didn't create isoweek map with slice of transactions")
			}
		})
	}
}

func Test_sameDay(t *testing.T) {
	type args struct {
		date1 time.Time
		date2 time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false if not the same day",
			args: args{
				date1: time.Now(),
				date2: time.Now().Add(24 * time.Hour),
			},
			want: false,
		},
		{
			name: "true if the same day",
			args: args{
				date1: time.Now(),
				date2: time.Now(),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameDay(tt.args.date1, tt.args.date2); got != tt.want {
				t.Errorf("sameDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
