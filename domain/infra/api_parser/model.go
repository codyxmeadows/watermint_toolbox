package api_parser

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/app"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

func ParseModelString(v interface{}, j string) error {
	if !gjson.Valid(j) {
		return errors.New("invalid json")
	}
	g := gjson.Parse(j)
	return ParseModel(v, g)
}

func ParseModelRaw(v interface{}, j json.RawMessage) error {
	if !gjson.ValidBytes(j) {
		return errors.New("invalid json")
	}
	g := gjson.ParseBytes(j)
	return ParseModel(v, g)
}

func ParseModel(v interface{}, j gjson.Result) error {
	vv := reflect.ValueOf(v).Elem()
	vt := vv.Type()

	log := app.Root().Log().With(zap.String("valueType", vt.Name()))
	debug := app.Root().IsDebug()

	for i := vt.NumField() - 1; i >= 0; i-- {
		vtf := vt.Field(i)
		vvf := vv.Field(i)

		if vtf.Name == "Raw" && vvf.Type().Kind() == reflect.TypeOf(json.RawMessage{}).Kind() {
			vvf.SetBytes(json.RawMessage(j.Raw))
			continue
		}

		p := vtf.Tag.Get("path")
		if p == "" {
			continue
		}
		pp := strings.Split(p, ",")
		path := pp[0]
		required := false
		if len(pp) > 1 && pp[1] == "required" {
			required = true
		}

		jv := j.Get(path)
		if !jv.Exists() {
			if required {
				log.Error("Missing required field", zap.String("field", vtf.Name), zap.String("path", p))
				if debug {
					log.Debug("Entry JSON", zap.String("Entry", j.Raw))
				}
				return errors.New("missing required field")
			}
			continue
		}

		switch vtf.Type.Kind() {
		case reflect.String:
			vvf.SetString(jv.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vvf.SetInt(jv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vvf.SetUint(jv.Uint())
		case reflect.Bool:
			vvf.SetBool(jv.Bool())
		case reflect.Float32, reflect.Float64:
			vvf.SetFloat(jv.Float())

		default:
			log.Error("unexpected type found", zap.String("type.kind", vtf.Type.Kind().String()))
			return errors.New("unexpected type found")
		}
	}
	return nil
}
