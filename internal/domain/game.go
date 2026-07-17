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
}
