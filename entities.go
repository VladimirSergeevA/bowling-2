package main

import (
	"sync"
	"time"
)

type Player struct {
	Id          int
	Throws      string
	EstPlayTime float64
	MaxWaitTime float64
	Score       int
	StartTime   time.Time
	JoinTime    time.Time
}
type Lane struct {
	Id     int
	Player *Player
}
type Manager struct {
	Lanes           []Lane
	FreeLanes       chan int
	IncPlayers      chan *Player
	Queue           []*Player
	FinishedPlayers []*Player
	LastLeftId      int
	mu              sync.Mutex
}
