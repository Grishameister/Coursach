package domain

type Error int

const (
	Success Error = 1 + iota
	NotFound
	KeyExist
)

type Route int

const (
	WriteStateData Route = 1 + iota
	ReadStatData
)

type Counters struct {
	People     int `json:"people" msgpack:"people"`
	Excavators int `json:"excavators" msgpack:"excavators"`
	Bulldozers int `json:"bulldozers" msgpack:"buldozers"`
}

type DataResponse struct {
	Date      string   `msgpack:"date" json:"date"`
	Counters  Counters `msgpack:"counters" json:"counters"`
	TypeRoute Route    `msgpack:"type" json:"type_route"`
	Error     Error    `msgpack:"error" json:"error"`
}

type DataRequest struct {
	Date      string   `msgpack:"date"`
	Counters  Counters `msgpack:"counters"`
	TypeRoute Route    `msgpack:"type"`
}

const (
	StatusOK                  = Status("OK")
	StatusNotEnoughBuldozers  = Status("NotEnoughBuldozers")
	StatusNotEnoughExcavators = Status("NotEnoughExcavators")
	StatusNotEnoughPeople     = Status("NotEnoughPeople")
)

type StatusRequest struct {
	Statuses string `json:"statuses" msgpack:"status"`
	Date     string `json:"date" msgpack:"date"`
}

type StatusChannel struct {
	Date     string   `json:"date"`
	Statuses []string `json:"statuses"`
}

type Status string
type Statuses map[Status]struct{}
