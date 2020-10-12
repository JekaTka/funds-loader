package fundsprocessor

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		r io.Reader
	}
	buf := new(bytes.Buffer)
	tests := []struct {
		name string
		args args
		want FundProcessor
	}{
		{
			name: "create FundProcessor",
			args: args{buf},
			want: &fundprocessor{buf, make(mProcessor)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fundprocessor_ProcessTo(t *testing.T) {
	type fields struct {
		r          io.Reader
		mProcessor mProcessor
	}
	b := new(bytes.Buffer)
	b.Write([]byte(`{"id":"15887","customer_id":"528","accepted":true}`))
	b.Write([]byte("\n"))

	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{
			name: "should read from io.Reader and write to io.Writer",
			fields: fields{
				bytes.NewReader([]byte(`{"id":"15887","customer_id":"528","load_amount":"$3318.47","time":"2000-01-01T00:00:00Z"}`)),
				make(mProcessor),
			},
			wantW: b.String(),
			wantErr: false,
		},
		{
			name: "should return error if wrong data",
			fields: fields{
				bytes.NewReader([]byte(`abs`)),
				make(mProcessor),
			},
			wantW: "",
			wantErr: true,
		},
		{
			name: "should return error if wrong data (json)",
			fields: fields{
				bytes.NewReader([]byte(`{"id":"15887","customer_id":"528","load_amount":"$3a3b18.47","time":"2000-01-01T00:00:00Z"}`)),
				make(mProcessor),
			},
			wantW: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &fundprocessor{
				r:          tt.fields.r,
				mProcessor: tt.fields.mProcessor,
			}
			w := &bytes.Buffer{}
			err := fp.ProcessTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ProcessTo() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
