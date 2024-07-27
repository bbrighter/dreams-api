package entity

type Count struct {
	ID    uint
	Count int
}

type Counts []Count

type CountResponse struct {
	ID    uint `json:"id" validate:"required"`
	Count int  `json:"count" validate:"required"`
}

type CountsResponse struct {
	Categories []CountResponse `json:"categories" validate:"required"`
	Persons    []CountResponse `json:"persons" validate:"required"`
}

func (counts Counts) ToResponse() []CountResponse {
	var countResponses = []CountResponse{}
	for _, c := range counts {
		var countResponse = CountResponse(c)
		countResponses = append(countResponses, countResponse)
	}
	return countResponses
}
