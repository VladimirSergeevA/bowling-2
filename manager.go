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
			player.JoinTime = time.Now()
			mgr.Queue = append(mgr.Queue, player)
			mgr.mu.Unlock()

		case laneID := <-mgr.FreeLanes:
			avLanes = append(avLanes, laneID)

		case <-time.After(100 * time.Millisecond):
			mgr.mu.Lock()
			now := time.Now()
			n := 0
			for _, p := range mgr.Queue {
				if now.Sub(p.JoinTime).Seconds() > p.MaxWaitTime {
					mgr.LastLeftId = p.Id
					p.Score = -1
					mgr.FinishedPlayers = append(mgr.FinishedPlayers, p)
				} else {
					mgr.Queue[n] = p
					n++
				}
			}
			mgr.Queue = mgr.Queue[:n]
			mgr.mu.Unlock()
		}
	}
}

func (mgr *Manager) play(laneid int, pl *Player) {
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

	pl.StartTime = time.Now()

	mgr.mu.Lock()
	mgr.Lanes[laneid].Player = pl
	mgr.mu.Unlock()

	time.Sleep(time.Duration(pl.EstPlayTime * float64(time.Second)))

	mgr.mu.Lock()
	mgr.Lanes[laneid].Player = nil
	mgr.FinishedPlayers = append(mgr.FinishedPlayers, pl)
	mgr.mu.Unlock()
	mgr.FreeLanes <- laneid
}
