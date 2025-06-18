package entity

type Count struct {
	ID    uint
	Count int
}

type Counts []Count

type CountResponse struct {
	ID    uint `json:"id" binding:"required"`
	Count int  `json:"count" binding:"required"`
}

type CountsResponse struct {
	Categories []CountResponse `json:"categories" binding:"required"`
	Persons    []CountResponse `json:"persons" binding:"required"`
}

func (counts Counts) ToResponse() []CountResponse {
	var countResponses = []CountResponse{}
	for _, c := range counts {
		var countResponse = CountResponse(c)
		countResponses = append(countResponses, countResponse)
	}
	return countResponses
}
