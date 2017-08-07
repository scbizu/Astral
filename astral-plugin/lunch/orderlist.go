package lunch

import (
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

//Recipe ...
type Recipe struct {
	Name  string   `yaml:"name"`
	Price int      `yaml:"price"`
	Star  int      `yaml:"star"`
	Tag   []string `yaml:"tag"`
}

//Restaurant ...
type Restaurant struct {
	Baocan []Recipe `yaml:"baocan"`
}

func getRecipe(count int) (orderLists []string) {
	r := loadData()
	meatCount := count / 2
	meat := getTagRecipe(r.Baocan, "meat")
	sortmeat := sortByStar(meat)
	vege := getTagRecipe(r.Baocan, "vege")
	sortvege := sortByStar(vege)
	//not enough meat
	if len(sortmeat) < meatCount {
		meatCount = len(sortmeat)
		for _, sr := range sortmeat {
			orderLists = append(orderLists, sr.Name)
		}
		//not enough vege
		if len(sortvege) < count-meatCount {
			for _, sr := range sortvege {
				orderLists = append(orderLists, sr.Name)
			}
		} else {
			for i := 0; i < count-meatCount; i++ {
				orderLists = append(orderLists, sortvege[i].Name)
			}
		}
	} else if len(sortvege) < count-meatCount {
		for _, sr := range sortvege {
			orderLists = append(orderLists, sr.Name)
		}
		if len(sortmeat) < count-len(sortvege) {
			for _, sr := range sortmeat {
				orderLists = append(orderLists, sr.Name)
			}
		} else {
			for i := 0; i < count-len(sortvege); i++ {
				orderLists = append(orderLists, sortmeat[i].Name)
			}
		}
	} else {
		for i := 0; i < meatCount; i++ {
			orderLists = append(orderLists, sortmeat[i].Name)
		}
		for i := 0; i < count-meatCount; i++ {
			orderLists = append(orderLists, sortvege[i].Name)
		}
	}
	return
}

func convertRecipe(raw []string) (res []int) {
	data := loadData()
	for idx, recipe := range data.Baocan {
		for _, order := range raw {
			if recipe.Name == order {
				res = append(res, idx)
			}
		}
	}
	return
}

func showRecipe(orders []int) (showRecipeStrs string) {
	r := loadData()
	for idx, recipe := range r.Baocan {
		for _, order := range orders {
			if idx == order {
				showRecipeStrs = showRecipeStrs + strconv.Itoa(idx) + " : " + recipe.Name + SPACE + " ¥" + strconv.Itoa(recipe.Price) + ENTER
			}
		}
	}
	return
}

func showRecipeByName(orders []string) (showRecipeStrs string) {
	if len(orders) > 0 {
		showRecipeStrs = showRecipeStrs + "Your Order List:" + ENTER
	}
	ordersStr := strings.Join(orders, "\n")
	showRecipeStrs = showRecipeStrs + ordersStr
	return
}

func getAllRecipeInfo() (showRecipeStrs string) {
	r := loadData()
	if len(r.Baocan) > 0 {
		showRecipeStrs = showRecipeStrs + "Order List: " + ENTER
	}
	for idx, recipe := range r.Baocan {
		showRecipeStrs = showRecipeStrs + strconv.Itoa(idx) + " : " + recipe.Name + SPACE + " ¥" + strconv.Itoa(recipe.Price) + SPACE + " ✨" + strconv.Itoa(recipe.Star) + ENTER
	}
	return
}

func loadData() *Restaurant {
	data, err := ioutil.ReadFile(DataFilePath)
	if err != nil {
		log.Println(err)
	}
	r := new(Restaurant)
	err = yaml.Unmarshal(data, r)
	if err != nil {
		log.Println(err)
	}
	return r
}

func getTagRecipe(re []Recipe, tag string) (res []Recipe) {
	for _, r := range re {
		tagStr := strings.Join(r.Tag, ",")
		if strings.Contains(tagStr, tag) {
			res = append(res, r)
		}
	}
	return
}

func sortByStar(recipes []Recipe) []Recipe {
	sort.Slice(recipes, func(i int, j int) bool {
		return recipes[i].Star > recipes[j].Star
	})
	return recipes
}
