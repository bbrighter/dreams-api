package entity

import (
	"strings"
	"time"
)

type Dream struct {
	ID          uint
	Date        time.Time
	Description string
	Visible     bool `gorm:"default:true"`
	Finalized   bool
	Categories  Categories `gorm:"many2many:categories_dreams;"`
}

type Dreams []Dream

type DreamsResponse struct {
	Dreams []DreamMetaResponse `json:"dreams" binding:"required"`
}

type DreamMetaResponse struct {
	ID        uint      `json:"id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Finalized bool      `json:"finalized" binding:"required"`
	Visible   bool      `json:"visible" binding:"required"`
	CategoriesResponse
}

// Response when querying one dream
type DreamResponse struct {
	DreamMetaResponse
	Description string `json:"description" binding:"required"`
}

func (d Dream) ToResponse() DreamResponse {
	return DreamResponse{
		DreamMetaResponse: DreamMetaResponse{
			ID:                 d.ID,
			Date:               d.Date,
			Visible:            d.Visible,
			Finalized:          d.Finalized,
			CategoriesResponse: d.Categories.ToResponse(),
		},
		Description: d.Description,
	}
}

func (d Dreams) ToResponse() DreamsResponse {
	resps := []DreamMetaResponse{}
	for _, dream := range d {
		resps = append(resps,
			DreamMetaResponse{
				ID:                 dream.ID,
				Date:               dream.Date,
				Visible:            dream.Visible,
				Finalized:          dream.Finalized,
				CategoriesResponse: dream.Categories.ToResponse(),
			},
		)
	}
	return DreamsResponse{Dreams: resps}
}

type Includes int

const (
	IncludePersons Includes = iota
	IncludeCategories
)

func ParseIncludes(raw string) []Includes {
	valid := map[string]Includes{
		"categories": IncludeCategories,
		"persons":    IncludePersons,
	}

	set := make(map[Includes]struct{})
	for _, part := range strings.Split(raw, ",") {
		if include, ok := valid[part]; ok {
			set[include] = struct{}{}
		}
	}

	var includes []Includes
	for key := range set {
		includes = append(includes, key)
	}
	return includes
}
