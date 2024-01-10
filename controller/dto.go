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

func dreamToDreamResponse(d store.Dream) DreamResponse {
	var tags []TagResponse = []TagResponse{}
	for _, tag := range d.Tags {
		tags = append(tags, tagToTagResponse(tag))
	}
	return DreamResponse{
		ID:          d.ID,
		Date:        d.Date,
		Description: d.Description,
		Tags:        tags,
	}
}

func dreamsToDreamsResponse(dreams []store.Dream) DreamsResponse {
	var respList = []DreamResponse{}
	for _, d := range dreams {
		dream := dreamToDreamResponse(d)
		respList = append(respList, dream)
	}
	return DreamsResponse{Dreams: respList}
}
