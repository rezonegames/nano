package protocol

type TableInfo struct {
	TableId		int32 	`json:"tableId"`
	Name 		string	`json:"name"`
	Desc 		string	`json:"desc"`
	Owner		int64	`json:"owner"`
	OwnerName	string	`json:"ownerName"`
}

type CreateTableRequest struct {
	Name	string `json:"name"`
	Desc	string `json:"desc"`
}

type CreateTableResponse struct {
	Code		int			`json:"code"`
	TableInfo	TableInfo	`json:"tableInfo"`
}

type JoinTableRequest struct {
	TableId	int32 	`json:"tableId"`
}

type JoinTableResponse struct {
	Code	int    		`json:"code"`
	Table TableInfo 	`json:"table"`
}

type OnStart struct {
}

type Reward struct {
	Items 	[]Item	`json:"item"`
	UId 	int64 	`json:"uId"`
}

type StopAndSettleBroadcast struct {
	Rewards 	[]Reward 	`json:"rewards"`
}

type ReadyRequest struct {
	
}

type ReadyResponse struct {
	Code 	int 	`json:"code"`
}

type OverRequest struct {
	
}

type OverResponse struct {
	Code 	int 	`json:"code"`
}

type OnCountdown struct {
	Countdown	int `json:"countdown"`
}

type OnReady struct {
	User UserInfo	`json:"user"`
}

type OnJoinTable struct {
	User UserInfo 	`json:"user"`
}

type OnOver struct {
	User UserInfo	`json:"user"`
}

type UpdateMessage struct {
	ID       int     `json:"id"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
}