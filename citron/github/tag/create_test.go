package tag

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestCreate_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &Create{})
}
