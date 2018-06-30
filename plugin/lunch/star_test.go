package lunch

import (
	"testing"
)

func TestStarRecipe(t *testing.T) {
	DataFilePath = "recipe.yaml"
	err := starRecipe([]int{1, 2, 3, 4}, 4)
	if err != nil {
		t.Error(err)
	}
}
