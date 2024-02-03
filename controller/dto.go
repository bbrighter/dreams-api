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

func tagToTagResponse(t store.Tag) TagResponse {
	return TagResponse{
		ID:    t.ID,
		Title: t.Title,
	}
}

func tagsToTagsResponse(t []store.Tag) TagsResponse {
	var tags []TagResponse = []TagResponse{}
	for _, tag := range t {
		tags = append(tags, tagToTagResponse(tag))
	}
	return TagsResponse{Tags: tags}
}

func personToPersonResponse(p store.Person) PersonResponse {
	return PersonResponse{
		ID:   p.ID,
		Name: p.Name,
	}
}

func dreamToDreamMetaResponse(d store.Dream) DreamMetaResponse {
	return DreamMetaResponse{
		ID:   d.ID,
		Date: d.Date,
	}
}

func dreamToDreamResponse(d store.Dream) DreamResponse {
	var dreamMetaResponse DreamMetaResponse = dreamToDreamMetaResponse(d)
	var tags []TagResponse = []TagResponse{}
	for _, tag := range d.Tags {
		tags = append(tags, tagToTagResponse(tag))
	}
	var persons []PersonResponse = []PersonResponse{}
	for _, person := range d.Persons {
		persons = append(persons, personToPersonResponse(person))
	}
	return DreamResponse{
		DreamMetaResponse: dreamMetaResponse,
		Description:       d.Description,
		Tags:              tags,
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
