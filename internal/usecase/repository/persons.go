package repository

import (
	"github.com/bbrighter/dreams-api/internal/entity"
	"gorm.io/gorm"
)

type PersonsRepo struct {
	Repo *gorm.DB
}

func NewPersonsRepo(db *gorm.DB) *PersonsRepo {
	return &PersonsRepo{Repo: db}
}

func (r *PersonsRepo) List() entity.Persons {
	var persons entity.Persons
	r.Repo.Find(&persons)
	return persons
}

func (r *PersonsRepo) AddToDream(name string, dream entity.Dream) (entity.Persons, error) {
	var persons entity.Persons
	if rowsAffected := r.Repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return persons, entity.ErrorNotFound
	}

	var person = entity.Person{Name: name, Dreams: entity.Dreams{dream}}
	r.Repo.Where(&entity.Person{Name: name}).First(&person)

	err := r.Repo.Save(&person).Error
	if err != nil {
		return persons, err
	}

	r.Repo.Find(&persons)

	return persons, nil
}

func (r *PersonsRepo) RemoveFromDream(person entity.Person, dream entity.Dream) (entity.Persons, error) {
	var persons entity.Persons
	if rowsAffected := r.Repo.First(&dream).RowsAffected; rowsAffected == 0 {
		return persons, entity.ErrorNotFound
	}
	if rowsAffected := r.Repo.First(&person).RowsAffected; rowsAffected == 0 {
		return persons, entity.ErrorNotFound
	}

	if err := r.Repo.Model(&dream).Association("Persons").Delete(person); err != nil {
		return persons, err
	}

	var usedPerson entity.Person
	r.Repo.Where(&entity.Person{Name: person.Name}).Preload("Dreams").Find(&usedPerson)
	if len(usedPerson.Dreams) == 0 {
		if err := r.Repo.Delete(usedPerson).Error; err != nil {
			return persons, err
		}
	}
	r.Repo.Find(&persons)

	return persons, nil
}
