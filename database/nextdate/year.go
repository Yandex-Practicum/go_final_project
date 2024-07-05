package nextdate

// calcY возвращает следующую дату, исходя из правила repeat "Y"
func calcY() string {
	nextDateDT := startDate.AddDate(1, 0, 0)
	for nextDateDT.Before(now) || nextDateDT.Equal(now) {
		nextDateDT = nextDateDT.AddDate(1, 0, 0)
	}
	nextDate = nextDateDT.Format(dateFormat)
	return nextDate
}
