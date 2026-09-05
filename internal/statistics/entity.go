package statistics

type CountByCatAndMonth struct {
	CategoryId uint `gorm:"column:category_id"`
	CountByMonth
}

type CountByMonth struct {
	Month string `gorm:"column:month" binding:"required"`
	Count int64  `gorm:"column:count" binding:"required"`
}
type CountByMonths []CountByMonth
