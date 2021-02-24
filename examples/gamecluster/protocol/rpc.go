package protocol

type NewUserRequest struct {
	UId  int64  `json:"uId"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}
