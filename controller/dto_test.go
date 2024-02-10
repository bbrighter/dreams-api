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
		Categories: []store.Category{
			{ID: 1, Name: "Category", Dreams: []store.Dream{{ID: 1}}},
		},
		Persons: []store.Person{
			{ID: 1, Name: "Name", Dreams: []store.Dream{{ID: 1}}},
		},
	}

	//
	var resp DreamResponse = dreamToDreamResponse(dream)
	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "desc", resp.Description)
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), resp.Date)
	assert.Len(t, resp.Categories, 1)
	assert.Equal(t, "Category", resp.Categories[0].Name)
	assert.EqualValues(t, 1, resp.Categories[0].ID)
	assert.Len(t, resp.Persons, 1)
	assert.EqualValues(t, 1, resp.Persons[0].ID)
	assert.Equal(t, "Name", resp.Persons[0].Name)
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
			Categories:  []store.Category{{ID: 1, Name: "Category"}},
			Persons:     []store.Person{{ID: 1, Name: "Name"}},
		},
	}
	var resp DreamsResponse = dreamsToDreamsResponse(dreams)
	assert.Len(t, resp.Dreams, 1)
}

func TestCategoryToCategoryResponse(t *testing.T) {
	t.Parallel()
	var category store.Category = store.Category{ID: 1, Name: "Category", Dreams: []store.Dream{{ID: 1}}}

	resp := categoryToCategoryResponse(category)

	assert.Equal(t, "Category", resp.Name)
	assert.EqualValues(t, 1, resp.ID)
}

func TestCategoiessToCategoriesResponse(t *testing.T) {
	t.Parallel()

	var categories []store.Category
	var category1 store.Category = store.Category{ID: 1, Name: "Category 1", Dreams: []store.Dream{{ID: 1}}}
	var category2 store.Category = store.Category{ID: 2, Name: "Category 2", Dreams: []store.Dream{{ID: 1}}}
	categories = append(categories, category1, category2)

	resp := categoriesToCategoriesResponse(categories)

	assert.Len(t, resp.Categories, 2)
	assert.Equal(t, "Category 1", resp.Categories[0].Name)
}

func TestPersonToPersonResponse(t *testing.T) {
	t.Parallel()

	var person store.Person = store.Person{ID: 1, Name: "Name", Dreams: []store.Dream{{ID: 1}}}

	resp := personToPersonResponse(person)

	assert.EqualValues(t, resp.ID, 1)
	assert.Equal(t, resp.Name, "Name")
}

func TestPersonsToPersonsResponse(t *testing.T) {
	t.Parallel()

	person1 := store.Person{ID: 1, Name: "Name", Dreams: []store.Dream{{ID: 1}}}
	person2 := store.Person{ID: 2, Name: "Name 2", Dreams: []store.Dream{{ID: 1}}}
	persons := []store.Person{person1, person2}

	resp := personsToPersonsResponse(persons)

	assert.Len(t, resp.Persons, 2)
}

func TestCategoryCountsToCategoryCountResponse(t *testing.T) {
	t.Parallel()

	count1 := store.Count{ID: 1, Count: 10}
	count2 := store.Count{ID: 2, Count: 20}
	counts := []store.Count{count1, count2}

	resp := countsToCounts(counts)

	assert.Len(t, resp, 2)
	assert.EqualValues(t, 1, resp[0].ID)
	assert.EqualValues(t, 10, resp[0].Count)
	assert.EqualValues(t, 2, resp[1].ID)
	assert.EqualValues(t, 20, resp[1].Count)
}
