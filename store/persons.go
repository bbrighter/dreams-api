package store

type Person struct {
	ID     uint
	Name   string
	Dreams []Dream `gorm:"many2many:people_dreams;"`
}

func (repo Repo) GetAllPersons() []Person {
	var persons []Person
	repo.db.Find(&persons)
	return persons
}

func (repo Repo) AddPersonToDream(name string, dreamId uint) (uint, error) {
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return 0, ErrorNotFound
	}

	var person Person = Person{Name: name, Dreams: []Dream{dream}}
	repo.db.Where(&Person{Name: name}).First(&person)

	err := repo.db.Save(&person).Error
	return person.ID, err
}

func (repo Repo) RemovePersonFromDream(personId uint, dreamId uint) ([]Person, error) {
	var persons []Person
	var dream Dream = Dream{ID: dreamId}
	if rowsAffected := repo.db.First(&dream).RowsAffected; rowsAffected == 0 {
		return persons, ErrorNotFound
	}
	var person Person = Person{ID: personId}
	if rowsAffected := repo.db.First(&person).RowsAffected; rowsAffected == 0 {
		return persons, ErrorNotFound
	}

	if err := repo.db.Model(&dream).Association("Persons").Delete(person); err != nil {
		return persons, err
	}

	var usedPerson Person
	repo.db.Where(&Person{Name: person.Name}).Preload("Dreams").Find(&usedPerson)
	if len(usedPerson.Dreams) == 0 {
		if err := repo.db.Debug().Delete(usedPerson).Error; err != nil {
			return persons, err
		}
	}
	repo.db.Find(&persons)

	return persons, nil
}
