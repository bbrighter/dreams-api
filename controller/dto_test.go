package controller

import (
	"testing"
	"time"

	"github.com/bbrighter/dreams-api/store"
	"github.com/stretchr/testify/assert"
)

func TestDreamRequestBodyToDream(t *testing.T) {
	t.Parallel()
	var body DreamRequestBody
	var dream store.Dream
	var desc string = "desc"

	body = DreamRequestBody{
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: &desc,
	}
	dream = body.dreamRequestBodyToDream()

	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), dream.Date)
	assert.Equal(t, "desc", dream.Description)

	body = DreamRequestBody{
		Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	dream = body.dreamRequestBodyToDream()
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), dream.Date)
	assert.Equal(t, "", dream.Description)

}

func TestDreamToDreamResponse(t *testing.T) {
	t.Parallel()
	var dream store.Dream = store.Dream{
		ID:          1,
		Date:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "desc",
		Tags: []store.Tag{
			{ID: 1, Title: "Tag", Dreams: []store.Dream{{ID: 1}}},
		},
	}

	//
	var resp DreamResponse = dreamToDreamResponse(dream)
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "desc", resp.Description)
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), resp.Date)
	assert.Len(t, resp.Tags, 1)
	assert.Equal(t, "Tag", resp.Tags[0].Title)
	assert.EqualValues(t, 1, resp.Tags[0].ID)
}

func TestDreamsToDreamsResponse(t *testing.T) {
	t.Parallel()

	var dreams []store.Dream

	// Empty input gives {dreams: []}
	var emptyResp DreamsResponse = dreamsToDreamsResponse(dreams)
	assert.NotNil(t, emptyResp.Dreams)
	assert.Len(t, emptyResp.Dreams, 0)

	// Non-empty input gives {dreams: [...]}
	dreams = []store.Dream{
		{
			ID:          1,
			Date:        time.Now(),
			Description: "desc",
			Tags:        []store.Tag{{ID: 1, Title: "Tag"}},
		},
	}
	var resp DreamsResponse = dreamsToDreamsResponse(dreams)
	assert.Len(t, resp.Dreams, 1)
	assert.Len(t, resp.Dreams[0].Tags, 1)
}

func TestTagToTagResponse(t *testing.T) {
	t.Parallel()
	var tag store.Tag = store.Tag{ID: 1, Title: "Tag", Dreams: []store.Dream{{ID: 1}}}

	resp := tagToTagResponse(tag)

	assert.Equal(t, "Tag", resp.Title)
	assert.EqualValues(t, 1, resp.ID)
}

func TestTagsToTagsResponse(t *testing.T) {
	t.Parallel()

	var tags []store.Tag
	var tag1 store.Tag = store.Tag{ID: 1, Title: "Tag 1", Dreams: []store.Dream{{ID: 1}}}
	var tag2 store.Tag = store.Tag{ID: 2, Title: "Tag 2", Dreams: []store.Dream{{ID: 1}}}
	tags = append(tags, tag1, tag2)

	resp := tagsToTagsResponse(tags)

	assert.Len(t, resp.Tags, 2)
	assert.Equal(t, "Tag 1", resp.Tags[0].Title)
}
