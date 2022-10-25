package category

import (
	"atus/backend/config"
	"atus/backend/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Name string

const (
	Movie   Name = "MOVIE"
	TV      Name = "TV"
	Docu    Name = "DOCU"
	App     Name = "APP"
	Game    Name = "GAME"
	Audio   Name = "AUDIO"
	EBook   Name = "EBOOK"
	XXX     Name = "XXX"
	Unknown Name = "UNKNOWN"
)

var allCategoryNames = []Name{Movie, TV, Docu, App, Game, Audio, EBook, XXX, Unknown}

type Category struct {
	Name     Name
	Enabled  bool
	Includes []string
	Excludes []string
	MaxSize  int64
}

const categoryEnabledConfigKey = "FILTERS__CATEGORY_%s_ENABLED"
const categoryIncludesConfigKey = "FILTERS__CATEGORY_%s_INCLUDES"
const categoryExcludesConfigKey = "FILTERS__CATEGORY_%s_EXCLUDES"
const categoryMaxSizeConfigKey = "FILTERS__CATEGORY_%s_MAX_SIZE"

func Get(name Name) (*Category, error) {

	// check if category is valid
	var found bool
	for _, n := range allCategoryNames {
		if n == name {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("unknown category name %s", name)
	}

	category := &Category{
		Name:    name,
		Enabled: config.GetBool(fmt.Sprintf(categoryEnabledConfigKey, name)),
		MaxSize: config.GetInt64(fmt.Sprintf(categoryMaxSizeConfigKey, name)),
	}

	// includes
	if includesStr := config.GetString(fmt.Sprintf(categoryIncludesConfigKey, name)); includesStr != "" {
		if err := json.Unmarshal([]byte(includesStr), &category.Includes); err != nil {
			return nil, err
		}
	}

	// excludes
	if excludesStr := config.GetString(fmt.Sprintf(categoryExcludesConfigKey, name)); excludesStr != "" {
		if err := json.Unmarshal([]byte(excludesStr), &category.Excludes); err != nil {
			return nil, err
		}
	}

	return category, nil

}

func GetAll() ([]*Category, error) {

	allCategoryNames := []Name{Movie, TV, XXX, Docu, App, Game, Audio, EBook, Unknown}

	var categories []*Category

	for _, name := range allCategoryNames {
		category, err := Get(name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil

}

func (c *Category) Save() error {

	includesBytes, err := json.Marshal(c.Includes)
	if err != nil {
		return err
	}

	excludesBytes, err := json.Marshal(c.Excludes)
	if err != nil {
		return err
	}

	config.Set(fmt.Sprintf(categoryEnabledConfigKey, c.Name), c.Enabled)
	config.Set(fmt.Sprintf(categoryIncludesConfigKey, c.Name), string(includesBytes))
	config.Set(fmt.Sprintf(categoryExcludesConfigKey, c.Name), string(excludesBytes))
	config.Set(fmt.Sprintf(categoryMaxSizeConfigKey, c.Name), c.MaxSize)

	return nil

}

// checks if a release is filtered by a category
func (c *Category) Accepts(rlsName string, rlsSize int64) (bool, error) {

	if !c.Enabled {
		return false, errors.New("category is disabled")
	}

	if c.MaxSize > 0 && rlsSize > c.MaxSize {
		return false, fmt.Errorf("release exceeds max size (%dGiB > %dGiB)", rlsSize/helpers.GiB, c.MaxSize/helpers.GiB)
	}

	lowerRlsName := strings.ToLower(rlsName)

	for _, include := range c.Includes {
		if strings.Contains(lowerRlsName, include) {
			return true, nil
		}
	}

	for _, exclude := range c.Excludes {
		if strings.Contains(lowerRlsName, exclude) {
			return false, fmt.Errorf("excluded by filter %s", exclude)
		}
	}

	return true, nil

}
