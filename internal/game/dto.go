package game

type ClientMsg struct {
	Type     string
	Value    string
	Sequence int
}

type MoveMsg struct {
	Player   string
	YPos     float32
	Sequence int
}
