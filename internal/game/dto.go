package game

type ClientMsg struct {
	Type  string
	Value string
}

type MoveMsg struct {
	Player string
	YPos   float32
}
