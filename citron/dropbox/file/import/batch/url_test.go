package batch

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestUrl_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &Url{})
}
