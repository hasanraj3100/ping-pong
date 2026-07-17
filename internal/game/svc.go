package game

import (
	"hasanraj3100/ping-pong/internal/domain"
)

const (
	paddleStep = 10
	maxPaddleY = 220
)

type GameService struct {
	repo        GameRepo
	broadcaster Broadcaster
}

func NewGameService(repo GameRepo, broadcaster Broadcaster) *GameService {
	return &GameService{repo: repo, broadcaster: broadcaster}
}

func (s *GameService) CreateGame(creator domain.User) domain.Game {
	g := s.repo.Create(creator)
	s.broadcaster.Broadcast(g.ID, g)
	return g
}

func (s *GameService) GetGame(gameID string) (domain.Game, bool) {
	return s.repo.FindByID(gameID)
}

func (s *GameService) JoinGame(gameID string, joiner domain.User) (domain.Game, error) {
	g, ok := s.repo.FindByID(gameID)
	if !ok {
		return domain.Game{}, ErrGameNotFound
	}

	if g.Player1.UUID == joiner.UUID {
		return domain.Game{}, ErrAlreadyJoined
	}

	if g.Player2.UUID != "" {
		return domain.Game{}, ErrGameFull
	}

	g.Player2 = joiner
	updated := s.repo.Update(g)
	s.broadcaster.Broadcast(updated.ID, updated)
	return updated, nil
}

func clampY(y float32) float32 {
	if y < 0 {
		return 0
	}
	if y > maxPaddleY {
		return maxPaddleY
	}
	return y
}

func (s *GameService) HandleClientMsg(msg ClientMsg, user domain.User, room string) {
	if msg.Type != "move" {
		return
	}

	var delta float32
	switch msg.Value {
	case "u":
		delta = -paddleStep
	case "d":
		delta = paddleStep
	default:
		return
	}

	g, ok := s.repo.FindByID(room)
	if !ok {
		return
	}

	var role string
	switch user.UUID {
	case g.Player1.UUID:
		role = "player1"
		g.State.Player1YPosition = clampY(g.State.Player1YPosition + delta)
	case g.Player2.UUID:
		role = "player2"
		g.State.Player2YPosition = clampY(g.State.Player2YPosition + delta)
	default:
		return // not a participant in this game
	}

	updated := s.repo.Update(g)

	yPos := updated.State.Player1YPosition
	if role == "player2" {
		yPos = updated.State.Player2YPosition
	}

	s.broadcaster.Broadcast(room, MoveMsg{Player: role, YPos: yPos})
}
