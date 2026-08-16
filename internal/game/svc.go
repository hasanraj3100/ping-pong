package game

import (
	"math/rand/v2"
	"time"

	"hasanraj3100/ping-pong/internal/domain"
)

const (
	paddleStep = 10
	maxPaddleY = 220

	courtWidth  = 500
	courtHeight = 300

	paddleOffset = 10
	paddleWidth  = 15
	paddleHeight = 80

	ballInitialVX = 150
	ballInitialVY = 90

	paddleHitSpeedMultiplier = 1.1

	winningScore = 10

	tickRate = 50 * time.Millisecond
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
	s.broadcaster.Broadcast(updated.ID, Player2Joined{Type: "match_found", Opponent: updated.Player2.Name})
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
	switch msg.Type {
	case "move":
		s.HandleMove(msg, user, room)
	case "ready":
		s.HandleReady(msg, user, room)
	}
}

func (s *GameService) HandleMove(msg ClientMsg, user domain.User, room string) {
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

	s.broadcaster.Broadcast(room, MoveMsg{Type: "move", Player: role, YPos: yPos, Sequence: msg.Sequence})
}

func (s *GameService) HandleReady(msg ClientMsg, user domain.User, room string) {
	g, ok := s.repo.FindByID(room)
	if !ok {
		return
	}

	if g.State.Winner != "" {
		return
	}

	switch user.UUID {
	case g.Player1.UUID:
		if g.State.Player1Ready {
			return
		}
		g.State.Player1Ready = true
	case g.Player2.UUID:
		if g.State.Player2Ready {
			return
		}
		g.State.Player2Ready = true
	default:
		return
	}

	starting := g.State.Player1Ready && g.State.Player2Ready
	if starting {
		resetBall(&g.State)
	}

	updated := s.repo.Update(g)
	s.broadcaster.Broadcast(room, ReadyMsg{Type: "ready", Player1Ready: updated.State.Player1Ready, Player2Ready: updated.State.Player2Ready})

	if starting {
		go s.runBallLoop(room)
	}
}

func (s *GameService) runBallLoop(room string) {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()
	dt := float32(tickRate.Seconds())

	for range ticker.C {
		g, ok := s.repo.FindByID(room)
		if !ok {
			return
		}

		state := &g.State
		state.BallX += state.BallVX * dt
		state.BallY += state.BallVY * dt

		for state.BallY < 0 || state.BallY > courtHeight {
			if state.BallY < 0 {
				state.BallY = -state.BallY
				state.BallVY = -state.BallVY
			} else if state.BallY > courtHeight {
				state.BallY = 2*courtHeight - state.BallY
				state.BallVY = -state.BallVY
			}
		}

		if state.BallX < 0 {
			state.Player2Score++
		} else if state.BallX > courtWidth {
			state.Player1Score++
		}

		if state.BallX < 0 || state.BallX > courtWidth {
			resetBall(state)
			state.Player1Ready = false
			state.Player2Ready = false

			if state.Player1Score >= winningScore {
				state.Winner = "player1"
			} else if state.Player2Score >= winningScore {
				state.Winner = "player2"
			}

			updated := s.repo.Update(g)
			s.broadcaster.Broadcast(room, ScoreMsg{Type: "score", Player1Score: updated.State.Player1Score, Player2Score: updated.State.Player2Score})
			// ball is parked at center until both players ready up again, so report it at rest
			s.broadcaster.Broadcast(room, BallMsg{Type: "ball", X: updated.State.BallX, Y: updated.State.BallY, VX: 0, VY: 0})

			if updated.State.Winner != "" {
				s.broadcaster.Broadcast(room, GameOverMsg{Type: "game_over", Winner: updated.State.Winner})
			} else {
				s.broadcaster.Broadcast(room, ReadyMsg{Type: "ready", Player1Ready: updated.State.Player1Ready, Player2Ready: updated.State.Player2Ready})
			}
			return
		}

		leftPaddleEdge := float32(paddleOffset + paddleWidth)
		rightPaddleEdge := float32(courtWidth - paddleOffset - paddleWidth)

		if state.BallVX < 0 && state.BallX <= leftPaddleEdge && paddleHit(state.BallY, state.Player1YPosition) {
			state.BallX = leftPaddleEdge
			state.BallVX = -state.BallVX * paddleHitSpeedMultiplier
			state.BallVY *= paddleHitSpeedMultiplier
		} else if state.BallVX > 0 && state.BallX >= rightPaddleEdge && paddleHit(state.BallY, state.Player2YPosition) {
			state.BallX = rightPaddleEdge
			state.BallVX = -state.BallVX * paddleHitSpeedMultiplier
			state.BallVY *= paddleHitSpeedMultiplier
		}

		updated := s.repo.Update(g)
		s.broadcaster.Broadcast(room, BallMsg{Type: "ball", X: updated.State.BallX, Y: updated.State.BallY, VX: updated.State.BallVX, VY: updated.State.BallVY})
	}
}

func paddleHit(ballY, paddleY float32) bool {
	return ballY >= paddleY && ballY <= paddleY+paddleHeight
}

func resetBall(state *domain.GameData) {
	state.BallX = courtWidth / 2
	state.BallY = courtHeight / 2

	state.BallVX = ballInitialVX
	if rand.IntN(2) == 0 {
		state.BallVX = -state.BallVX
	}

	state.BallVY = ballInitialVY
	if rand.IntN(2) == 0 {
		state.BallVY = -state.BallVY
	}
}
