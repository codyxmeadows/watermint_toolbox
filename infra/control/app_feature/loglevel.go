package app_feature

import (
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/runtime/es_env"
	"github.com/watermint/toolbox/infra/app"
)

func ConsoleLogLevel(test, debug bool) esl.Level {
	switch {
	case es_env.IsEnabled(app.EnvNameDebugVerbose):
		return esl.LevelDebug

	case es_env.IsEnabled(app.EnvNameTestQuiet):
		return esl.LevelWarn

	case app.IsProduction() && test:
		return esl.LevelQuiet

	case debug:
		return esl.LevelDebug

	default:
		return esl.LevelInfo
	}
}
