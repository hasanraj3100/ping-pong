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

type BallMsg struct {
	Type string
	X    float32
	Y    float32
	VX   float32
	VY   float32
}

type ScoreMsg struct {
	Type         string
	Player1Score int
	Player2Score int
}

type GameOverMsg struct {
	Type   string
	Winner string
}

type PresenceMsg struct {
	Type      string
	Player    string
	Connected bool
}

type StateMsg struct {
	Type         string
	Player1Score int
	Player2Score int
	Player1Ready bool
	Player2Ready bool
	Player1YPos  float32
	Player2YPos  float32
	BallX        float32
	BallY        float32
	BallVX       float32
	BallVY       float32
	Winner       string
}
