package main

import (
	"bowling-2/utils"
	"fmt"
	"time"
)

func NewManager(ttlLanes int) *Manager {
	freeLanes := make(chan int, ttlLanes)
	incPlayers := make(chan *Player)
	lanes := make([]Lane, ttlLanes)
	for i := 0; i < ttlLanes; i++ {
		lanes[i].Id = i
		freeLanes <- i
	}
	mgr := &Manager{
		Lanes:      lanes,
		FreeLanes:  freeLanes,
		IncPlayers: incPlayers,
	}
	return mgr
}

func (mgr *Manager) Run() {
	// var queue []*Player
	var avLanes []int
	for {
		mgr.mu.Lock()
		for len(mgr.Queue) > 0 && len(avLanes) > 0 {
			player := mgr.Queue[0]
			mgr.Queue = mgr.Queue[1:]

			laneID := avLanes[0]
			avLanes = avLanes[1:]

			go mgr.play(laneID, player)
		}
		mgr.mu.Unlock()

		select {
		case player := <-mgr.IncPlayers:
			mgr.mu.Lock()
			mgr.Queue = append(mgr.Queue, player)
			mgr.mu.Unlock()

		case laneID := <-mgr.FreeLanes:
			avLanes = append(avLanes, laneID)
		}
	}
}

func (mgr *Manager) play(laneid int, pl *Player) {
	mgr.mu.Lock()
	mgr.Lanes[laneid].Player = pl
	mgr.mu.Unlock()

	time.Sleep(time.Duration(pl.EstPlayTime * float64(time.Second)))
	th, er := utils.Inp(pl.Throws)
	if er != nil {
		fmt.Printf("ошибка игрока-%d %v", pl.Id, er)
	} else {
		sc, er := utils.Scr(th)
		if er != nil {
			fmt.Printf("ошибка игрока-%d %v", pl.Id, er)
		} else {
			pl.Score = sc
		}
	}
	mgr.mu.Lock()
	mgr.Lanes[laneid].Player = nil
	mgr.mu.Unlock()
	mgr.FreeLanes <- laneid
}
