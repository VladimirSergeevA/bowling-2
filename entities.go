package main

type Player struct {
	id            int
	throws        string
	est_play_time float64
	max_wait_time float64
	score         int
}
type Lane struct {
	id     int
	player *Player
}
