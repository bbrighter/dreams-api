package entity

type CountByCat struct {
	CategoryId uint  `json:"id" gorm:"column:category_id" binding:"required"`
	Count      int64 `json:"count" gorm:"column:count" binding:"required"`
}

type CategoriesCount []CountByCat

type CategoriesCountResponse struct {
	Categories []CountByCat `json:"categories" binding:"required"`
}

func (cc CategoriesCount) ToResponse() CategoriesCountResponse {
	return CategoriesCountResponse{Categories: cc}
}

type CountByCatAndMonth struct {
	CategoryId uint `gorm:"column:category_id"`
	CountByMonth
}

type CountByMonth struct {
	Month string `gorm:"column:month" binding:"required"`
	Count int64  `gorm:"column:count" binding:"required"`
}
type CountByMonths []CountByMonth

type Statistic struct {
	Month      string       `json:"month" binding:"required"`
	Categories []CountByCat `json:"categories" binding:"required"`
	DreamCount int64        `json:"dreamCount" binding:"required"`
}

type Statistics struct {
	Statistics []Statistic `json:"statistics" binding:"required"`
}

func (cs CountByMonths) ToResponse(ms []CountByCatAndMonth) Statistics {
	var statistics = []Statistic{}
	for _, c := range cs {
		var cats = []CountByCat{}
		for _, m := range ms {
			if c.Month == m.Month {
				cats = append(cats, CountByCat{CategoryId: m.CategoryId, Count: m.Count})
			}
		}
		statistics = append(statistics, Statistic{Month: c.Month, Categories: cats, DreamCount: c.Count})
	}
	return Statistics{Statistics: statistics}
}
