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
	Persons     Persons    `gorm:"many2many:people_dreams;"`
}

type Dreams []Dream

type DreamsResponse struct {
	Dreams []DreamMetaResponse `json:"dreams" validate:"required"`
}

type DreamMetaResponse struct {
	ID         uint               `json:"id" validate:"required"`
	Date       time.Time          `json:"date" validate:"required"`
	Finalized  bool               `json:"finalized" validate:"required"`
	Visible    bool               `json:"visible" validate:"required"`
	Categories []CategoryResponse `json:"categories,omitempty" validate:"optional"`
	Persons    []PersonResponse   `json:"persons,omitempty" validate:"optional"`
}

// Response when querying one dream
type DreamResponse struct {
	DreamMetaResponse
	Description string `json:"description" validate:"required"`
}

func (d Dream) ToResponse() DreamResponse {
	return DreamResponse{
		DreamMetaResponse: DreamMetaResponse{
			ID:         d.ID,
			Date:       d.Date,
			Visible:    d.Visible,
			Finalized:  d.Finalized,
			Categories: d.Categories.ToList(),
			Persons:    d.Persons.ToList(),
		},
		Description: d.Description,
	}
}

func (d Dreams) ToResponse() DreamsResponse {
	resps := []DreamMetaResponse{}
	for _, dream := range d {
		resps = append(resps,
			DreamMetaResponse{
				ID:         dream.ID,
				Date:       dream.Date,
				Visible:    dream.Visible,
				Finalized:  dream.Finalized,
				Categories: dream.Categories.ToList(),
				Persons:    dream.Persons.ToList(),
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
