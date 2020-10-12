package fundsprocessor

import (
	"bufio"
	"encoding/json"
	"io"
)

type fundprocessor struct {
	r io.Reader
	mProcessor mProcessor
}

func New(r io.Reader) FundProcessor {
	return &fundprocessor{r, make(mProcessor)}
}

func (fp *fundprocessor) ProcessTo(w io.Writer) error {
	scanner := bufio.NewScanner(fp.r)
	for scanner.Scan() {
		lf := new(loadfund)

		if err := json.Unmarshal(scanner.Bytes(), lf); err != nil {
			return err
		}

		if err := lf.Parse(); err != nil {
			return err
		}

		pf := &processedfund{
			ID:         lf.ID,
			CustomerID: lf.CustomerID,
			Accepted:   fp.mProcessor.Process(lf),
		}
		b, err := json.Marshal(pf)
		if err != nil {
			return err
		}

		w.Write(b)
		w.Write([]byte("\n"))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
