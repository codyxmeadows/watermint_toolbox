package es_filepath

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_desktop"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/control/app_definitions"
	"github.com/watermint/toolbox/infra/security/sc_random"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"
)

const (
	VarDropboxPersonal           = "DropboxPersonal"
	VarDropboxBusiness           = "DropboxBusiness"
	VarDropboxBusinessOrPersonal = "DropboxBusinessOrPersonal"
	VarDropboxPersonalOrBusiness = "DropboxPersonalOrBusiness"
	VarHome                      = "Home"
	VarUsername                  = "Username"
	VarHostname                  = "Hostname"
	VarExecPath                  = "ExecPath"
	VarRand8                     = "Rand8"
	VarYear                      = "Year"
	VarMonth                     = "Month"
	VarDay                       = "Day"
	VarDate                      = "Date"
	VarTime                      = "Time"
	VarDateUTC                   = "DateUTC"
	VarTimeUTC                   = "TimeUTC"
	VarRandom                    = "Random"
)

var (
	PathVariables = []string{
		VarDropboxPersonal,
		VarDropboxBusiness,
		VarDropboxBusinessOrPersonal,
		VarDropboxPersonalOrBusiness,
		VarHome,
		VarUsername,
		VarHostname,
		VarExecPath,
		VarRand8,
		VarYear,
		VarMonth,
		VarDay,
		VarDate,
		VarTime,
		VarDateUTC,
		VarTimeUTC,
	}
	isWindows = app_definitions.IsWindows()
)

func Rel(basePath, targetPath string) (rel string, err error) {
	l := esl.Default()

	isSeparator := func(c rune) bool {
		switch {
		case c == '/', c == '\\':
			return true
		case c == ':' && isWindows:
			return true
		default:
			return false
		}
	}

	bpr := []rune(basePath)
	tpr := []rune(targetPath)

	bl := len(bpr)
	tl := len(tpr)

	l = l.With(esl.Int("basePathLen", bl), esl.Int("targetPathLen", tl))

	if bl < 1 || tl < 1 {
		l.Debug("Empty path")
		return "", errors.New("empty path")
	}

	if isSeparator(bpr[bl-1]) {
		bpr = bpr[:bl-1]
		bl = len(bpr)
	}
	if isSeparator(tpr[tl-1]) {
		tpr = tpr[:tl-1]
		tl = len(tpr)
	}

	if tl == bl {
		same := true
		for i := 0; i < tl; i++ {
			if unicode.ToLower(bpr[i]) != unicode.ToLower(tpr[i]) {
				same = false
			}
		}
		if same {
			return ".", nil
		}
	}
	if tl <= bl {
		l.Debug("Target path is shorter or than base path, or same length")
		return "", errors.New("target path is shorter than base path")
	}

	errMsg := "target path does not have same base path"

	for i := 0; i < bl; i++ {
		if unicode.ToLower(bpr[i]) != unicode.ToLower(tpr[i]) {
			return "", errors.New(errMsg)
		}
	}
	if isSeparator(bpr[bl-1]) {
		return string(tpr[bl:]), nil
	}
	if isSeparator(tpr[bl]) {
		return string(tpr[bl+1:]), nil
	}
	return "", errors.New(errMsg)
}

// Replace chars that is not usable for path with '_'
func Escape(p string) string {
	illegal := []string{
		"<",
		">",
		":",
		"\"",
		"|",
		"?",
		"*",
		"/",
		"\\",
	}

	o := p
	for _, il := range illegal {
		o = strings.ReplaceAll(o, il, "_")
	}
	return o
}

type FormatError struct {
	Reason string
	Key    string
}

func (z *FormatError) Error() string {
	return z.Value()
}
func (z *FormatError) Value() string {
	return "{{." + z.Key + "}}: " + z.Reason
}

type Pair struct {
	Key   string
	Value string
}

// Format path if a path contains pattern like `{{.DropboxPersonal}}`.
func FormatPathWithPredefinedVariables(path string, pairs ...Pair) (string, error) {
	predefined := make(map[string]func() (string, error))
	predefined[VarDropboxPersonal] = func() (s string, e error) {
		p, _, _ := sv_desktop.New().Lookup()
		if p != nil {
			return p.Path, nil
		}
		return "", errors.New("personal dropbox desktop folder not found")
	}
	predefined[VarDropboxBusiness] = func() (s string, e error) {
		_, p, _ := sv_desktop.New().Lookup()
		if p != nil {
			return p.Path, nil
		}
		return "", errors.New("business dropbox desktop folder not found")
	}
	predefined[VarDropboxBusinessOrPersonal] = func() (string, error) {
		p, b, _ := sv_desktop.New().Lookup()
		if b != nil {
			return b.Path, nil
		}
		if p != nil {
			return p.Path, nil
		}
		return "", errors.New("dropbox desktop folder not found")
	}
	predefined[VarDropboxPersonalOrBusiness] = func() (string, error) {
		p, b, _ := sv_desktop.New().Lookup()
		if p != nil {
			return p.Path, nil
		}
		if b != nil {
			return b.Path, nil
		}
		return "", errors.New("dropbox desktop folder not found")
	}
	predefined[VarHome] = func() (s string, e error) {
		u, err := user.Current()
		if err == nil {
			return u.HomeDir, nil
		}
		return "", errors.New("unable to retrieve current user home")
	}
	predefined[VarUsername] = func() (s string, e error) {
		h, err := user.Current()
		if err == nil {
			return Escape(h.Username), nil
		}
		return "", errors.New("unable to retrieve hostname")
	}
	predefined[VarHostname] = func() (s string, e error) {
		h, err := os.Hostname()
		if err == nil {
			return Escape(h), nil
		}
		return "", errors.New("unable to retrieve hostname")
	}
	predefined[VarExecPath] = func() (s string, err error) {
		s = filepath.Dir(os.Args[0])
		return s, nil
	}
	predefined[VarRand8] = func() (s string, err error) {
		return fmt.Sprintf("%08d", rand.Intn(100_000_000)), nil
	}
	predefined[VarYear] = func() (s string, e error) {
		s = time.Now().Local().Format("2006")
		return s, nil
	}
	predefined[VarMonth] = func() (s string, e error) {
		s = time.Now().Local().Format("01")
		return s, nil
	}
	predefined[VarDay] = func() (s string, e error) {
		s = time.Now().Local().Format("02")
		return s, nil
	}
	predefined[VarDate] = func() (s string, e error) {
		s = time.Now().Local().Format("2006-01-02")
		return s, nil
	}
	predefined[VarTime] = func() (s string, e error) {
		s = time.Now().Local().Format("15-04-05")
		return s, nil
	}
	predefined[VarDateUTC] = func() (s string, e error) {
		s = time.Now().UTC().Format("2006-01-02")
		return s, nil
	}
	predefined[VarTimeUTC] = func() (s string, e error) {
		s = time.Now().UTC().Format("15-04-05")
		return s, nil
	}
	predefined[VarRandom] = func() (s string, e error) {
		s = sc_random.MustGetSecureRandomString(6)
		return s, nil
	}
	predefined["AlwaysErrorForTest"] = func() (s string, e error) {
		return "", errors.New("always error")
	}
	data := make(map[string]string)
	for _, p := range pairs {
		data[p.Key] = p.Value
	}

	for k, vf := range predefined {
		ptn := "{{." + k + "}}"
		if strings.Index(path, ptn) >= 0 {
			v, err := vf()
			if err != nil {
				return "", &FormatError{
					Reason: err.Error(),
					Key:    k,
				}
			}
			data[k] = v
		}
	}

	if len(data) < 1 {
		return path, nil
	}

	var buf bytes.Buffer
	pathTmpl, err := template.New("path").Parse(path)
	if err != nil {
		return "", err
	}

	err = pathTmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
