package qt_errors

import (
	"errors"
	"github.com/watermint/toolbox/essentials/file/es_filesystem"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/control/app_definitions"
)

var (
	// Marker error: Skip end to end test
	ErrorSkipEndToEndTest = errors.New("skip end to end test")

	// Marker error: The test requires human interaction. Will not do automated test
	ErrorHumanInteractionRequired = errors.New("human interaction required")

	// Marker error: The test requires no test
	ErrorNoTestRequired = errors.New("no test required")

	// Marker error: The test will be done separately as a scenario test.
	ErrorScenarioTest = errors.New("scenario test")

	// Marker error: The test is not yet implemented
	ErrorImplementMe = errors.New("implement me")

	// Marker error: The test requires some resource, but the resource is not available.
	ErrorNotEnoughResource = errors.New("not enough resource")

	// Unsupported UI
	ErrorUnsupportedUI = errors.New("unsupported UI for this auth scope")

	// Marker error: Mock
	ErrorMock = errors.New("mock error")
)

// Returns nil even err != nil if the error type is ignorable.
func ErrorsForTest(l esl.Logger, err error) (resolvedErr error, cont bool) {
	if err == nil {
		return nil, true
	}

	if fsErr, ok := err.(es_filesystem.FileSystemError); ok {
		if fsErr.IsMockError() {
			l.Debug("Mock error")
			return nil, false
		}
	}

	switch err {
	case ErrorSkipEndToEndTest:
		l.Debug("Skip: skip end to end test")
		return nil, false

	case ErrorNoTestRequired:
		l.Debug("Skip: No test required for this recipe")
		return nil, false

	case ErrorHumanInteractionRequired:
		l.Debug("Skip: Human interaction required for this test")
		return nil, false

	case ErrorNotEnoughResource:
		l.Debug("Skip: Not enough resource")
		return nil, false

	case ErrorScenarioTest:
		l.Debug("Skip: Implemented as scenario test")
		return nil, false

	case ErrorImplementMe:
		l.Debug("Test is not implemented for this recipe")
		return nil, false

	case ErrorUnsupportedUI:
		l.Debug("Test is not compatible for testing UI")
		return nil, false

	case ErrorMock:
		l.Debug("Mock test")
		return nil, false

	case app_definitions.ErrorUserCancelled:
		l.Debug("User cancelled")
		return nil, false

	default:
		return err, false
	}
}
