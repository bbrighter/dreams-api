package entity

type Person struct {
	ID     uint
	Name   string
	Dreams []Dream `gorm:"many2many:people_dreams;"`
}

type Persons []Person

type PersonResponse struct {
	ID   uint   `json:"id"  validate:"required"`
	Name string `json:"name"  validate:"required"`
}

type PersonsResponse struct {
	Persons []PersonResponse `json:"persons" validate:"required"`
}

func (p Person) ToResponse() PersonResponse {
	return PersonResponse{ID: p.ID, Name: p.Name}
}

func (ps Persons) ToResponse() PersonsResponse {
	resps := []PersonResponse{}
	for _, p := range ps {
		resps = append(resps, p.ToResponse())
	}
	return PersonsResponse{Persons: resps}
}
