package el_ja

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/watermint/toolbox/essentials/cache/ec_file"
	"github.com/watermint/toolbox/essentials/log/esl"
	"os"
)

type DictionaryContainer interface {
	// Load loads dictionary by name.
	Load(name string) (d *dict.Dict, err error)

	// LoadIpa loads IPA dictionary.
	LoadIpa() (d *dict.Dict, err error)

	// LoadUni loads UniDic dictionary.
	LoadUni() (d *dict.Dict, err error)

	// NewTokenizer creates a tokenizer.
	NewTokenizer(name string, omitBosEos bool) (t *tokenizer.Tokenizer, err error)

	// NewIpaTokenizer creates a tokenizer with IPA dictionary.
	NewIpaTokenizer(omitBosEos bool) (t *tokenizer.Tokenizer, err error)

	// NewUniTokenizer creates a tokenizer with UniDic dictionary.
	NewUniTokenizer(omitBosEos bool) (t *tokenizer.Tokenizer, err error)
}

func NewContainer(cache ec_file.File, logger esl.Logger) DictionaryContainer {
	return &dictionaryContainerImpl{
		cache:  cache,
		logger: logger,
	}
}

type dictionaryContainerImpl struct {
	cache  ec_file.File
	logger esl.Logger
}

func (z dictionaryContainerImpl) NewTokenizer(name string, omitBosEos bool) (t *tokenizer.Tokenizer, err error) {
	dic, err := z.Load(name)
	if err != nil {
		return nil, err
	}
	options := make([]tokenizer.Option, 0)
	if omitBosEos {
		options = append(options, tokenizer.OmitBosEos())
	}
	return tokenizer.New(dic, options...)
}

func (z dictionaryContainerImpl) NewIpaTokenizer(omitBosEos bool) (t *tokenizer.Tokenizer, err error) {
	return z.NewTokenizer("ipa", omitBosEos)
}

func (z dictionaryContainerImpl) NewUniTokenizer(omitBosEos bool) (t *tokenizer.Tokenizer, err error) {
	return z.NewTokenizer("uni", omitBosEos)
}

func (z dictionaryContainerImpl) LoadIpa() (d *dict.Dict, err error) {
	return z.Load("ipa")
}

func (z dictionaryContainerImpl) LoadUni() (d *dict.Dict, err error) {
	return z.Load("uni")
}

func (z dictionaryContainerImpl) Load(name string) (d *dict.Dict, err error) {
	l := z.logger.With(esl.String("name", name))
	var url string
	switch name {
	case "ipa":
		url = "https://raw.githubusercontent.com/watermint/kagome-dict/main/ipa/ipa.dict"
	case "uni":
		url = "https://raw.githubusercontent.com/watermint/kagome-dict/main/uni/uni.dict"
	default:
		return nil, fmt.Errorf("unknown dictionary: %s", name)
	}
	path, err := z.cache.Get("nlp-ja-kagome", name, url)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		l.Debug("Unable to read file", esl.Error(err))
		return nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		l.Debug("Unable to open zip", esl.Error(err))
		return nil, err
	}
	d, err = dict.Load(zr, true)
	if err != nil {
		l.Debug("Unable to load dictionary", esl.Error(err))
		return nil, err
	}
	return d, nil
}
