package domain

import "time"

type GetStat struct {
	Date time.Time
}

type Counters struct {
	People     int `json:"people"`
	Excavators int `json:"excavators"`
	Bulldozers int `json:"bulldozers"`
}

type StatFromServer struct {
	Date     time.Time `json:"date"`
	Counters Counters  `json:"counters"`
}

const (
	StatusOK                  = Status("OK")
	StatusNotEnoughBuldozers  = Status("NotEnoughBuldozers")
	StatusNotEnoughExcavators = Status("NotEnoughExcavators ")
	StatusNotEnoughPeople     = Status("NotEnoughPeople")
)

type Status string
type Statuses map[Status]struct{}
