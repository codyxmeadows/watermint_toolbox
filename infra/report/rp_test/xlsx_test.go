package rp_test

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/report/rp_writer_impl"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"path/filepath"
	"testing"
)

func TestXlsx_Rotate(t *testing.T) {
	if qt_endtoend.IsSkipEndToEndTest() {
		t.Skipped()
		return
	}

	type Data struct {
		Index int    `json:"index"`
		Value string `json:"value"`
	}

	qtr_endtoend.TestWithControl(t, func(ctl app_control.Control) {
		name := "xlsx_less_than_threshold"
		x := rp_writer_impl.NewXlsxWriter(name, ctl)
		err := x.Open(ctl, &Data{})
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < rp_writer_impl.XlsxMaxRows; i++ {
			x.Row(&Data{
				Index: i,
				Value: fmt.Sprintf("%04d", i),
			})
		}

		x.Close()

		f, err := xlsx.OpenFile(filepath.Join(ctl.Workspace().Report(), name+".xlsx"))
		if err != nil {
			t.Error(err)
			return
		}

		sheet, ok := f.Sheet[name]
		if !ok {
			t.Error(fmt.Sprintf("Sheet `%s` not found", name))
			return
		}

		for i := 0; i < rp_writer_impl.XlsxMaxRows; i++ {
			row := sheet.Row(i + 1).Cells
			if row[0].Value != fmt.Sprintf("%d", i) {
				t.Error("Invalid row found", i)
			}
		}
	})

	qtr_endtoend.TestWithControl(t, func(ctl app_control.Control) {
		name := "xlsx_equals_to_threshold"
		x := rp_writer_impl.NewXlsxWriter(name, ctl)
		err := x.Open(ctl, &Data{})
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i <= rp_writer_impl.XlsxMaxRows; i++ {
			x.Row(&Data{
				Index: i,
				Value: fmt.Sprintf("%04d", i),
			})
		}

		x.Close()

		{
			f, err := xlsx.OpenFile(filepath.Join(ctl.Workspace().Report(), name+"_0000.xlsx"))
			if err != nil {
				t.Error(err)
				return
			}

			sheet, ok := f.Sheet[name]
			if !ok {
				t.Error(fmt.Sprintf("Sheet `%s` not found", name))
				return
			}

			for i := 0; i <= rp_writer_impl.XlsxMaxRows; i++ {
				row := sheet.Row(i + 1).Cells
				if row[0].Value != fmt.Sprintf("%d", i) {
					t.Error("Invalid row found", i)
				}
			}
		}

		{
			_, err := xlsx.OpenFile(filepath.Join(ctl.Workspace().Report(), name+"_0001.xlsx"))
			if err == nil {
				t.Error("should not exist")
				return
			}
		}
	})

	qtr_endtoend.TestWithControl(t, func(ctl app_control.Control) {
		name := "xlsx_rotate"
		x := rp_writer_impl.NewXlsxWriter(name, ctl)
		err := x.Open(ctl, &Data{})
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i <= rp_writer_impl.XlsxMaxRows*2; i++ {
			x.Row(&Data{
				Index: i,
				Value: fmt.Sprintf("%04d", i),
			})
		}

		x.Close()

		{
			f, err := xlsx.OpenFile(filepath.Join(ctl.Workspace().Report(), name+"_0000.xlsx"))
			if err != nil {
				t.Error(err)
				return
			}

			sheet, ok := f.Sheet[name]
			if !ok {
				t.Error(fmt.Sprintf("Sheet `%s` not found", name))
				return
			}

			for i := 0; i <= rp_writer_impl.XlsxMaxRows; i++ {
				row := sheet.Row(i + 1).Cells
				if row[0].Value != fmt.Sprintf("%d", i) {
					t.Error("Invalid row found", i)
				}
			}
		}

		{
			offset := rp_writer_impl.XlsxMaxRows + 1
			f, err := xlsx.OpenFile(filepath.Join(ctl.Workspace().Report(), name+"_0001.xlsx"))
			if err != nil {
				t.Error(err)
				return
			}

			sheet, ok := f.Sheet[name]
			if !ok {
				t.Error(fmt.Sprintf("Sheet `%s` not found", name))
				return
			}

			for i := 0; i < rp_writer_impl.XlsxMaxRows; i++ {
				row := sheet.Row(i + 1).Cells
				if row[0].Value != fmt.Sprintf("%d", i+offset) {
					t.Error("Invalid row found", i+offset)
				}
			}
		}
	})
}
