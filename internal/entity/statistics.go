package entity

type CountByCat struct {
	CategoryId uint  `json:"id" gorm:"column:category_id"`
	Count      int64 `json:"count" gorm:"column:count"`
}

type CategoriesCount []CountByCat

type CategoriesCountResponse struct {
	Categories []CountByCat `json:"categories"`
}

func (cc CategoriesCount) ToResponse() CategoriesCountResponse {
	return CategoriesCountResponse{Categories: cc}
}

type CountByCatAndMonth struct {
	CategoryId uint `gorm:"column:category_id"`
	CountByMonth
}

type CountByMonth struct {
	Month string `gorm:"column:month"`
	Count int64  `gorm:"column:count"`
}
type CountByMonths []CountByMonth

type Statistic struct {
	Month      string       `json:"month"`
	Categories []CountByCat `json:"categories"`
	DreamCount int64        `json:"dreamCount"`
}

type Statistics []Statistic

func (cs CountByMonths) ToResponse(ms []CountByCatAndMonth) Statistics {
	var statistics = Statistics{}
	for _, c := range cs {
		for _, m := range ms {
			var cats = []CountByCat{}
			if c.Month == m.Month {
				cats = append(cats, CountByCat{CategoryId: m.CategoryId, Count: m.Count})
			}
			statistics = append(statistics, Statistic{Month: c.Month, Categories: cats, DreamCount: c.Count})
		}
	}
	return statistics
}
