package fundsprocessor

import "io"

type FundProcessor interface {
	ProcessTo(w io.Writer) error
}
