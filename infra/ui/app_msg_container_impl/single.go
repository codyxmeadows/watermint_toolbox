package app_msg_container_impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/watermint/toolbox/essentials/go/es_lang"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/control/app_resource"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/quality/infra/qt_msgusage"
	"strings"
	"text/template"
)

// Load single language container
func NewSingle(la es_lang.Lang) (c app_msg_container.Container, err error) {
	return newFromBytes(la, func(name string) (i []byte, e error) {
		return app_resource.Bundle().Messages().Bytes(name)
	})
}

func NewSingleWithMessagesForTest(msgs map[string]string) app_msg_container.Container {
	return &sglContainer{
		messages: msgs,
		msgLang:  es_lang.Default,
	}
}

func newFromBytes(la es_lang.Lang, loader func(name string) ([]byte, error)) (c app_msg_container.Container, err error) {
	l := esl.Default().With(esl.String("lang", la.String()))

	resName := fmt.Sprintf("messages%s.json", la.Suffix())
	l = l.With(esl.String("name", resName))
	resData, err := loader(resName)
	if err != nil {
		return nil, err
	}
	resMsgs := make(map[string]string)
	if err = json.Unmarshal(resData, &resMsgs); err != nil {
		l.Error("Unable to unmarshal message resource", esl.Error(err))
		return nil, err
	}

	return &sglContainer{
		messages: resMsgs,
		msgLang:  la,
	}, nil
}

// Single language container.
type sglContainer struct {
	messages map[string]string
	msgLang  es_lang.Lang
}

func (z sglContainer) Lang() es_lang.Lang {
	return z.msgLang
}

func (z sglContainer) Text(key string) string {
	if msg, ok := z.messages[key]; !ok {
		qt_msgusage.Record().NotFound(key)
		return AltText(key)
	} else {
		qt_msgusage.Record().Touch(key)
		return msg
	}
}

func (z sglContainer) Exists(msg app_msg.Message) bool {
	switch m := msg.(type) {
	case app_msg.MessageComplex:
		for _, mm := range m.Messages() {
			if !z.ExistsKey(mm.Key()) {
				return false
			}
		}
		return true

	default:
		return z.ExistsKey(msg.Key())
	}
}

func (z sglContainer) ExistsKey(key string) bool {
	_, ok := z.messages[key]
	return ok
}

func (z sglContainer) compileMessage(m app_msg.Message, msg string) string {
	l := esl.Default()
	key := m.Key()
	params := make(map[string]interface{})
	for _, p := range m.Params() {
		for k, v := range p {
			params[k] = v
		}
	}
	t, err := template.New(key).Parse(msg)
	if err != nil {
		l.Warn("Unable to compile message",
			esl.String("key", key),
			esl.String("msg", msg),
			esl.Error(err),
		)
		return AltCompile(m)
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, params); err != nil {
		l.Warn("Unable to format message",
			esl.String("key", key),
			esl.String("msg", msg),
			esl.Error(err),
		)
		return AltCompile(m)
	}

	return buf.String()
}

func (z sglContainer) compileComplex(messages []app_msg.Message) string {
	compiled := make([]string, 0)
	for _, msg := range messages {
		compiled = append(compiled, z.Compile(msg))
	}
	return strings.Join(compiled, " ")
}

func (z sglContainer) Compile(m app_msg.Message) string {
	key := m.Key()
	switch m0 := m.(type) {
	case app_msg.MessageComplex:
		qt_msgusage.Record().Touch(key)
		return z.compileComplex(m0.Messages())

	case app_msg.MessageOptional:
		if msg, ok := z.messages[key]; !ok {
			if m0.Optional() {
				return ""
			} else {
				qt_msgusage.Record().NotFound(key)
				return AltCompile(m)
			}
		} else {
			return z.compileMessage(m, msg)
		}

	default:
		if msg, ok := z.messages[key]; !ok {
			qt_msgusage.Record().NotFound(key)
			return AltCompile(m)
		} else {
			return z.compileMessage(m, msg)
		}
	}
}
