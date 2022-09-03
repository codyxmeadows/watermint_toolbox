package sv_group_member

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_async"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_client"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_list"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_group"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_group_member"
	"github.com/watermint/toolbox/essentials/api/api_request"
	"github.com/watermint/toolbox/essentials/encoding/es_json"
	"sync"
)

type GroupDirectory interface {
	List(groupId string) (members []*mo_group_member.Member, err error)
}

type GroupMemberReader interface {
	List() (members []*mo_group_member.Member, err error)
}

type GroupMember interface {
	GroupMemberReader
	Add(members ...MemberOpt) (group *mo_group.Group, err error)
	Remove(members ...MemberOpt) (group *mo_group.Group, err error)
}

type MemberOpt func(opt *memberOpts) *memberOpts
type memberOpts struct {
	tag          string
	teamMemberId string
	email        string
}

func ByEmail(email string) MemberOpt {
	return func(opt *memberOpts) *memberOpts {
		opt.tag = "email"
		opt.email = email
		return opt
	}
}
func ByTeamMemberId(teamMemberId string) MemberOpt {
	return func(opt *memberOpts) *memberOpts {
		opt.tag = "team_member_id"
		opt.teamMemberId = teamMemberId
		return opt
	}
}

func New(ctx dbx_client.Client, group *mo_group.Group) GroupMember {
	return &groupMemberImpl{
		ctx:     ctx,
		groupId: group.GroupId,
	}
}

func NewByGroupId(ctx dbx_client.Client, groupId string) GroupMember {
	return &groupMemberImpl{
		ctx:     ctx,
		groupId: groupId,
	}
}

func NewCachedReader(ctx dbx_client.Client, groupId string) GroupMemberReader {
	return &cachedReader{
		group:   NewByGroupId(ctx, groupId),
		mutex:   sync.Mutex{},
		members: nil,
		lastErr: nil,
	}
}

func NewCachedDirectory(ctx dbx_client.Client) GroupDirectory {
	return &cachedDirectory{
		groups: make(map[string]GroupMemberReader),
		mutex:  sync.Mutex{},
		ctx:    ctx,
	}
}

type cachedDirectory struct {
	groups map[string]GroupMemberReader
	mutex  sync.Mutex
	ctx    dbx_client.Client
}

func (z *cachedDirectory) List(groupId string) (members []*mo_group_member.Member, err error) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	if group, ok := z.groups[groupId]; ok {
		return group.List()
	}

	group := NewCachedReader(z.ctx, groupId)
	z.groups[groupId] = group
	return group.List()
}

type cachedReader struct {
	group   GroupMemberReader
	mutex   sync.Mutex
	members []*mo_group_member.Member
	lastErr error
}

func (z *cachedReader) List() (members []*mo_group_member.Member, err error) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	if z.lastErr != nil {
		return nil, z.lastErr
	}

	if z.members != nil {
		return z.members, nil
	}

	members, err = z.group.List()
	if err != nil {
		z.lastErr = err
		return nil, err
	}
	z.members = members
	return z.members, nil
}

type groupMemberImpl struct {
	ctx     dbx_client.Client
	groupId string
}

func (z *groupMemberImpl) List() (members []*mo_group_member.Member, err error) {
	type GS struct {
		Tag     string `json:".tag"`
		GroupId string `json:"group_id"`
	}
	p := struct {
		Group GS  `json:"group"`
		Limit int `json:"limit,omitempty"`
	}{
		Group: GS{
			Tag:     "group_id",
			GroupId: z.groupId,
		},
	}

	members = make([]*mo_group_member.Member, 0)
	res := z.ctx.List("team/groups/members/list", api_request.Param(p)).Call(
		dbx_list.Continue("team/groups/members/list/continue"),
		dbx_list.ResultTag("members"),
		dbx_list.UseHasMore(),
		dbx_list.OnEntry(func(entry es_json.Json) error {
			gm := &mo_group_member.Member{}
			if err := entry.Model(gm); err != nil {
				return err
			}
			members = append(members, gm)
			return nil
		}),
	)
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	return members, nil
}

func (z *groupMemberImpl) Add(members ...MemberOpt) (group *mo_group.Group, err error) {
	type GS struct {
		Tag     string `json:".tag"`
		GroupId string `json:"group_id"`
	}
	type U struct {
		Tag          string `json:".tag"`
		TeamMemberId string `json:"team_member_id,omitempty"`
		Email        string `json:"email,omitempty"`
	}
	type M struct {
		User       U      `json:"user"`
		AccessType string `json:"access_type"`
	}

	mm := make([]*M, 0)
	for _, m := range members {
		mo := &memberOpts{}
		m(mo)
		mm = append(mm, &M{
			AccessType: "member",
			User: U{
				Tag:          mo.tag,
				TeamMemberId: mo.teamMemberId,
				Email:        mo.email,
			},
		})
	}
	p := struct {
		Group         GS   `json:"group"`
		Members       []*M `json:"members"`
		ReturnMembers bool `json:"return_members"`
	}{
		Group: GS{
			Tag:     "group_id",
			GroupId: z.groupId,
		},
		Members:       mm,
		ReturnMembers: false,
	}

	res := z.ctx.Async("team/groups/members/add", api_request.Param(p)).Call(
		dbx_async.Status("team/groups/job_status/get"))
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	group = &mo_group.Group{}
	err = res.Success().Json().Model(group)
	return
}

func (z *groupMemberImpl) Remove(members ...MemberOpt) (group *mo_group.Group, err error) {
	type GS struct {
		Tag     string `json:".tag"`
		GroupId string `json:"group_id"`
	}
	type U struct {
		Tag          string `json:".tag"`
		TeamMemberId string `json:"team_member_id,omitempty"`
		Email        string `json:"email,omitempty"`
	}
	users := make([]*U, 0)
	for _, m := range members {
		mo := &memberOpts{}
		m(mo)
		users = append(users, &U{
			Tag:          mo.tag,
			TeamMemberId: mo.teamMemberId,
			Email:        mo.email,
		})
	}

	p := struct {
		Group         GS   `json:"group"`
		Users         []*U `json:"users"`
		ReturnMembers bool `json:"return_members,omitempty"`
	}{
		Group: GS{
			Tag:     "group_id",
			GroupId: z.groupId,
		},
		Users: users,
	}

	res := z.ctx.Async("team/groups/members/remove", api_request.Param(p)).Call(
		dbx_async.Status("team/groups/job_status/get"))
	if err, fail := res.Failure(); fail {
		return nil, err
	}
	group = &mo_group.Group{}
	err = res.Success().Json().Model(group)
	return
}
