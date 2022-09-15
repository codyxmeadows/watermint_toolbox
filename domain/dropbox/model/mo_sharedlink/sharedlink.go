package mo_sharedlink

import (
	"encoding/json"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_sharedfolder_member"
	"github.com/watermint/toolbox/essentials/api/api_parser"
	"github.com/watermint/toolbox/essentials/log/esl"
)

type SharedLink interface {
	SharedLinkId() string
	LinkTag() string
	LinkUrl() string
	LinkName() string
	LinkExpires() string
	LinkVisibility() string
	LinkPathLower() string
	File() (file *File, ok bool)
	Folder() (folder *Folder, ok bool)
	EntryRaw() json.RawMessage
	Metadata() *Metadata
}

func newMetadata(raw json.RawMessage) *Metadata {
	ce := &Metadata{}
	if err := api_parser.ParseModelRaw(ce, raw); err != nil {
		esl.Default().Debug("Unable to parse json", esl.Error(err), esl.ByteString("raw", raw))
		return ce
	}
	ce.Raw = raw
	return ce
}

type Metadata struct {
	Raw        json.RawMessage
	Id         string `path:"id" json:"id"`
	Tag        string `path:"\\.tag" json:"tag"`
	Url        string `path:"url" json:"url"`
	Name       string `path:"name" json:"name"`
	Expires    string `path:"expires" json:"expires"`
	PathLower  string `path:"path_lower" json:"path_lower"`
	Visibility string `path:"link_permissions.resolved_visibility.\\.tag" json:"visibility"`
}

func (z *Metadata) Metadata() *Metadata {
	return z
}

func (z *Metadata) EntryRaw() json.RawMessage {
	return z.Raw
}

func (z *Metadata) LinkTag() string {
	return z.Tag
}

func (z *Metadata) File() (file *File, ok bool) {
	if z.Tag == "file" {
		file := &File{}
		if err := api_parser.ParseModelRaw(file, z.Raw); err != nil {
			return nil, false
		}
		return file, true
	}
	return nil, false
}

func (z *Metadata) Folder() (folder *Folder, ok bool) {
	if z.Tag == "folder" {
		folder := &Folder{}
		if err := api_parser.ParseModelRaw(folder, z.Raw); err != nil {
			return nil, false
		}
		return folder, true
	}
	return nil, false
}

func (z *Metadata) SharedLinkId() string {
	return z.Id
}

func (z *Metadata) LinkUrl() string {
	return z.Url
}

func (z *Metadata) LinkName() string {
	return z.Name
}

func (z *Metadata) LinkExpires() string {
	return z.Expires
}

func (z *Metadata) LinkPathLower() string {
	return z.PathLower
}

func (z *Metadata) LinkVisibility() string {
	return z.Visibility
}

type File struct {
	Raw            json.RawMessage
	Id             string `path:"id"`
	Tag            string `path:"\\.tag"`
	Url            string `path:"url"`
	Name           string `path:"name"`
	ClientModified string `path:"client_modified"`
	ServerModified string `path:"server_modified"`
	Revision       string `path:"rev"`
	Expires        string `path:"expires"`
	PathLower      string `path:"path_lower"`
	Size           int    `path:"size"`
	Visibility     string `path:"link_permissions.resolved_visibility.\\.tag"`
}

func (z *File) Metadata() *Metadata {
	return newMetadata(z.Raw)
}

func (z *File) EntryRaw() json.RawMessage {
	return z.Raw
}

func (z *File) SharedLinkId() string {
	return z.Id
}

func (z *File) LinkTag() string {
	return z.Tag
}

func (z *File) LinkUrl() string {
	return z.Url
}

func (z *File) LinkName() string {
	return z.Name
}

func (z *File) LinkExpires() string {
	return z.Expires
}

func (z *File) LinkPathLower() string {
	return z.PathLower
}

func (z *File) LinkVisibility() string {
	return z.LinkVisibility()
}

func (z *File) File() (file *File, ok bool) {
	return z, true
}

func (z *File) Folder() (folder *Folder, ok bool) {
	return nil, false
}

type Folder struct {
	Raw        json.RawMessage
	Id         string `path:"id"`
	Tag        string `path:"\\.tag"`
	Url        string `path:"url"`
	Name       string `path:"name"`
	Expires    string `path:"expires"`
	PathLower  string `path:"path_lower"`
	Visibility string `path:"link_permissions.resolved_visibility.\\.tag"`
}

func (z *Folder) Metadata() *Metadata {
	return newMetadata(z.Raw)
}

func (z *Folder) EntryRaw() json.RawMessage {
	return z.Raw
}

func (z *Folder) SharedLinkId() string {
	return z.Id
}

func (z *Folder) LinkTag() string {
	return z.Tag
}

func (z *Folder) LinkUrl() string {
	return z.Url
}

func (z *Folder) LinkName() string {
	return z.Name
}

func (z *Folder) LinkExpires() string {
	return z.Expires
}

func (z *Folder) LinkVisibility() string {
	return z.Visibility
}

func (z *Folder) LinkPathLower() string {
	return z.PathLower
}

func (z *Folder) File() (file *File, ok bool) {
	return nil, false
}

func (z *Folder) Folder() (folder *Folder, ok bool) {
	return z, true
}

type SharedLinkMember struct {
	Raw          json.RawMessage
	SharedLinkId string `path:"sharedlink.id" json:"shared_link_id"`
	Tag          string `path:"sharedlink.\\.tag" json:"tag"`
	Url          string `path:"sharedlink.url" json:"url"`
	Name         string `path:"sharedlink.name" json:"name"`
	Expires      string `path:"sharedlink.expires" json:"expires"`
	PathLower    string `path:"sharedlink.path_lower" json:"path_lower"`
	Visibility   string `path:"sharedlink.link_permissions.resolved_visibility.\\.tag" json:"visibility"`
	AccountId    string `path:"member.profile.account_id" json:"account_id"`
	TeamMemberId string `path:"member.profile.team_member_id" json:"team_member_id"`
	Email        string `path:"member.profile.email" json:"email"`
	Status       string `path:"member.profile.status.\\.tag" json:"status"`
	Surname      string `path:"member.profile.name.surname" json:"surname"`
	GivenName    string `path:"member.profile.name.given_name" json:"given_name"`
}

func (z *SharedLinkMember) SharedLink() (link SharedLink) {
	link = &Metadata{}
	if err := api_parser.ParseModelPathRaw(link, z.Raw, "sharedlink"); err != nil {
		esl.Default().Warn("unexpected data format", esl.String("entry", string(z.Raw)), esl.Error(err))
		// return empty
		return link
	}
	return link
}

func (z *SharedLinkMember) Member() (member mo_sharedfolder_member.Member) {
	member = &mo_sharedfolder_member.Metadata{}
	if err := api_parser.ParseModelPathRaw(member, z.Raw, "member"); err != nil {
		esl.Default().Warn("unexpected data format", esl.String("entry", string(z.Raw)), esl.Error(err))
		// return empty
		return member
	}
	return member
}

func NewSharedLinkMember(link SharedLink, member *mo_member.Member) (slm *SharedLinkMember) {
	raws := make(map[string]json.RawMessage)
	raws["sharedlink"] = link.EntryRaw()
	raws["member"] = member.Raw
	raw := api_parser.CombineRaw(raws)

	slm = &SharedLinkMember{}
	if err := api_parser.ParseModelRaw(slm, raw); err != nil {
		esl.Default().Warn("unexpected data format", esl.Error(err))
		// return empty
		return slm
	}
	return slm
}
