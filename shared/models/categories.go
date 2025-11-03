package models

// /categories
type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type CategoriesDbResponse struct {
	CategoriesResponse
}

type Category struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SingleCategoryResponse struct {
	*Category
}
