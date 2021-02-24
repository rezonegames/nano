package protocol

type UserInfo struct {
	UId     int64  `json:"uId"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
