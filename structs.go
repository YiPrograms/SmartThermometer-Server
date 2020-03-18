package main

type config struct {
	Key      string
	SheetsID string
	TimeZone string
}

type card struct {
	UID  string
	Num  int
	Name string
}

type queryRequest struct {
	UID string
	Key string
}

type queryResponse struct {
	Num  int
	Name string
}

type registerRequest struct {
	UID string
	Num int
	Key string
}

type placeRequest struct {
	Num  int
	Temp float32
	Key  string
}
