package sv_group_member

import (
	"github.com/watermint/toolbox/domain/service/sv_group"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/quality/infra/qt_api"
	"testing"
)

func TestGroupMemberImpl_List(t *testing.T) {
	qt_api.DoTestBusinessManagement(func(ctx api_context.Context) {
		gsv := sv_group.New(ctx)
		groups, err := gsv.List()
		if err != nil {
			t.Error(err)
			return
		}

		for i, group := range groups {
			if i > 10 {
				break
			}
			msv := New(ctx, group)
			members, err := msv.List()
			if err != nil {
				t.Error(err)
			}
			for _, m := range members {
				if m.TeamMemberId == "" || m.AccessType == "" {
					t.Error("invalid")
				}
			}
		}
	})
}
