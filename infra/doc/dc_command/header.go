package dc_command

import (
	"github.com/watermint/toolbox/infra/doc/dc_section"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
)

type MsgHeader struct {
	WarnDeprecation app_msg.Message
}

var (
	MHeader = app_msg.Apply(&MsgHeader{}).(*MsgHeader)
)

func NewHeader(spec rc_recipe.Spec) dc_section.Section {
	return &Header{
		spec: spec,
	}
}

type Header struct {
	spec rc_recipe.Spec
}

func (z Header) Title() app_msg.Message {
	return app_msg.Raw(z.spec.CliPath())
}

func (z Header) Body(ui app_ui.UI) {
	if z.spec.IsDeprecated() {
		ui.Info(MHeader.WarnDeprecation)
	}
	ui.Info(app_msg.Join(z.spec.Title(), z.spec.Remarks()))
	ui.Break()
	ui.Info(z.spec.Desc())
}
