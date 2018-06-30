package lunch

import (
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

func starRecipe(orderlist []int, star int) error {
	data := loadData()
	for idx, recipe := range data.Baocan {
		for _, ol := range orderlist {
			if ol == idx {
				recipe.Star = (recipe.Star + star) / 2
			}
		}
	}
	ya, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(DataFilePath, ya, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
