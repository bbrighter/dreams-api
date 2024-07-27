package usecase

import "github.com/bbrighter/dreams-api/internal/entity"

type PersonsUseCase struct {
	repo PersonsRepo
}

func NewPersonsUseCase(r PersonsRepo) *PersonsUseCase {
	return &PersonsUseCase{repo: r}
}

func (u *PersonsUseCase) GetAll() entity.Persons {
	return u.repo.GetAll()
}

func (u *PersonsUseCase) AddToDream(personName string, dreamId uint) (entity.Persons, error) {
	var dream = entity.Dream{ID: dreamId}
	return u.repo.AddToDream(personName, dream)
}

func (u *PersonsUseCase) RemoveFromDream(personId uint, dreamId uint) (entity.Persons, error) {
	var person = entity.Person{ID: personId}
	var dream = entity.Dream{ID: dreamId}
	return u.repo.RemoveFromDream(person, dream)
}
