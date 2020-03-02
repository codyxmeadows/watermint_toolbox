package sharedlink

import (
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/model/mo_sharedlink"
	"github.com/watermint/toolbox/domain/service/sv_sharedlink"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

type Delete struct {
	Peer       rc_conn.ConnUserFile
	Path       mo_path.DropboxPath
	Recursive  bool
	SharedLink rp_model.TransactionReport
}

func (z *Delete) Preset() {
	z.SharedLink.SetModel(&mo_sharedlink.Metadata{}, nil)
}

func (z *Delete) Exec(c app_control.Control) error {
	if err := z.SharedLink.Open(); err != nil {
		return err
	}

	if z.Recursive {
		return z.removeRecursive(c)
	} else {
		return z.removePathAt(c)
	}
}

func (z *Delete) removePathAt(c app_control.Control) error {
	ui := c.UI()
	l := c.Log()
	links, err := sv_sharedlink.New(z.Peer.Context()).ListByPath(z.Path)
	if err != nil {
		return err
	}
	if len(links) < 1 {
		ui.InfoK("recipe.sharedlink.delete.info.no_links_at_the_path", app_msg.P{
			"Path": z.Path.Path(),
		})
		return nil
	}

	var lastErr error
	for _, link := range links {
		ui.InfoK("recipe.sharedlink.delete.progress", app_msg.P{
			"Url":  link.LinkUrl(),
			"Path": link.LinkPathLower(),
		})
		err = sv_sharedlink.New(z.Peer.Context()).Remove(link)
		if err != nil {
			l.Debug("Unable to remove link", zap.Error(err), zap.Any("link", link))
			z.SharedLink.Failure(err, link)
			lastErr = err
		} else {
			z.SharedLink.Success(link, nil)
		}
	}
	return lastErr
}

func (z *Delete) removeRecursive(c app_control.Control) error {
	ui := c.UI()
	l := c.Log().With(zap.String("path", z.Path.Path()))
	links, err := sv_sharedlink.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}
	if len(links) < 1 {
		ui.InfoK("recipe.sharedlink.delete.info.no_links_at_the_path", app_msg.P{
			"Path": z.Path.Path(),
		})
		return nil
	}

	var lastErr error
	for _, link := range links {
		l = l.With(zap.String("linkPath", link.LinkPathLower()))
		rel, err := filepath.Rel(strings.ToLower(z.Path.Path()), link.LinkPathLower())
		if err != nil {
			l.Debug("Skip due to path calc error", zap.Error(err))
			continue
		}
		if strings.HasPrefix(rel, "..") {
			l.Debug("Skip due to non related path")
			continue
		}

		ui.InfoK("recipe.sharedlink.delete.progress", app_msg.P{
			"Url":  link.LinkUrl(),
			"Path": link.LinkPathLower(),
		})
		err = sv_sharedlink.New(z.Peer.Context()).Remove(link)
		if err != nil {
			l.Debug("Unable to remove link", zap.Error(err), zap.Any("link", link))
			z.SharedLink.Failure(err, link)
			lastErr = err
		} else {
			z.SharedLink.Success(link, nil)
		}
	}
	return lastErr
}

func (z *Delete) Test(c app_control.Control) error {
	return qt_errors.ErrorImplementMe
}
