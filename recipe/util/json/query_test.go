package json

import (
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestQuery_Exec(t *testing.T) {
	qtr_endtoend.TestRecipe(t, &Query{})
}
