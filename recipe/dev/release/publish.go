package release

import (
	"errors"
	"fmt"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_auth"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn_impl"
	"github.com/watermint/toolbox/domain/github/api/gh_conn"
	"github.com/watermint/toolbox/domain/github/api/gh_context"
	"github.com/watermint/toolbox/domain/github/api/gh_context_impl"
	"github.com/watermint/toolbox/domain/github/model/mo_release"
	"github.com/watermint/toolbox/domain/github/service/sv_reference"
	"github.com/watermint/toolbox/domain/github/service/sv_release"
	"github.com/watermint/toolbox/domain/github/service/sv_release_asset"
	"github.com/watermint/toolbox/essentials/file/es_filehash"
	"github.com/watermint/toolbox/essentials/lang"
	"github.com/watermint/toolbox/essentials/log/esl"
	mo_path2 "github.com/watermint/toolbox/essentials/model/mo_path"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_auth_impl"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"github.com/watermint/toolbox/ingredient/release/homebrew"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"github.com/watermint/toolbox/quality/infra/qt_runtime"
	"github.com/watermint/toolbox/recipe/dev/test"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrorBuildIsNotReadyForRelease = errors.New("the build does not satisfy release criteria")
	ErrorNoArtifactFound           = errors.New("no artifact found")
)

const (
	homebrewRepoOwner  = "watermint"
	homebrewRepoName   = "homebrew-toolbox"
	homebrewRepoBranch = "current"
)

type Publish struct {
	rc_recipe.RemarkConsole
	rc_recipe.RemarkSecret
	ArtifactPath mo_path2.FileSystemPath
	Branch       string
	ConnGithub   gh_conn.ConnGithubRepo
	SkipTests    bool

	HeadingReleaseTheme       app_msg.Message
	HeadingChanges            app_msg.Message
	ListSpecChange            app_msg.Message
	HeadingDocument           app_msg.Message
	ListReadme                app_msg.Message
	HeadingBinary             app_msg.Message
	BinaryTableHeaderFilename app_msg.Message
	BinaryTableHeaderSize     app_msg.Message
	BinaryTableHeaderMD5      app_msg.Message
	BinaryTableHeaderSHA256   app_msg.Message
	TagCommitMessage          app_msg.Message
	ReleaseName               app_msg.Message
	ReleaseNameBeta           app_msg.Message
	ReleaseNameDev            app_msg.Message
}

type ArtifactSum struct {
	Filename string
	Size     int64
	MD5      string
	SHA256   string
}

func (z *Publish) Preset() {
	z.Branch = "master"
}

func (z *Publish) artifactAssets(c app_control.Control) (paths []string, sizes map[string]int64, err error) {
	l := c.Log()

	entries, err := ioutil.ReadDir(z.ArtifactPath.Path())
	if err != nil {
		return nil, nil, err
	}
	paths = make([]string, 0)
	sizes = make(map[string]int64)
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "tbx-"+app.Version) || !strings.HasSuffix(e.Name(), ".zip") {
			l.Debug("Ignore non artifact file", esl.Any("file", e))
			continue
		}
		path := filepath.Join(z.ArtifactPath.Path(), e.Name())
		paths = append(paths, path)
		sizes[path] = e.Size()
	}
	return paths, sizes, nil
}

func (z *Publish) verifyArtifacts(c app_control.Control) (a []*ArtifactSum, err error) {
	l := c.Log()
	a = make([]*ArtifactSum, 0)

	assets, assetSize, err := z.artifactAssets(c)

	h := es_filehash.NewHash(l)
	for _, p := range assets {
		sum := &ArtifactSum{
			Filename: filepath.Base(p),
			Size:     assetSize[p],
		}
		sum.MD5, err = h.MD5(p)
		if err != nil {
			l.Debug("Unable to calc MD5", esl.Error(err))
			return nil, err
		}
		sum.SHA256, err = h.SHA256(p)
		if err != nil {
			l.Debug("Unable to calc SHA256", esl.Error(err))
			return nil, err
		}
		a = append(a, sum)
	}
	if len(a) < 1 {
		return nil, ErrorNoArtifactFound
	}
	return a, nil
}

func (z *Publish) releaseNotes(c app_control.Control, sum []*ArtifactSum) (relNote string, err error) {
	l := c.Log()
	baseUrl := "https://github.com/watermint/toolbox/blob/" + app.Version

	md := app_ui.MakeMarkdown(c.WithLang("en").Messages(), func(mui app_ui.UI) {
		mui.Header(z.HeadingReleaseTheme)
		mui.Break()

		mui.Header(z.HeadingChanges)

		for _, la := range lang.Supported {
			mui.Info(z.ListSpecChange.
				With("Link", baseUrl+"/doc/generated"+la.Suffix()+"/changes.md").
				With("Lang", la.Self()),
			)
		}

		mui.Break()
		mui.Header(z.HeadingDocument)

		for _, la := range lang.Supported {
			name := "README" + la.Suffix() + ".md"
			mui.Info(z.ListReadme.
				With("Name", name).
				With("Link", baseUrl+"/"+name).
				With("Lang", la.Self()),
			)
		}

		mui.Break()
		mui.Header(z.HeadingBinary)
		mui.WithTable("Binaries", func(mit app_ui.Table) {
			mit.Header(z.BinaryTableHeaderFilename, z.BinaryTableHeaderSize, z.BinaryTableHeaderMD5, z.BinaryTableHeaderSHA256)

			for _, s := range sum {
				mit.RowRaw(s.Filename, fmt.Sprintf("%d", s.Size), s.MD5, s.SHA256)
			}
		})
	})

	relNotesPath := filepath.Join(c.Workspace().Report(), "release_notes.md")
	err = ioutil.WriteFile(relNotesPath, []byte(md), 0644)
	if err != nil {
		l.Debug("Unable to write release notes", esl.Error(err), esl.String("path", relNotesPath))
		return "", err
	}
	l.Info("Release note created", esl.String("path", relNotesPath))

	return md, nil
}

func (z *Publish) endToEndTest(c app_control.Control) error {
	l := c.Log()
	if c.Feature().IsTest() {
		l.Info("Skip tests")
		return nil
	}

	if c.Feature().IsProduction() {
		l.Info("Prepare resources")
		if !api_auth_impl.IsLegacyCacheAvailable(c, app.PeerEndToEndTest, []string{
			api_auth.DropboxTokenFull,
			api_auth.DropboxTokenBusinessAudit,
			api_auth.DropboxTokenBusinessManagement,
			api_auth.DropboxTokenBusinessFile,
			api_auth.DropboxTokenBusinessInfo,
		}, dbx_auth.NewLegacyApp(c)) {
			return qt_errors.ErrorNotEnoughResource
		}
	}

	l.Info("Ensure end to end resource availability")
	if !dbx_conn_impl.IsEndToEndTokenAllAvailable(c) {
		l.Error("At least one of end to end resource is not available.")
		return qt_errors.ErrorNotEnoughResource
	}

	l.Info("Testing all end to end test")
	err := rc_exec.Exec(c, &test.Recipe{}, func(r rc_recipe.Recipe) {
		m := r.(*test.Recipe)
		m.All = true
	})
	return err
}

func (z *Publish) ghCtx(c app_control.Control) gh_context.Context {
	if c.Feature().IsTest() {
		return gh_context_impl.NewMock("mock", c)
	} else {
		return z.ConnGithub.Context()
	}
}

func (z *Publish) createTag(c app_control.Control) error {
	l := c.Log().With(
		esl.String("owner", app.RepositoryOwner),
		esl.String("repository", app.RepositoryName),
		esl.String("version", app.Version),
		esl.String("hash", app.Hash))
	svt := sv_reference.New(z.ghCtx(c), app.RepositoryOwner, app.RepositoryName)
	l.Debug("Create tag")
	tag, err := svt.Create(
		"refs/tags/"+app.Version,
		app.Hash,
	)
	if err != nil && err != qt_errors.ErrorMock {
		l.Debug("Unable to create tag", esl.Error(err))
		return err
	}
	if err == qt_errors.ErrorMock {
		return nil
	}
	l.Info("The tag created", esl.String("tag", tag.Ref))
	return nil
}

func (z *Publish) createReleaseDraft(c app_control.Control, relNote string) (rel *mo_release.Release, err error) {
	l := c.Log().With(
		esl.String("owner", app.RepositoryOwner),
		esl.String("repository", app.RepositoryName),
		esl.String("version", app.Version),
		esl.String("hash", app.Hash))
	ui := c.UI()

	relName := ""
	switch app.ReleaseStage() {
	case app.StageDev:
		relName = ui.Text(z.ReleaseNameDev.With("Version", app.Version))
	case app.StageBeta:
		relName = ui.Text(z.ReleaseNameBeta.With("Version", app.Version))
	case app.StageRelease:
		relName = ui.Text(z.ReleaseName.With("Version", app.Version))
	}

	svr := sv_release.New(z.ghCtx(c), app.RepositoryOwner, app.RepositoryName)
	rel, err = svr.CreateDraft(
		app.Version,
		relName,
		relNote,
		z.Branch,
	)
	if err != nil && err != qt_errors.ErrorMock {
		l.Debug("Unable to create release draft", esl.Error(err))
		return nil, err
	}
	if err == qt_errors.ErrorMock {
		return &mo_release.Release{}, nil
	}
	l.Info("Release created", esl.String("release", rel.Url))
	return rel, nil
}

func (z *Publish) uploadAssets(c app_control.Control, rel *mo_release.Release) error {
	l := c.Log()
	assets, _, err := z.artifactAssets(c)
	if err != nil {
		return err
	}

	sva := sv_release_asset.New(z.ghCtx(c), app.RepositoryOwner, app.RepositoryName, rel.Id)
	for _, p := range assets {
		l.Info("Uploading asset", esl.String("path", p))
		a, err := sva.Upload(mo_path2.NewExistingFileSystemPath(p))
		if err != nil && err != qt_errors.ErrorMock {
			return err
		}
		if err == qt_errors.ErrorMock {
			continue
		}
		l.Info("Uploaded", esl.Any("asset", a.Name))
		if strings.Contains(a.Name, "mac") {
			l.Info("updating Homebrew formula", esl.String("asset", a.Name))
			if err := z.updateHomebrewFormula(c, p); err != nil {
				return err
			}
		}
	}
	return nil
}

func (z *Publish) updateHomebrewFormula(c app_control.Control, path string) error {
	name := filepath.Base(path)
	return rc_exec.Exec(c, &homebrew.Formula{}, func(r rc_recipe.Recipe) {
		m := r.(*homebrew.Formula)
		m.Owner = homebrewRepoOwner
		m.Repository = homebrewRepoName
		m.Branch = homebrewRepoBranch
		m.AssetPath = mo_path2.NewExistingFileSystemPath(path)
		m.DownloadUrl = "https://github.com/watermint/toolbox/releases/download/" + app.Version + "/" + name
		m.Message = "Release " + app.Version
		m.FormulaName = "toolbox.rb"
	})
}

func (z *Publish) Exec(c app_control.Control) error {
	l := c.Log()
	ready := true

	if app.IsProduction() {
		l.Info("Verify embedded resources")
		qt_runtime.Suite(c)
	} else {
		l.Info("Run as dev mode")
		ready = false
	}

	if z.SkipTests {
		ready = false
	} else {
		err := z.endToEndTest(c)
		if err != nil {
			return err
		}
	}

	sum, err := z.verifyArtifacts(c)
	if err != nil {
		return err
	}

	relNote, err := z.releaseNotes(c, sum)
	if err != nil {
		return err
	}

	if err := z.createTag(c); err != nil {
		return nil
	}

	rel, err := z.createReleaseDraft(c, relNote)
	if err != nil {
		return err
	}

	if err := z.uploadAssets(c, rel); err != nil {
		return err
	}

	if !ready {
		l.Warn("The build does not satisfy release criteria")
		return ErrorBuildIsNotReadyForRelease
	}
	l.Info("The build is ready to publish")
	return nil
}

func (z *Publish) Test(c app_control.Control) error {
	d, err := qt_file.MakeTestFolder("release-publish", false)
	if err != nil {
		return err
	}

	platforms := []string{"linux", "mac", "win"}
	for _, platform := range platforms {
		app.Version = "dev-test"
		err = ioutil.WriteFile(filepath.Join(d, "tbx-"+app.Version+"-"+platform+".zip"), []byte("Test artifact"), 0644)
		if err != nil {
			c.Log().Warn("Unable to create test artifact", esl.Error(err))
			return err
		}
	}
	defer os.RemoveAll(d)

	err = rc_exec.ExecMock(c, &Publish{}, func(r rc_recipe.Recipe) {
		m := r.(*Publish)
		m.ArtifactPath = mo_path2.NewFileSystemPath(d)
	})
	if err, _ = qt_errors.ErrorsForTest(c.Log(), err); err != ErrorBuildIsNotReadyForRelease && err != nil {
		return err
	}
	return nil
}
