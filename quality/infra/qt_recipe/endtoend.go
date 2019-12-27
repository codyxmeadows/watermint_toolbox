package qt_recipe

import (
	"encoding/csv"
	rice "github.com/GeertJohan/go.rice"
	"github.com/pkg/profile"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_impl"
	"github.com/watermint/toolbox/infra/control/app_root"
	"github.com/watermint/toolbox/infra/network/nw_ratelimit"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/infra/ui/app_msg_container_impl"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"github.com/watermint/toolbox/infra/util/ut_memory"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"github.com/watermint/toolbox/quality/infra/qt_missingmsg_impl"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	TestTeamFolderName = "watermint-toolbox-test"
)

func Resources(t *testing.T) (bx, web *rice.Box, mc app_msg_container.Container, ui app_ui.UI) {
	bx = rice.MustFindBox("../../../resources")
	web = rice.MustFindBox("../../../web")

	mc = app_msg_container_impl.NewContainer(bx)
	ui = app_ui.NewConsole(mc, qt_missingmsg_impl.NewMessageTest(t), true)
	return
}

func findTestResource() (resource gjson.Result, found bool) {
	l := app_root.Log()
	p, found := os.LookupEnv("TOOLBOX_TESTRESOURCE")
	if !found {
		return gjson.Parse("{}"), false
	}
	l = l.With(zap.String("path", p))
	b, err := ioutil.ReadFile(p)
	if err != nil {
		l.Debug("unable to read file", zap.Error(err))
		return gjson.Parse("{}"), false
	}
	if !gjson.ValidBytes(b) {
		l.Debug("invalid file content", zap.ByteString("resource", b))
		return gjson.Parse("{}"), false
	}
	return gjson.ParseBytes(b), true
}

func TestWithControl(t *testing.T, twc func(ctl app_control.Control)) {
	nw_ratelimit.SetTestMode(true)
	bx, web, mc, ui := Resources(t)

	ctl := app_control_impl.NewSingle(ui, bx, web, mc, false, []rc_recipe.Recipe{}, []rc_recipe.Recipe{})
	cs := ctl.(*app_control_impl.Single)
	if res, found := findTestResource(); found {
		var err error
		ctl, err = cs.NewTestControl(res)
		if err != nil {
			t.Error("Unable to create new test control", err)
			return
		}
	}
	err := ctl.Up(app_control.Test(), app_control.Concurrency(runtime.NumCPU()))
	if err != nil {
		os.Exit(app_control.FatalStartup)
	}
	defer ctl.Down()

	twc(ctl)
}

func RecipeError(l *zap.Logger, err error) error {
	switch err.(type) {
	case *qt_endtoend.ErrorNoTestRequired:
		l.Info("Skip: No test required for this recipe")
		return nil

	case *qt_endtoend.ErrorHumanInteractionRequired:
		l.Info("Skip: Human interaction required for this test")
		return nil

	case *qt_endtoend.ErrorNotEnoughResource:
		l.Info("Skip: Not enough resource")
		return nil

	case *qt_endtoend.ErrorScenarioTest:
		l.Info("Skip: Implemented as scenario test")
		return nil

	case *qt_endtoend.ErrorImplementMe:
		l.Warn("Test is not implemented for this recipe")
		return nil

	default:
		return err
	}
}

func TestRecipe(t *testing.T, re rc_recipe.Recipe) {
	nw_ratelimit.SetTestMode(true)
	TestWithControl(t, func(ctl app_control.Control) {
		l := ctl.Log()
		l.Debug("Start testing")
		pr := profile.Start(
			profile.ProfilePath(ctl.Workspace().Log()),
			profile.MemProfile,
		)

		err := re.Test(ctl)

		pr.Stop()
		ut_memory.DumpStats(l)

		if err == nil {
			return
		}

		if re := RecipeError(l, err); re != nil {
			t.Error(re)
		}
	})
}

type RowTester func(cols map[string]string) error

func TestRows(ctl app_control.Control, reportName string, tester RowTester) error {
	l := ctl.Log().With(zap.String("reportName", reportName))
	csvFile := filepath.Join(ctl.Workspace().Report(), reportName+".csv")

	l.Debug("Start loading report", zap.String("csvFile", csvFile))

	cf, err := os.Open(csvFile)
	if err != nil {
		l.Warn("Unable to open report CSV", zap.Error(err))
		return err
	}
	defer cf.Close()
	csf := csv.NewReader(cf)
	var header []string
	isFirstLine := true

	for {
		cols, err := csf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Warn("An error occurred during read report file", zap.Error(err))
			return err
		}
		if isFirstLine {
			header = cols
			isFirstLine = false
		} else {
			colMap := make(map[string]string)
			for i, h := range header {
				colMap[h] = cols[i]
			}
			if err := tester(colMap); err != nil {
				l.Warn("Tester returned an error", zap.Error(err), zap.Any("cols", colMap))
				return err
			}
		}
	}

	return nil
}
