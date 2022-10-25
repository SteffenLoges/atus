package atus

import (
	"atus/backend/category"
	"atus/backend/logger"
)

func (a *ATUS) GetAllCategories() []*category.Category {
	var categories []*category.Category

	a.categories.Range(func(key, value interface{}) bool {
		categories = append(categories, value.(*category.Category))
		return true
	})

	return categories
}

func (a *ATUS) GetCategoryByName(name category.Name) *category.Category {
	c, ok := a.categories.Load(name)
	if !ok {
		return nil
	}
	return c.(*category.Category)
}

func (a *ATUS) UpdateCategory(c *category.Category) error {
	if err := c.Save(); err != nil {
		return err
	}

	a.categories.Store(c.Name, c)

	logger.Infof("category %s updated", c.Name)

	return nil
}
