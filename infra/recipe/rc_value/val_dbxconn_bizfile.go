package rc_value

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn_impl"
	"github.com/watermint/toolbox/essentials/encoding/es_json"
	"github.com/watermint/toolbox/essentials/go/es_reflect"
	"github.com/watermint/toolbox/infra/api/api_conn"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"reflect"
)

func newValueDbxConnBusinessFile(peerName string) rc_recipe.Value {
	v := &ValueDbxConnBusinessFile{peerName: peerName}
	v.conn = dbx_conn_impl.NewConnBusinessFile(peerName)
	return v
}

type ValueDbxConnBusinessFile struct {
	conn     dbx_conn.ConnBusinessFile
	peerName string
}

func (z *ValueDbxConnBusinessFile) Spec() (typeName string, typeAttr interface{}) {
	return es_reflect.Key(app.Pkg, z.conn), nil
}

func (z *ValueDbxConnBusinessFile) ValueText() string {
	return z.peerName
}

func (z *ValueDbxConnBusinessFile) Accept(t reflect.Type, v0 interface{}, name string) rc_recipe.Value {
	if t.Implements(reflect.TypeOf((*dbx_conn.ConnBusinessFile)(nil)).Elem()) {
		return newValueDbxConnBusinessFile(z.peerName)
	}
	return nil
}

func (z *ValueDbxConnBusinessFile) Bind() interface{} {
	return &z.peerName
}

func (z *ValueDbxConnBusinessFile) Init() (v interface{}) {
	return z.conn
}

func (z *ValueDbxConnBusinessFile) ApplyPreset(v0 interface{}) {
	z.conn = v0.(dbx_conn.ConnBusinessFile)
	z.peerName = z.conn.PeerName()
}

func (z *ValueDbxConnBusinessFile) Apply() (v interface{}) {
	z.conn.SetPeerName(z.peerName)
	return z.conn
}

func (z *ValueDbxConnBusinessFile) Capture(ctl app_control.Control) (v interface{}, err error) {
	return z.peerName, nil
}

func (z *ValueDbxConnBusinessFile) Restore(v es_json.Json, ctl app_control.Control) error {
	if peerName, found := v.String(); found {
		z.conn = dbx_conn_impl.NewConnBusinessFile(peerName)
		z.peerName = peerName
		return nil
	} else {
		return rc_recipe.ErrorValueRestoreFailed
	}
}

func (z *ValueDbxConnBusinessFile) SpinUp(ctl app_control.Control) error {
	return z.conn.Connect(ctl)
}

func (z *ValueDbxConnBusinessFile) SpinDown(ctl app_control.Control) error {
	return nil
}

func (z *ValueDbxConnBusinessFile) Conn() (conn api_conn.Connection, valid bool) {
	return z.conn, true
}

func (z *ValueDbxConnBusinessFile) Debug() interface{} {
	return map[string]string{
		"peerName": z.peerName,
		"scope":    z.conn.ScopeLabel(),
	}
}
