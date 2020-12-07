package protocol

type JoinRoomRequest struct {
	Offset 	int		`json:"offset"`
}

type JoinRoomResponse struct {
	Code	int    			`json:"code"`
	//Result string `json:"result"`
	Tables []TableInfo 	`json:"tableInfo"`
}

type LeaveRoomRequest struct {
}

type LeaveRoomResponse struct {
	Code	int    	`json:"code"`
	UId 	int64 	`json:"uId"`
	Content string	`json:"content"`
}

type QuickStartRequest struct {
}

type QuickStartResponse struct {
	Code	int    			`json:"code"`
}

type CancelQuickStartRequest struct {
}

type CancelQuickStartResponse struct {
	Code	int    			`json:"code"`
}

type OnCancelQuickStart struct {
}