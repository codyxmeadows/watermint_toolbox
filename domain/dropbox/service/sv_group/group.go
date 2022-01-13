package sv_group

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_async"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_list"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_group"
	"github.com/watermint/toolbox/essentials/encoding/es_json"
	"github.com/watermint/toolbox/infra/api/api_request"
	"strings"
)

var (
	ErrorGroupNotFoundForGroupId = errors.New("group not found for the group id")
	ErrorGroupNotFoundForName    = errors.New("group not found for the name")
)

type Group interface {
	Resolve(groupId string) (g *mo_group.Group, err error)
	ResolveByName(groupName string) (g *mo_group.Group, err error)
	List() (g []*mo_group.Group, err error)
	Create(name string, opt ...CreateOpt) (g *mo_group.Group, err error)
	Remove(groupId string) error
	Update(group *mo_group.Group) (g *mo_group.Group, err error)
}

type CreateOpt func(opt *createOpts) *createOpts
type createOpts struct {
	mgmtType string
}

func CompanyManaged() CreateOpt {
	return func(opt *createOpts) *createOpts {
		opt.mgmtType = "company_managed"
		return opt
	}
}
func UserManaged() CreateOpt {
	return func(opt *createOpts) *createOpts {
		opt.mgmtType = "user_managed"
		return opt
	}
}
func ManagementType(mgmtType string) CreateOpt {
	return func(opt *createOpts) *createOpts {
		opt.mgmtType = mgmtType
		return opt
	}
}

func New(ctx dbx_context.Context) Group {
	g := &implGroup{
		ctx: ctx,
	}
	return g
}
func NewCached(ctx dbx_context.Context) Group {
	g := &cachedGroup{
		impl: &implGroup{
			ctx: ctx,
		},
		groups: nil,
	}
	return g
}

type cachedGroup struct {
	impl   Group
	groups []*mo_group.Group
}

func (z *cachedGroup) Resolve(groupId string) (g *mo_group.Group, err error) {
	if z.groups == nil {
		if _, err := z.List(); err != nil {
			return nil, err
		}
	}
	for _, g := range z.groups {
		if g.GroupId == groupId {
			return g, nil
		}
	}
	return nil, ErrorGroupNotFoundForGroupId
}

func (z *cachedGroup) ResolveByName(groupName string) (g *mo_group.Group, err error) {
	if z.groups == nil {
		if _, err := z.List(); err != nil {
			return nil, err
		}
	}
	gn := strings.ToLower(groupName)
	for _, g := range z.groups {
		if strings.ToLower(g.GroupName) == gn {
			return g, nil
		}
	}
	return nil, ErrorGroupNotFoundForName
}

func (z *cachedGroup) List() (g []*mo_group.Group, err error) {
	if z.groups == nil {
		z.groups, err = z.impl.List()
		if err != nil {
			return nil, err
		}
	}
	return z.groups, nil
}

func (z *cachedGroup) Create(name string, opt ...CreateOpt) (g *mo_group.Group, err error) {
	z.groups = nil // invalidate cache
	return z.impl.Create(name, opt...)
}

func (z *cachedGroup) Remove(groupId string) error {
	z.groups = nil // invalidate cache
	return z.impl.Remove(groupId)
}

func (z *cachedGroup) Update(group *mo_group.Group) (g *mo_group.Group, err error) {
	z.groups = nil // invalidate cache
	return z.impl.Update(group)
}

type implGroup struct {
	ctx   dbx_context.Context
	limit int
}

func (z *implGroup) ResolveByName(groupName string) (g *mo_group.Group, err error) {
	groups, err := z.List()
	if err != nil {
		return nil, err
	}
	gn := strings.ToLower(groupName)
	for _, g := range groups {
		if strings.ToLower(g.GroupName) == gn {
			return g, nil
		}
	}
	return nil, ErrorGroupNotFoundForName
}

func (z *implGroup) Create(name string, opt ...CreateOpt) (g *mo_group.Group, err error) {
	co := &createOpts{}
	for _, o := range opt {
		o(co)
	}

	type MT struct {
		Tag string `json:".tag"`
	}
	p := struct {
		GroupName           string `json:"group_name"`
		GroupManagementType MT     `json:"group_management_type"`
	}{
		GroupName: name,
		GroupManagementType: MT{
			Tag: co.mgmtType,
		},
	}
	res := z.ctx.Post("team/groups/create", api_request.Param(p))
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	g = &mo_group.Group{}
	err = res.Success().Json().Model(g)
	return
}

func (z *implGroup) Remove(groupId string) error {
	p := struct {
		Tag     string `json:".tag"`
		GroupId string `json:"group_id"`
	}{
		Tag:     "group_id",
		GroupId: groupId,
	}
	res := z.ctx.Async("team/groups/delete", api_request.Param(p)).Call(
		dbx_async.Status("team/groups/job_status/get"),
	)
	if err, fail := res.Failure(); fail {
		return err
	}
	return nil
}

func (z *implGroup) List() (groups []*mo_group.Group, err error) {
	groups = make([]*mo_group.Group, 0)
	p := struct {
		Limit int `json:"limit,omitempty"`
	}{
		Limit: z.limit,
	}

	res := z.ctx.List("team/groups/list", api_request.Param(p)).Call(
		dbx_list.Continue("team/groups/list/continue"),
		dbx_list.UseHasMore(),
		dbx_list.ResultTag("groups"),
		dbx_list.OnEntry(func(entry es_json.Json) error {
			g := &mo_group.Group{}
			if err := entry.Model(g); err != nil {
				return err
			}
			groups = append(groups, g)
			return nil
		}),
	)
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	return groups, nil
}

func (z *implGroup) Resolve(groupId string) (g *mo_group.Group, err error) {
	p := struct {
		Tag      string   `json:".tag"`
		GroupIds []string `json:"group_ids"`
	}{
		Tag:      "group_ids",
		GroupIds: []string{groupId},
	}
	res := z.ctx.Post("team/groups/get_info", api_request.Param(p))
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	g = &mo_group.Group{}
	err = res.Success().Json().FindModel(es_json.PathArrayFirst, g)
	return
}

func (z *implGroup) Update(group *mo_group.Group) (g *mo_group.Group, err error) {
	type GroupSelector struct {
		Tag     string `json:".tag"`
		GroupId string `json:"group_id"`
	}
	u := struct {
		Group                  *GroupSelector `json:"group"`
		NewGroupName           string         `json:"new_group_name,omitempty"`
		NewGroupExternalId     string         `json:"new_group_external_id"`
		NewGroupManagementType string         `json:"new_group_management_type,omitempty"`
	}{
		Group: &GroupSelector{
			Tag:     "group_id",
			GroupId: group.GroupId,
		},
		NewGroupName:           group.GroupName,
		NewGroupExternalId:     group.GroupExternalId,
		NewGroupManagementType: group.GroupManagementType,
	}
	res := z.ctx.Post("team/groups/update", api_request.Param(u))
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	g = &mo_group.Group{}
	err = res.Success().Json().Model(g)
	return
}
