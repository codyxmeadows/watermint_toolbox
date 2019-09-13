package app_report

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"sync"
)

func NewUI(name string, row interface{}, ctl app_control.Control) (Report, error) {
	parser := NewColumn(row, ctl)
	r := &UI{
		ctl:    ctl,
		table:  ctl.UI().InfoTable(name),
		parser: parser,
	}
	return r, nil
}

type UI struct {
	ctl    app_control.Control
	table  app_ui.Table
	parser Column
	index  int
	mutex  sync.Mutex
}

func (z *UI) Success(input interface{}, result interface{}) {
	z.Row(TransactionRow{
		Status: z.ctl.UI().Text(msgSuccess.Key(), msgSuccess.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *UI) Failure(reason app_msg.Message, input interface{}, result interface{}) {
	z.Row(TransactionRow{
		Status: z.ctl.UI().Text(msgFailure.Key(), msgFailure.Params()...),
		Reason: z.ctl.UI().Text(reason.Key(), reason.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *UI) Skip(reason app_msg.Message, input interface{}, result interface{}) {
	z.Row(TransactionRow{
		Status: z.ctl.UI().Text(msgSkip.Key(), msgFailure.Params()...),
		Reason: z.ctl.UI().Text(reason.Key(), reason.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *UI) Row(row interface{}) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	if z.index == 0 {
		z.table.HeaderRaw(z.parser.Header()...)
	}
	z.table.RowRaw(z.parser.ValuesAsString(row)...)
	z.index++
}

func (z *UI) Flush() {
	z.table.Flush()
}

func (z *UI) Close() {
	z.table.Flush()
}
