package db

type User struct {
	UId      int64  `bson:"_id"`
	Diamond  int    `json:"diamond"`
	DeviceId string `json:"deviceid"`
	Name     string `json:"name"`
	Pic      string `json:"pic"`
}

type Counter struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}
