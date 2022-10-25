package predb

import (
	"context"
	"time"
)

type Pre struct {
	At          time.Time
	Category    *Category
	CategoryRaw string
}

// predb.ovh allows for a maximum of 30 requests per minute
// so we keep track of the last request time and sleep if necessary
var lastRequest time.Time
var waitTimeBetweenRequests time.Duration = time.Second * 3

func GetRelease(ctx context.Context, rlsName string) (*Pre, error) {

	if time.Since(lastRequest) < waitTimeBetweenRequests {
		// this function will always be called from the same goroutine
		// so sleeping is not a problem
		time.Sleep(waitTimeBetweenRequests - time.Since(lastRequest))
	}
	lastRequest = time.Now()

	externalData, err := getExternalData(ctx, rlsName)
	if err != nil {
		return nil, err
	}

	// predb.ovh's categories are a total mess.
	// there are over 1,700 unique categories so we have to do some manual mapping
	normalizedExternalCategoryName := normalizeExternalCategory(externalData.Cat, externalData.URL)

	// now that we have the normalized category name, we can use it to find a category that fits our needs
	category := getCategory(externalData.Name, normalizedExternalCategoryName)

	pre := &Pre{
		At:          time.Unix(externalData.PreAt, 0),
		Category:    category,
		CategoryRaw: externalData.Cat,
	}

	return pre, nil
}
