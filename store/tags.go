package store

type Tag struct {
	ID     uint
	Title  string
	Dreams []Dream `gorm:"many2many:tags_dreams;"`
}

func (repo Repo) GetTags() []Tag {
	var tags []Tag
	repo.db.Find(&tags)
	return tags
}

func (repo Repo) AddTagToDream(tagTitle string, dreamId uint) ([]Tag, error) {
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return nil, ErrorNotFound
	}

	var tag Tag = Tag{Title: tagTitle, Dreams: []Dream{dream}}
	repo.db.Where(&Tag{Title: tagTitle}).First(&tag)
	err := repo.db.Save(&tag).Error

	var tags []Tag
	repo.db.Find(&tags)
	return tags, err
}

func (repo Repo) RemoveTagFromDream(tagId uint, dreamId uint) ([]Tag, error) {
	var tags []Tag = []Tag{}

	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return tags, ErrorNotFound
	}
	var tag Tag = Tag{ID: tagId}
	if rowsAffected := repo.db.First(&tag).RowsAffected; rowsAffected == 0 {
		return tags, ErrorNotFound
	}

	repo.db.Model(&dream).Association("Tags").Delete(tag)
	tags = repo.removeTagsIfNeeded([]Tag{tag})

	return tags, nil
}

func (repo Repo) removeTagsIfNeeded(tags []Tag) []Tag {
	for _, tag := range tags {
		var usedTag Tag
		repo.db.Where(&Tag{Title: tag.Title}).Preload("Dreams").Find(&usedTag)
		if len(usedTag.Dreams) == 0 {
			repo.db.Delete(tag)
		}
	}
	repo.db.Find(&tags)
	return tags
}
