package repository

import (
	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

func (s *RepoTestSuite) TestListCategoriesCount() {
	var cats entity.CategoriesCount
	var err error
	cats, err = s.mgmt.List(s.ctx)
	s.NoError(err)
	s.Len(cats, 0)

	var cat1 = entity.Category{Type: entity.TypeCategory}
	s.db.Create(&cat1)
	cats, err = s.mgmt.List(s.ctx)
	s.NoError(err)
	s.Len(cats, 0)

	s.db.Create(&entity.Dream{Categories: entity.Categories{cat1}})
	cats, err = s.mgmt.List(s.ctx)
	s.NoError(err)
	s.Len(cats, 1)
	s.EqualValues(1, cats[0].Count)

	s.db.Create(&entity.Dream{Categories: entity.Categories{cat1}})
	cats, err = s.mgmt.List(s.ctx)
	s.NoError(err)
	s.Len(cats, 1)
	s.EqualValues(2, cats[0].Count)

	s.db.Create(&entity.Dream{Categories: entity.Categories{{Type: entity.TypePerson}}})
	cats, err = s.mgmt.List(s.ctx)
	s.NoError(err)
	s.Len(cats, 2)
	s.EqualValues(2, cats[0].Count)
	s.EqualValues(1, cats[1].Count)
}

func (s *RepoTestSuite) TestCategoryUpdate() {
	tests := map[string]struct {
		id           uint
		categoryType entity.CategoryType
		name         string
		expectedErr  error
	}{
		"not found": {
			id:           100,
			categoryType: entity.TypeCategory,
			expectedErr:  entity.ErrorNotFound,
		},
		"ok": {
			id:           1,
			categoryType: entity.TypeCategory,
			name:         "new name",
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			s.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Name: "cat", Type: entity.TypePerson}}})

			updates := make(map[string]any)
			if test.categoryType != "" {
				updates["type"] = test.categoryType
			}
			if test.name != "" {
				updates["name"] = test.name
			}
			err := s.mgmt.Update(s.ctx, test.id, updates)
			if test.expectedErr == nil {
				s.NoError(err)
			} else {
				s.ErrorIs(err, test.expectedErr)
			}
		})
	}
}

func (s *RepoTestSuite) TestCategoryDelete() {
	tests := map[string]struct {
		catId       uint
		expectedErr error
	}{
		"ok":           {catId: 10},
		"not found":    {catId: 100, expectedErr: gorm.ErrRecordNotFound},
		"still in use": {catId: 1, expectedErr: gorm.ErrForeignKeyViolated},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{entity.Category{ID: 1, Name: "cat", Type: entity.TypePerson}}})
			s.db.Create(&entity.Category{ID: 10, Type: entity.TypeCategory})

			err := s.mgmt.Delete(s.ctx, test.catId)
			if test.expectedErr != nil {
				s.ErrorIs(err, test.expectedErr)
				return
			}
			s.NoError(err)

			numberOfCats, _ := gorm.G[entity.Category](s.db).Where("id = ?", test.catId).Count(s.ctx, "*")
			s.EqualValues(0, numberOfCats)
		})
	}
}

func (s *RepoTestSuite) TestMerge() {
	type CategoryDream struct {
		CategoryID uint
		DreamID    uint
	}

	var cat1 = entity.Category{ID: 1, Name: "cat1", Type: entity.TypeCategory}
	var cat2 = entity.Category{ID: 2, Name: "cat2", Type: entity.TypeCategory}
	var cat3 = entity.Category{ID: 3, Name: "cat3", Type: entity.TypeCategory}

	tests := map[string]struct {
		fromId             uint
		toId               uint
		newName            string
		expectedError      error
		expectedCategories entity.Categories
		expectedRelations  []CategoryDream
	}{
		"ok": {fromId: 1, toId: 3, newName: "new",
			expectedCategories: entity.Categories{{ID: 3, Name: "new", Type: entity.TypeCategory}, cat2},
			expectedRelations:  []CategoryDream{{CategoryID: 3, DreamID: 2}, {CategoryID: 2, DreamID: 1}, {CategoryID: 3, DreamID: 1}},
		},
		"ok, but both cats in same dream": {fromId: 1, toId: 2, newName: "new",
			expectedCategories: entity.Categories{{ID: 2, Name: "new", Type: entity.TypeCategory}, cat3},
			expectedRelations:  []CategoryDream{{CategoryID: 3, DreamID: 2}, {CategoryID: 2, DreamID: 1}},
		},
	}
	for name, test := range tests {
		s.Run(name, func() {
			s.db.Create(&entity.Dream{ID: 1, Categories: entity.Categories{cat1, cat2}})
			s.db.Create(&entity.Dream{ID: 2, Categories: entity.Categories{cat3}})

			err := s.mgmt.Merge(s.ctx, test.fromId, test.toId, test.newName)
			if test.expectedError == nil {
				s.NoError(err)
			} else {
				s.ErrorContains(err, test.expectedError.Error())
			}
			var cats entity.Categories
			s.db.Find(&cats)
			s.ElementsMatch(test.expectedCategories, cats)

			var actual []CategoryDream
			s.db.Table("categories_dreams").Find(&actual)
			s.ElementsMatch(actual, test.expectedRelations)
		})

	}
}
