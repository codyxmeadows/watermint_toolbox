package sharedfolder

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_sharedfolder"
	"github.com/watermint/toolbox/domain/service/sv_sharedfolder"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/quality/qt_test"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_test"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/report/rp_spec_impl"
)

type ListVO struct {
	Peer app_conn.ConnUserFile
}

const (
	reportList = "sharedfolder"
)

type List struct {
}

func (z *List) Reports() []rp_spec.ReportSpec {
	return []rp_spec.ReportSpec{
		rp_spec_impl.Spec(reportList, &mo_sharedfolder.SharedFolder{}),
	}
}

func (z *List) Test(c app_control.Control) error {
	lvo := &ListVO{}
	if !app_test.ApplyTestPeers(c, lvo) {
		return qt_test.NotEnoughResource()
	}
	if err := z.Exec(app_kitchen.NewKitchen(c, lvo)); err != nil {
		return err
	}
	return app_test.TestRows(c, "sharedfolder", func(cols map[string]string) error {
		if _, ok := cols["shared_folder_id"]; !ok {
			return errors.New("shared_folder_id is not found")
		}
		return nil
	})
}

func (*List) Requirement() app_vo.ValueObject {
	return &ListVO{}
}

func (z *List) Exec(k app_kitchen.Kitchen) error {
	var vo interface{} = k.Value()
	lvo := vo.(*ListVO)
	conn, err := lvo.Peer.Connect(k.Control())
	if err != nil {
		return err
	}

	k.Log().Debug("Scanning folders")
	folders, err := sv_sharedfolder.New(conn).List()
	if err != nil {
		return err
	}

	// Write report
	rep, err := rp_spec_impl.New(z, k.Control()).Open(reportList)
	if err != nil {
		return err
	}
	defer rep.Close()

	for _, folder := range folders {
		rep.Row(folder)
	}
	return nil
}
