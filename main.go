package main

import (
	"bowling-2/utils"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func genRandGame() string {
	chunks := []string{"X ", "9- ", "5/ ", "72 ", "81 ", "- - ", "X ", "5/ "}
	endl := []string{"X X X", "9-", "5/ 5", "11"}
	res := ""
	for i := 0; i < 9; i++ {
		res += chunks[rand.Intn(len(chunks))]
	}
	res += endl[rand.Intn(len(endl))]
	return res
}

func genPlayers(mgr *Manager, count int) {
	ids := 1
	for i := 0; i < count; i++ {
		thr := genRandGame()
		plTime := 0.5 + rand.Float64()*3.0
		p := &Player{
			Id:          ids,
			Throws:      thr,
			EstPlayTime: plTime,
			MaxWaitTime: float64(rand.Intn(3) + 1),
		}
		// fmt.Printf("новый игрок %d пришел на время %f\n", p.Id, p.EstPlayTime)
		mgr.IncPlayers <- p
		ids++
		wait := time.Duration(200+rand.Intn(200)) * time.Millisecond
		time.Sleep(wait)
	}
}

func display(mgr *Manager, done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
		}
		fmt.Print("\033[H\033[2J")
		mgr.mu.Lock()
		for _, lane := range mgr.Lanes {
			if lane.Player == nil {
				fmt.Printf("дорожка %d \t свободна\n", lane.Id)
			} else {
				p := lane.Player
				now := time.Now()
				diff := now.Sub(p.StartTime)
				sec := diff.Seconds()
				all := len(p.Throws)
				var showLen int
				if sec >= p.EstPlayTime {
					showLen = all
				} else {
					ratio := sec / p.EstPlayTime
					showLen = int(float64(all) * ratio)
				}
				if showLen > all {
					showLen = all
				}
				if showLen < 0 {
					showLen = 0
				}
				currentThrows := p.Throws[0:showLen]
				th, _ := utils.Inp(currentThrows)
				currentScore := utils.ScrPart(th)
				fmt.Printf("дорожка %d \t игрок-%d \t счет: %d \t табло: %-30s\n", lane.Id, p.Id, currentScore, currentThrows)
			}
		}
		fmt.Println()
		fmt.Print("очередь: ")
		if len(mgr.Queue) == 0 {
			fmt.Print("empty")
		} else {
			for i, p := range mgr.Queue {
				if i < 10 {
					fmt.Printf("%d ", p.Id)
				} else {
					fmt.Printf("... (еще %d в очереди)", len(mgr.Queue)-i)
					break
				}
			}
		}
		if mgr.LastLeftId > 0 {
			fmt.Printf("\n...игрок-%d вышел из очереди\n", mgr.LastLeftId)
			// mgr.LastLeftId = 0
		} else {
			fmt.Println()
		}
		mgr.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func printLeaderboard(mgr *Manager) {
	var players []*Player
	for _, p := range mgr.FinishedPlayers {
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	fmt.Println("\n\n топ игроков за сегодня")

	// for i, p := range players  {
	for i := 0; i < len(players); i++ {
		p := players[i]
		fmt.Printf("%d. Player %d \t Score: %d\n", i+1, p.Id, p.Score)
	}
}

func main() {
	// input := "X X X X X X X X X X X 2"
	// thr, err := utils.Inp(input)
	// if err != nil {
	// 	fmt.Println("Ошибка ввода:", err)
	// 	return
	// }
	// sc, err := utils.Scr(thr)
	// if err != nil {
	// 	fmt.Println("Ошибка подсчёта:", err)
	// }
	// fmt.Printf("Res: %d", sc)
	cnt := 50

	mgr := NewManager(5)
	done := make(chan bool)
	go mgr.Run()
	go display(mgr, done)
	genPlayers(mgr, cnt)

	for mgr.GetFinishedCount() < cnt {
		time.Sleep(500 * time.Millisecond)
	}

	done <- true
	// time.Sleep(200 * time.Millisecond)

	printLeaderboard(mgr)
}
