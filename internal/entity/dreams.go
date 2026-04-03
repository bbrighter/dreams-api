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
	Categories  Categories `gorm:"many2many:categories_dreams"`
	Rating      *int
}

type Dreams []Dream

type DreamsResponse struct {
	Dreams []DreamResponse `json:"dreams" binding:"required"`
}

// Response when querying one dream
type DreamResponse struct {
	ID          uint      `json:"id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Finalized   bool      `json:"finalized" binding:"required"`
	Visible     bool      `json:"visible" binding:"required"`
	Rating      *int      `json:"rating,omitempty"`
	Description string    `json:"description" binding:"required"`
	CategoriesResponse
}

func (d Dream) ToResponse() DreamResponse {
	return DreamResponse{
		ID:                 d.ID,
		Date:               d.Date,
		Visible:            d.Visible,
		Finalized:          d.Finalized,
		Rating:             d.Rating,
		Description:        d.Description,
		CategoriesResponse: d.Categories.ToResponse(),
	}
}

func (d Dreams) ToResponse() DreamsResponse {
	dreams := []DreamResponse{}
	for _, dream := range d {
		dreams = append(dreams, dream.ToResponse())
	}
	return DreamsResponse{Dreams: dreams}
}

type Includes int

const (
	IncludeCategories Includes = iota + 1
	IncludeDreamsCount
)

func ParseIncludes(raw string) []Includes {
	valid := map[string]Includes{
		"categories":  IncludeCategories,
		"dreamsCount": IncludeDreamsCount,
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
