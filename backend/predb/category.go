package predb

import "atus/backend/category"

type Category struct {
	Name category.Name
	Info interface{}
}

// We try to find a fitting category for the release
func getCategory(rlsName string, preDBCategoryName category.Name) *Category {

	// XXX
	if xxx := isXXX(rlsName); xxx != nil {
		return &Category{
			Name: category.XXX,
			Info: xxx,
		}
	}

	// Video
	video := isVideo(rlsName)
	isVideo := video != nil

	// TV Show
	if t := isTVSeries(rlsName, isVideo); t != nil {
		return &Category{
			Name: category.TV,
			Info: t,
		}
	}

	if isVideo {
		// We dont have season / episode infos but the predb thinks it's a tv show
		if preDBCategoryName == category.TV {
			return &Category{
				Name: category.TV,
			}
		}

		return &Category{
			Name: category.Movie,
		}
	}

	// We could not find a fitting category, so we return the predb category
	return &Category{
		Name: preDBCategoryName,
	}
}
