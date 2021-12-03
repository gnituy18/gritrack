package step

type Step struct {
	Summary   string
	Items     []*Item
	CreatedAt int64
}

type ItemType int

const (
	ItemTypeTime ItemType = iota + 1
)

type Item struct {
	Type      ItemType
	Desc      string
	CreatedAt int64

	Time ItemTime
}

type ItemTime struct {
	Duration int
}
