package app_msg_container_impl

import (
	"github.com/watermint/toolbox/essentials/lang"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/quality/infra/qt_missingmsg"
)

func NewMultilingual(las []lang.Lang, containers map[lang.Iso639One]app_msg_container.Container) app_msg_container.Container {
	return &mlContainer{
		priority:   las,
		containers: containers,
	}
}

type mlContainer struct {
	priority   []lang.Lang
	containers map[lang.Iso639One]app_msg_container.Container
}

func (z mlContainer) Text(key string) string {
	l := esl.Default()
	for _, la := range z.priority {
		if c, ok := z.containers[la.Code()]; ok {
			if c.ExistsKey(key) {
				return c.Text(key)
			}
		}
	}
	qt_missingmsg.Record().NotFound(key)
	l.Warn("Unable to find message resource", esl.String("key", key))
	return AltText(key)
}

func (z *mlContainer) Exists(msg app_msg.Message) bool {
	for _, la := range z.priority {
		if c, ok := z.containers[la.Code()]; ok {
			if c.Exists(msg) {
				return true
			}
		}
	}
	return false
}

func (z *mlContainer) ExistsKey(key string) bool {
	for _, la := range z.priority {
		if c, ok := z.containers[la.Code()]; ok {
			if c.ExistsKey(key) {
				return true
			}
		}
	}
	return false
}

func (z *mlContainer) Compile(m app_msg.Message) string {
	l := esl.Default()
	for _, la := range z.priority {
		if c, ok := z.containers[la.Code()]; ok {
			if c.Exists(m) {
				return c.Compile(m)
			}
		}
	}
	qt_missingmsg.Record().NotFound(m.Key())
	l.Warn("Unable to find message resource", esl.String("key", m.Key()))
	return AltCompile(m)
}
