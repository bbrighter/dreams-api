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
	Rating      *int
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
	Rating    *int      `json:"rating,omitempty"`
	CategoriesResponse
}

// Response when querying one dream
type DreamResponse struct {
	DreamMetaResponse
	Description string `json:"description" binding:"required"`
}

func (d Dream) toMetaResponse() DreamMetaResponse {
	return DreamMetaResponse{
		ID:                 d.ID,
		Date:               d.Date,
		Visible:            d.Visible,
		Finalized:          d.Finalized,
		CategoriesResponse: d.Categories.ToResponse(),
		Rating:             d.Rating,
	}
}

func (d Dream) ToResponse() DreamResponse {
	return DreamResponse{
		DreamMetaResponse: d.toMetaResponse(),
		Description:       d.Description,
	}
}

func (d Dreams) ToResponse() DreamsResponse {
	dreams := []DreamMetaResponse{}
	for _, dream := range d {
		dreams = append(dreams, dream.toMetaResponse())
	}
	return DreamsResponse{Dreams: dreams}
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
