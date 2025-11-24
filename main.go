package main

import (
	"bowling-2/utils"
	"fmt"
)

func main() {
	input := "X X X X X X X X X X X 2"
	thr, err := utils.Inp(input)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}
	sc, err := utils.Scr(thr)
	if err != nil {
		fmt.Println("Ошибка подсчёта:", err)
	}
	fmt.Printf("Res: %d", sc)
}
