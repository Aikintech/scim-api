package utils

type PaginationParams[T any] struct {
	Items    []T
	Search   func(item T) bool
	OrderBy  func(items []T, orderAsc bool)
	Page     int
	Limit    int
	OrderAsc bool
}

// PaginateItems paginates a list of items based on a search string, page number, and limit.
func PaginateItems[T any](params PaginationParams[T]) []T {
	// Filtering based on the search function
	var filteredItems []T
	for _, item := range params.Items {
		if params.Search(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	// Sorting based on the orderBy function
	if params.OrderBy != nil {
		params.OrderBy(filteredItems, params.OrderAsc)
	}

	// Paginating based on page and limit
	startIndex := (params.Page - 1) * params.Limit
	endIndex := startIndex + params.Limit

	// Ensure endIndex does not exceed the length of the slice
	if endIndex > len(filteredItems) {
		endIndex = len(filteredItems)
	}

	// Return the paginated items
	return filteredItems[startIndex:endIndex]
}
