package recipe

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestLicense_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &License{})
	t.Error("intended failure for test workflow")
}
