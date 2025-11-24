package main

import (
	"fmt"
	"math/rand"
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

func genPlayers(mgr *Manager) {
	ids := 1
	for {
		thr := genRandGame()
		plTime := 0.5 + rand.Float64()*3.0
		p := &Player{
			Id:          ids,
			Throws:      thr,
			EstPlayTime: plTime,
		}
		fmt.Printf("новый игрок %d пришел на время %f\n", p.Id, p.EstPlayTime)
		mgr.IncPlayers <- p
		ids++
		wait := time.Duration(200+rand.Intn(800)) * time.Millisecond
		time.Sleep(wait)
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

	mgr := NewManager(3)
	go mgr.Run()
	genPlayers(mgr)
}
