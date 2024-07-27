package entity

type Category struct {
	ID     uint
	Name   string
	Dreams []Dream `gorm:"many2many:categories_dreams;"`
}

type Categories []Category

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories" validate:"required"`
}

type CategoryResponse struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (c Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}
func (cats Categories) ToResponse() CategoriesResponse {
	resps := []CategoryResponse{}
	for _, c := range cats {
		resps = append(resps, c.ToResponse())
	}
	return CategoriesResponse{Categories: resps}
}
