package report_csv

import (
	"encoding/csv"
	"github.com/watermint/toolbox/report/report_column"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type CsvReport struct {
	logger       *zap.Logger
	ReportPath   string
	ReportHeader bool
	files        map[string]*os.File
	writers      map[string]*csv.Writer
	parsers      map[string]report_column.Row
}

func (z *CsvReport) prepare(row interface{}) (f *os.File, w *csv.Writer, p report_column.Row, err error) {
	name := report_column.RowName(row)
	if zp, ok := z.parsers[name]; ok {
		p = zp
	}
	if zw, ok := z.writers[name]; ok {
		w = zw
	}
	if zf, ok := z.files[name]; ok {
		f = zf
	}

	if p != nil && w != nil {
		return
	}

	open := func(name string) (f *os.File, w *csv.Writer, err2 error) {
		if z.ReportPath == "" {
			return nil, csv.NewWriter(os.Stdout), nil
		}
		filePath := filepath.Join(z.ReportPath, name+".csv")
		z.logger.Debug("Opening report file", zap.String("path", filePath))
		if zf, err := os.Create(filePath); err != nil {
			z.logger.Error(
				"unable to create report file, fallback to stdout",
				zap.String("path", filePath),
				zap.Error(err),
			)
			return nil, csv.NewWriter(os.Stdout), nil
		} else {
			return zf, csv.NewWriter(zf), nil
		}
	}

	if f != nil {
		f.Close()
		z.logger.Fatal("File opened but no writer and/or parser available")
	}
	f, w, err = open(name)
	if err != nil {
		return nil, nil, nil, err
	}
	p = report_column.NewRow(row, z.logger)

	z.files[name] = f
	z.writers[name] = w
	z.parsers[name] = p

	if z.ReportHeader {
		w.Write(p.Header())
	}
	return
}

func (z *CsvReport) Open(logger *zap.Logger) error {
	z.logger = logger
	if z.files == nil {
		z.files = make(map[string]*os.File)
	}
	if z.writers == nil {
		z.writers = make(map[string]*csv.Writer)
	}
	if z.parsers == nil {
		z.parsers = make(map[string]report_column.Row)
	}
	return nil
}

func (z *CsvReport) Report(row interface{}) error {
	_, w, p, err := z.prepare(row)
	if err != nil {
		return err
	}
	w.Write(p.Values(row))

	return nil
}

func (z *CsvReport) Close() {
	for _, w := range z.writers {
		w.Flush()
	}
	for _, f := range z.files {
		f.Close()
	}
}
