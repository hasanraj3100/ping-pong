package domain

type Game struct {
	ID      string
	Player1 User
	Player2 User
	State   GameData
}

type GameData struct {
	Player1YPosition float32
	Player2YPosition float32
	Player1Ready     bool
	Player2Ready     bool

	BallX  float32
	BallY  float32
	BallVX float32
	BallVY float32

	Player1Score int
	Player2Score int
	Winner       string
}
