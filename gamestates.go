package main

type State int
type GameStates struct {
	Menu, Main, Pause State
}

var States GameStates = GameStates{0, 1, 2}
