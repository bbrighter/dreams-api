package main

type Tag struct {
	ID      uint
	Title   string
	DreamID uint
}

func (repo Repo) createTag(tag Tag) (uint, error) {
	err := repo.db.Create(&tag).Error
	return tag.ID, err
}
