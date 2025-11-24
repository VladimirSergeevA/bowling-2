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
	var queue []*Player
	var avLanes []int
	for {
		for len(queue) > 0 && len(avLanes) > 0 {
			p := queue[0]
			queue = queue[1:]
			laneid := avLanes[0]
			avLanes = avLanes[1:]
			go mgr.play(laneid, p)
		}
		select {
		case pl := <-mgr.IncPlayers:
			queue = append(queue, pl)
		case laneid := <-mgr.FreeLanes:
			avLanes = append(avLanes, laneid)
		}
	}
}

func (mgr *Manager) play(laneid int, pl *Player) {
	mgr.Lanes[laneid].Player = pl
	time.Sleep(time.Duration(pl.EstPlayTime))
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
	mgr.Lanes[laneid].Player = nil
	mgr.FreeLanes <- laneid
}
