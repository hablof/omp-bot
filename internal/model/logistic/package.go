package logistic

var AllEntities = []*Package{
	{Title: "Коробка"},
	{Title: "Ящик"},
	{Title: "Конверт"},
	{Title: "Пачка"},
	{Title: "Банка"},
	{Title: "Коробка для конфет"},
	{Title: "Ящик для конфет"},
	{Title: "Конверт для тонких конфет"},
	{Title: "Пачка для конфет"},
	{Title: "Банка для конфет"},
}

type Package struct {
	Title string
}

func (s *Package) String() string {
	return s.Title
}

// check index inbounds, index starts with 1
func CheckInbounds(realWorldIndex uint64) bool {
	if realWorldIndex < 1 {
		return false
	}
	if realWorldIndex >= uint64(len(AllEntities))+1 {
		return false
	}

	return true
}
