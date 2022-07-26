package rp_writer_impl

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/report/rp_column"
	"github.com/watermint/toolbox/infra/report/rp_column_impl"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/report/rp_writer"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"sync"
)

type MsgUIWriter struct {
	ReportTitle app_msg.Message
}

var (
	MUIWriter = app_msg.Apply(&MsgUIWriter{}).(*MsgUIWriter)
)

func newUIWriter(name string, ctl app_control.Control) rp_writer.Writer {
	return &uiWriter{
		name: name,
		ctl:  ctl,
	}
}

type uiWriter struct {
	name           string
	ctl            app_control.Control
	table          app_ui.Table
	colModel       rp_column.Column
	showReportName bool
	index          int
	mutex          sync.Mutex
}

func (z *uiWriter) Name() string {
	return z.name
}

func (z *uiWriter) Row(r interface{}) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	if z.index == 0 {
		if z.showReportName {
			z.ctl.UI().Info(MUIWriter.ReportTitle.With("Title", z.name))
		}
		z.table.HeaderRaw(z.colModel.Header()...)
	}
	z.table.RowRaw(z.colModel.ValueStrings(r)...)
	z.index++
}

func (z *uiWriter) Open(ctl app_control.Control, model interface{}, opts ...rp_model.ReportOpt) error {
	z.ctl = ctl
	z.colModel = rp_column_impl.NewModel(model, opts...)
	z.table = ctl.UI().InfoTable(z.name)
	ro := &rp_model.ReportOpts{}
	for _, o := range opts {
		o(ro)
	}
	z.showReportName = ro.ShowReportTitle

	return nil
}

func (z *uiWriter) Close() {
	z.table.Flush()
}
