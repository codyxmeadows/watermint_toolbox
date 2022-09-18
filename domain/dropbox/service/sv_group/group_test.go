package sv_group

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_client"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_group"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

// ---- Mock tests for Cache

func TestCachedGroup_Create(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		_, err := sv.Create("test", CompanyManaged())
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestCachedGroup_List(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		_, err := sv.List()
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestCachedGroup_Remove(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		err := sv.Remove("test")
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestCachedGroup_Resolve(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		_, err := sv.Resolve("test")
		if err != ErrorGroupNotFoundForGroupId && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestCachedGroup_ResolveByName(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		_, err := sv.ResolveByName("test")
		if err != ErrorGroupNotFoundForName && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestCachedGroup_Update(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := NewCached(ctx)
		_, err := sv.Update(&mo_group.Group{})
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

/// ----- Mock tests for impl

func TestImplGroup_Create(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		_, err := sv.Create("test", ManagementType("company_managed"))
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestImplGroup_List(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		_, err := sv.List()
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})

}

func TestImplGroup_Remove(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		err := sv.Remove("test")
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestImplGroup_Resolve(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		_, err := sv.Resolve("test")
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestImplGroup_ResolveByName(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		_, err := sv.ResolveByName("test")
		if err != ErrorGroupNotFoundForName && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestImplGroup_Update(t *testing.T) {
	qtr_endtoend.TestWithDbxClient(t, func(ctx dbx_client.Client) {
		sv := New(ctx)
		_, err := sv.Update(&mo_group.Group{})
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}
