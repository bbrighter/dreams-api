package entity

import "time"

type Dream struct {
	ID          uint
	Date        time.Time
	Description string
	Visible     bool       `gorm:"default:true"`
	Categories  Categories `gorm:"many2many:categories_dreams;"`
	Persons     Persons    `gorm:"many2many:people_dreams;"`
}

type Dreams []Dream

type DreamsResponse struct {
	Dreams []DreamMetaResponse `json:"dreams" validate:"required"`
}

type DreamMetaResponse struct {
	ID      uint      `json:"id" validate:"required"`
	Date    time.Time `json:"date" validate:"required"`
	Visible bool      `json:"visible" validate:"required"`
}

// Response when querying one dream
type DreamResponse struct {
	DreamMetaResponse
	Description string             `json:"description" validate:"required"`
	Categories  CategoriesResponse `json:"categories" validate:"required"`
	Persons     PersonsResponse    `json:"persons" validate:"required"`
}

func (d Dream) ToResponse() DreamResponse {
	return DreamResponse{
		DreamMetaResponse: DreamMetaResponse{
			ID:      d.ID,
			Date:    d.Date,
			Visible: d.Visible,
		},
		Description: d.Description,
		Categories:  d.Categories.ToResponse(),
		Persons:     d.Persons.ToResponse(),
	}
}

func (d Dreams) ToResponse() DreamsResponse {
	resps := []DreamMetaResponse{}
	for _, dream := range d {
		resps = append(resps,
			DreamMetaResponse{
				ID:      dream.ID,
				Date:    dream.Date,
				Visible: dream.Visible,
			},
		)
	}
	return DreamsResponse{Dreams: resps}
}
