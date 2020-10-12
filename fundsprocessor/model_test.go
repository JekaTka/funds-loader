package fundsprocessor

import (
	"testing"
	"time"
)

func Test_loadfund_Parse(t *testing.T) {
	type fields struct {
		ID         string
		CustomerID customerID
		LoadAmount string
		Time       time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "correct data to Parse",
			fields: fields{
				ID: "1",
				CustomerID: customerID("1"),
				LoadAmount: "$1.2",
				Time: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "should error",
			fields: fields{
				ID: "1",
				CustomerID: customerID("1"),
				LoadAmount: "abc",
				Time: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf := &loadfund{
				ID:         tt.fields.ID,
				CustomerID: tt.fields.CustomerID,
				LoadAmount: tt.fields.LoadAmount,
				Time:       tt.fields.Time,
			}
			if err := lf.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
