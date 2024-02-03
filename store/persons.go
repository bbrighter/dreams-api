package store

type Person struct {
	ID     uint
	Name   string
	Dreams []Dream `gorm:"many2many:people_dreams;"`
}

func (repo Repo) GetPersonsForDream(dreamId uint) []Person {
	var persons []Person
	repo.db.Find(&persons)
	return persons
}

func (repo Repo) AddPersonToDream(name string, dreamId uint) error {
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return ErrorNotFound
	}

	var person Person = Person{Name: name, Dreams: []Dream{dream}}
	repo.db.Where(&Person{Name: name}).First(&person)

	return repo.db.Save(&person).Error
}

func (repo Repo) RemovePersonFromDream(personId uint, dreamId uint) error {
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return ErrorNotFound
	}
	var person Person = Person{ID: personId}
	if rowsAffected := repo.db.First(&person).RowsAffected; rowsAffected == 0 {
		return ErrorNotFound
	}

	if err := repo.db.Model(&dream).Association("Persons").Delete(person); err != nil {
		return err
	}

	var usedPerson Person
	repo.db.Debug().Where(&Person{Name: person.Name}).Preload("Dreams").Find(&usedPerson)
	if len(usedPerson.Dreams) == 0 {
		return repo.db.Debug().Delete(usedPerson).Error
	}

	return nil
}
