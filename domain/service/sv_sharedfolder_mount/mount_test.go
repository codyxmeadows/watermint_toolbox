package sv_sharedfolder_mount

import (
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/quality/infra/qt_api"
	"testing"
)

func TestMountImpl_List(t *testing.T) {
	qt_api.DoTestTokenFull(func(ctx api_context.Context) {
		svc := New(ctx)
		mounts, err := svc.List()
		if err != nil {
			t.Error(err)
			return
		}

		for _, m := range mounts {
			if m.SharedFolderId == "" || m.Name == "" {
				t.Error("invalid")
			}
		}
	})
}
