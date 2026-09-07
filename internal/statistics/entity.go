package statistics

type CountByCatAndMonth struct {
	ResultType string `gorm:"column:result_type"`
	CategoryId *uint  `gorm:"column:category_id"`
	Month      string `gorm:"column:month" binding:"required"`
	Count      int64  `gorm:"column:count" binding:"required"`
}
