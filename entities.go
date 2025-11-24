package main

import "sync"

type Player struct {
	Id          int
	Throws      string
	EstPlayTime float64
	MaxWaitTime float64
	Score       int
}
type Lane struct {
	Id     int
	Player *Player
}
type Manager struct {
	Lanes      []Lane
	FreeLanes  chan int
	IncPlayers chan *Player
	Queue      []*Player
	mu         sync.Mutex
}
