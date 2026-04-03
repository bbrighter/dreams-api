package repository

import (
	"context"
	"time"

	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestListDreams() {
	tests := map[string]struct {
		createDreamBefore bool
		showAll           bool
		expectedLength    int
		catLength         int
	}{
		"no dreams":              {},
		"one dream":              {createDreamBefore: true, expectedLength: 1, catLength: 2},
		"include private dreams": {createDreamBefore: true, showAll: true, expectedLength: 2, catLength: 2},
	}

	for name, test := range tests {
		s.Run(name, func() {
			if test.createDreamBefore {
				s.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{
					entity.Category{ID: 1, Type: entity.TypeCategory},
					entity.Category{ID: 2, Type: entity.TypePerson},
				}})
				s.db.Create(&entity.Dream{ID: 2, Visible: false})
				s.db.Model(&entity.Dream{ID: 2}).UpdateColumn("visible", false)
			}

			dreams, err := s.dreams.List(s.ctx, test.showAll)
			s.NoError(err)
			s.Len(dreams, test.expectedLength)
			var dream1 entity.Dream
			for _, dream := range dreams {
				if dream.ID == 1 {
					dream1 = dream
				}
			}
			s.Len(dream1.Categories, test.catLength)
		})
	}
}

func (s *RepoTestSuite) TestUpdateDream() {
	tests := map[string]struct {
		dreamId   uint
		updateMap map[string]any
		isError   bool
	}{
		"date, ok":        {dreamId: 1, updateMap: map[string]any{"date": time.Now()}},
		"description, ok": {dreamId: 1, updateMap: map[string]any{"description": "desc"}},
		"date + desc, ok": {dreamId: 1, updateMap: map[string]any{"date": time.Now(), "description": "desc"}},
		"rating, ok":      {dreamId: 1, updateMap: map[string]any{"rating": 3}},
		"finalized, ok":   {dreamId: 1, updateMap: map[string]any{"finalized": true}},
		"not found":       {dreamId: 100, updateMap: map[string]any{"date": time.Now()}, isError: true},
		"invalid column":  {dreamId: 1, updateMap: map[string]any{"invalid": "abc"}, isError: true},
		"invalid type":    {dreamId: 1, updateMap: map[string]any{"rating": "a"}, isError: true},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var err error
			err = s.db.Create(&entity.Dream{ID: 1}).Error
			s.NoError(err)

			err = s.dreams.Update(s.ctx, test.dreamId, test.updateMap)

			if test.isError {
				s.Error(err)
				return
			}
			s.NoError(err)
		})
	}
}

func (s *RepoTestSuite) TestDeleteDream() {
	err := s.db.Create(&entity.Dream{
		ID: 1,
		Categories: entity.Categories{
			entity.Category{ID: 1, Name: "Cat", Type: entity.TypeCategory},
			entity.Category{ID: 2, Name: "Person", Type: entity.TypePerson},
		},
	}).Error
	s.NoError(err)

	err = s.dreams.Delete(context.Background(), 1)
	s.NoError(err)

	var count int64
	s.db.Table("categories_dreams").Count(&count)
	s.EqualValues(count, 0)
}

func (s *RepoTestSuite) TestToggleVisibility() {
	err := s.db.Create(&entity.Dream{ID: 1}).Error
	s.NoError(err)

	err = s.dreams.ToggleVisibility(s.ctx, 1)
	s.NoError(err)

	var dream entity.Dream
	s.db.First(&dream)
	s.False(dream.Visible)

	err = s.dreams.ToggleVisibility(s.ctx, 200)
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *RepoTestSuite) TestCreateDream() {
	id, err := s.dreams.Create(s.ctx, &entity.Dream{
		ID: 1,
		Categories: entity.Categories{
			entity.Category{ID: 1, Name: "Cat", Type: entity.TypeCategory},
			entity.Category{ID: 2, Name: "Person", Type: entity.TypePerson},
		},
	})
	s.NoError(err)
	s.EqualValues(1, id)
}

func (s *RepoTestSuite) TestGetDream() {
	tests := map[string]struct {
		useNonExistingId bool
		expectedError    error
	}{
		"not found": {useNonExistingId: true, expectedError: gorm.ErrRecordNotFound},
		"ok":        {},
	}
	for name, test := range tests {
		s.Run(name, func() {
			var id uint = s.CreateDream(true, false).ID
			if test.useNonExistingId {
				id = 1000
			}
			dream, err := s.dreams.Get(s.ctx, id, false)
			if test.expectedError != nil {
				s.ErrorIs(err, test.expectedError)
				return
			}
			s.NoError(err)
			s.Len(dream.Categories, 1)
		})
	}
}
