package logistic

var AllEntities = []*Package{
	{
		Title:         "Коробка",
		Material:      "Картон",
		MaximumVolume: 5000,
		Reusable:      true,
	},
	{
		Title:         "Ящик",
		Material:      "Дерево",
		MaximumVolume: 10000,
		Reusable:      false,
	},
	{
		Title:         "Конверт",
		Material:      "Бумага",
		MaximumVolume: 100,
		Reusable:      false,
	},
	{
		Title:         "Пачка",
		Material:      "Бумага",
		MaximumVolume: 75,
		Reusable:      false,
	},
	{
		Title:         "Банка",
		Material:      "Стекло",
		MaximumVolume: 3000,
		Reusable:      true,
	},
	{
		Title:         "Коробка для конфет",
		Material:      "Картон",
		MaximumVolume: 200,
		Reusable:      false,
	},
	{
		Title:         "Мешок для картофеля",
		Material:      "Полипропилен",
		MaximumVolume: 120000,
		Reusable:      true,
	},
	{
		Title:         "Жестяная банка из под сгухи",
		Material:      "Жесть",
		MaximumVolume: 380,
		Reusable:      true,
	},
}

type Package struct {
	Title         string
	Material      string
	MaximumVolume float32 //cm^3
	Reusable      bool
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
