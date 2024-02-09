package controller

import (
	"github.com/bbrighter/dreams-api/store"
)

var ErrorNotFound = store.ErrorNotFound

func (body DreamRequestBody) dreamRequestBodyToDream() store.Dream {
	var desc string
	if body.Description != nil {
		desc = *body.Description
	}
	return store.Dream{
		Date:        body.Date,
		Description: desc,
	}
}

func categoryToCategoryResponse(t store.Category) CategoryResponse {
	return CategoryResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}

func categoriesToCategoriesResponse(c []store.Category) CategoriesResponse {
	var categories []CategoryResponse = []CategoryResponse{}
	for _, category := range c {
		categories = append(categories, categoryToCategoryResponse(category))
	}
	return CategoriesResponse{Categories: categories}
}

func personToPersonResponse(p store.Person) PersonResponse {
	return PersonResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

func personsToPersonsResponse(p []store.Person) PersonsResponse {
	var persons []PersonResponse = []PersonResponse{}
	for _, person := range p {
		persons = append(persons, personToPersonResponse(person))
	}
	return PersonsResponse{Persons: persons}
}

func dreamToDreamMetaResponse(d store.Dream) DreamMetaResponse {
	return DreamMetaResponse{
		ID:   d.ID,
		Date: d.Date,
	}
}

func dreamToDreamResponse(d store.Dream) DreamResponse {
	var dreamMetaResponse DreamMetaResponse = dreamToDreamMetaResponse(d)
	var categories []CategoryResponse = []CategoryResponse{}
	for _, category := range d.Categories {
		categories = append(categories, categoryToCategoryResponse(category))
	}
	var persons []PersonResponse = []PersonResponse{}
	for _, person := range d.Persons {
		persons = append(persons, personToPersonResponse(person))
	}
	return DreamResponse{
		DreamMetaResponse: dreamMetaResponse,
		Description:       d.Description,
		Categories:        categories,
		Persons:           persons,
	}
}

func dreamsToDreamsResponse(dreams []store.Dream) DreamsResponse {
	var respList = []DreamMetaResponse{}
	for _, d := range dreams {
		dream := dreamToDreamMetaResponse(d)
		respList = append(respList, dream)
	}
	return DreamsResponse{Dreams: respList}
}

func countsToCounts(count []store.Count) []Count {
	var counts []Count
	for _, c := range count {
		cc := Count{
			ID:    c.ID,
			Count: c.Count,
		}
		counts = append(counts, cc)
	}
	return counts
}
