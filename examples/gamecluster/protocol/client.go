package protocol

//login
type LoginRequest struct {
	//UId		int64	`json:"uId"`
	//Name 	string 	`json:"name"`
	//Password string `json:"password"`
	DeviceId string `json:"deviceId"`
}

type LoginResponse struct {
	UId		int64 	`json:"uId"`
	Diamond	int		`json:"diamond"`
	Name 	string 	`json:"name"`
	Pic 	string 	`json:"pic"`
}

type BindRequest struct {

}

type Ping struct {
	
}

type Pong struct {
	
}