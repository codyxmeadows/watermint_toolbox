package dbx_namespace

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/dbx_api"
	"github.com/watermint/toolbox/dbx_api/dbx_rpc"
)

type Namespace struct {
	NamespaceId   string          `json:"namespace_id" column:"namespace_id"`
	NamespaceType string          `json:"namespace_type" column:"namespace_type"`
	Name          string          `json:"name" column:"name"`
	TeamMemberId  string          `json:"team_member_id,omitempty" column:"team_member_id"`
	Namespace     json.RawMessage `json:"namespace"`
}

func ParseNamespace(n gjson.Result) (namespace *Namespace, annotation dbx_api.ErrorAnnotation, err error) {
	namespaceId := n.Get("namespace_id")
	if !namespaceId.Exists() {
		err = errors.New("required field `namespace_id` not found in the response")
		annotation = dbx_api.ErrorAnnotation{
			ErrorType: dbx_api.ErrorUnexpectedDataType,
			Error:     err,
		}
		return nil, annotation, err
	}

	ns := &Namespace{
		NamespaceId:   namespaceId.String(),
		NamespaceType: n.Get("namespace_type.\\.tag").String(),
		Name:          n.Get("name").String(),
		Namespace:     json.RawMessage(n.Raw),
	}
	return ns, dbx_api.Success, nil
}

type NamespaceList struct {
	OnError func(annotation dbx_api.ErrorAnnotation) bool
	OnEntry func(namespace *Namespace) bool
}

func (w *NamespaceList) List(c *dbx_api.Context) bool {
	list := dbx_rpc.RpcList{
		EndpointList:         "team/namespaces/list",
		EndpointListContinue: "team/namespaces/list/continue",
		UseHasMore:           true,
		ResultTag:            "namespaces",
		OnError:              w.OnError,
		OnEntry: func(namespace gjson.Result) bool {
			n, ea, _ := ParseNamespace(namespace)
			if ea.IsSuccess() {
				if w.OnEntry != nil {
					return w.OnEntry(n)
				} else {
					return true
				}
			}
			if ea.IsFailure() {
				if w.OnError != nil {
					return w.OnError(ea)
				} else {
					return false
				}
			}
			return false
		},
	}
	return list.List(c, nil)
}
