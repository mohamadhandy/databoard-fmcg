package helper

type SyncMessage struct {
	Data      interface{} `json:"data"`
	Table     string      `json:"table"`
	Operation string      `json:"operation"`
}
