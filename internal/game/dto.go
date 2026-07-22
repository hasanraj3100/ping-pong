package game

type ClientMsg struct {
	Type     string
	Value    string
	Sequence int
}

type Player2Joined struct {
	Type     string
	Opponent string
}

type MoveMsg struct {
	Type     string
	Player   string
	YPos     float32
	Sequence int
}

type ReadyMsg struct {
	Type         string
	Player1Ready bool
	Player2Ready bool
}
