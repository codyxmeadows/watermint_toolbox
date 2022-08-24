package uc_file_relocation

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_client"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_error"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file_relocation"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type MsgRelocation struct {
	ErrorConflictUseFileSync app_msg.Message
}

var (
	MRelocation = app_msg.Apply(&MsgRelocation{}).(*MsgRelocation)
)

type Relocation interface {
	// options: allow_shared_folder, allow_ownership_transfer, auto_rename
	Copy(from, to mo_path.DropboxPath) (err error)

	// options: allow_shared_folder, allow_ownership_transfer, auto_rename
	Move(from, to mo_path.DropboxPath) (err error)
}

func New(ctx dbx_client.Client) Relocation {
	return &relocationImpl{
		ctx: ctx,
	}
}

type relocationImpl struct {
	ctx dbx_client.Client
}

func (z *relocationImpl) Copy(from, to mo_path.DropboxPath) (err error) {
	return z.relocation(from, to, func(from, to mo_path.DropboxPath) (err error) {
		svc := sv_file_relocation.New(z.ctx)
		_, err = svc.Copy(from, to)
		return err
	})
}

func (z *relocationImpl) Move(from, to mo_path.DropboxPath) (err error) {
	return z.relocation(from, to, func(from, to mo_path.DropboxPath) (err error) {
		svc := sv_file_relocation.New(z.ctx)
		_, err = svc.Move(from, to)
		return err
	})
}

func (z *relocationImpl) relocation(from, to mo_path.DropboxPath,
	reloc func(from, to mo_path.DropboxPath) (err error)) (err error) {
	l := z.ctx.Log().With(esl.String("from", from.Path()), esl.String("to", to.Path()))

	svc := sv_file.NewFiles(z.ctx)

	fromEntry, err := svc.Resolve(from)
	if err != nil {
		l.Debug("Cannot resolve from", esl.Error(err))
		return err
	}
	var fromToTag string
	if to.LogicalPath() == "/" {
		fromToTag = fromEntry.Tag() + "-folder"
		l = l.With(esl.String("fromTag", fromEntry.Tag()), esl.String("toTag", "root"))
	} else {
		toEntry, err := svc.Resolve(to)
		if err != nil {
			es := dbx_error.NewErrors(err)
			if es.Path().IsNotFound() {
				l.Debug("To not found. Do relocate", esl.Error(err))
				return reloc(from, to)
			}
			l.Debug("Invalid path to relocate, or restricted", esl.Error(err), esl.String("summary", es.Summary()))
			return err
		}
		fromToTag = fromEntry.Tag() + "-" + toEntry.Tag()
		l = l.With(esl.String("fromTag", fromEntry.Tag()), esl.String("toTag", toEntry.Tag()))
	}

	switch fromToTag {
	case "file-file":
		l.Debug("Do relocate")
		return reloc(from, to)

	case "file-folder", "folder-folder":
		l.Debug("Do relocate into folder")
		toPath := to.ChildPath(fromEntry.Name())
		relocErr := reloc(from, toPath)
		dbxErr := dbx_error.NewErrors(relocErr)
		switch {
		case dbxErr == nil:
			return nil
		case dbxErr.To().IsConflict():
			z.ctx.UI().Error(MRelocation.ErrorConflictUseFileSync.
				With("From", fromEntry.Path().Path()).
				With("To", toPath.Path()))

			return relocErr
		default:
			return relocErr
		}

	case "folder-file":
		l.Debug("Not a folder")
		return errors.New("not a folder")

	default:
		return errors.New("unsupported file/folder type")
	}
}
