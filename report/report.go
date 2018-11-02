package report

import (
	"errors"
	"flag"
	"github.com/watermint/toolbox/report/report_csv"
	"github.com/watermint/toolbox/report/report_json"
	"go.uber.org/zap"
)

type Report interface {
	Open(logger *zap.Logger) error
	Close()
	Report(row interface{}) error
}

type Factory struct {
	logger       *zap.Logger
	report       Report
	ReportHeader bool
	ReportPath   string
	ReportFormat string
}

func (y *Factory) FlagConfig(f *flag.FlagSet) {
	descReportPath := "Output file path of the report (default: STDOUT)"
	f.StringVar(&y.ReportPath, "report-path", "", descReportPath)

	descReportFormat := "Output file for/**/mat (csv|json) (default: json)"
	f.StringVar(&y.ReportFormat, "report-format", "json", descReportFormat)

	descReportHeader := "Report with header (for csv)"
	f.BoolVar(&y.ReportHeader, "report-header", true, descReportHeader)
}

func (y *Factory) Open(logger *zap.Logger) error {
	y.logger = logger
	switch y.ReportFormat {
	case "csv":
		y.report = &report_csv.CsvReport{
			ReportPath:   y.ReportPath,
			ReportHeader: y.ReportHeader,
		}
		return y.report.Open(y.logger)

	case "json":
		y.report = &report_json.JsonReport{
			ReportPath: y.ReportPath,
		}
		return y.report.Open(y.logger)

	default:
		y.logger.Error(
			"unsupported report format",
			zap.String("specified_format", y.ReportFormat),
		)
		return errors.New("unsupported format")
	}
}

func (y *Factory) Report(row interface{}) error {
	if y.report == nil {
		y.logger.Fatal("open report before write")
		return errors.New("report was not opened")
	}

	return y.report.Report(row)
}

func (y *Factory) Close() {
	if y.report == nil {
		y.logger.Debug("Report already closed")
		return
	}
	y.report.Close()
}
