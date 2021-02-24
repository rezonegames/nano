package protocol

type ItemType int

const (
	Diamond    ItemType = 0
	Consumed            = 1
	UnConsumed          = 2
)

type Item struct {
	ItemType ItemType `json:"itemType"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Count    int      `json:"count"`
}
