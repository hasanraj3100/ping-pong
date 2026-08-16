// Start with default options
import kaplay from "./kaplay.mjs";

const SCREEN_W = 500;
const SCREEN_H = 300;
const PADDLE_W = 15;
const PADDLE_H = 80;
const BALL_SIZE = 10;
const PADDLE_SPEED = 200; // px/sec, matches server paddleStep(10) applied every 50ms while a key is held

window.p1TargetY = (SCREEN_H / 2) - (PADDLE_H / 2);
window.p2TargetY = (SCREEN_H / 2) - (PADDLE_H / 2);
window.p1ForceSnap = false;
window.p2ForceSnap = false;
window.ballX = SCREEN_W / 2;
window.ballY = SCREEN_H / 2;
window.ballVX = 0;
window.ballVY = 0;

kaplay({
  width: SCREEN_W,
  height: SCREEN_H,
  background: "#202940",
  scale: 2,
  canvas: document.getElementById("game-canvas"),
});

scene("main", () => {

  const player1 = add([
    rect(PADDLE_W, PADDLE_H),
    pos(10, (SCREEN_H / 2) - (PADDLE_H / 2)),
    color("#9A8678"),
    "player1",
  ]);


  const player2 = add([
    rect(PADDLE_W, PADDLE_H),
    pos(SCREEN_W - 10 - PADDLE_W, (SCREEN_H / 2) - (PADDLE_H / 2)),
    color("#9A8678"),
    "player2",
  ])


  function updatePaddle(paddle, targetKey, snapKey) {
    const clampedTarget = Math.max(0, Math.min(SCREEN_H - PADDLE_H, window[targetKey]));
    window[targetKey] = clampedTarget;

    if (window[snapKey]) {
      // own paddle's prediction was wrong: correct it immediately
      paddle.pos.y = clampedTarget;
      window[snapKey] = false;
    } else {
      // glide toward the target - covers own predicted movement and the
      // opponent's paddle, which has no local prediction to glide from otherwise
      const maxDelta = PADDLE_SPEED * dt();
      const diff = clampedTarget - paddle.pos.y;
      paddle.pos.y += Math.abs(diff) <= maxDelta ? diff : Math.sign(diff) * maxDelta;
    }
  }

  player1.onUpdate(() => {
    updatePaddle(player1, "p1TargetY", "p1ForceSnap");
  });

  player2.onUpdate(() => {
    updatePaddle(player2, "p2TargetY", "p2ForceSnap");
  });

  const ball = add([
    rect(BALL_SIZE, BALL_SIZE),
    pos(window.ballX - BALL_SIZE / 2, window.ballY - BALL_SIZE / 2),
    color("#E4DCCF"),
    "ball",
  ]);

  ball.onUpdate(() => {
    window.ballX += window.ballVX * dt();
    window.ballY += window.ballVY * dt();

    while (window.ballY < 0 || window.ballY > SCREEN_H) {
      if (window.ballY < 0) {
        window.ballY = -window.ballY;
        window.ballVY = -window.ballVY;
      } else if (window.ballY > SCREEN_H) {
        window.ballY = 2 * SCREEN_H - window.ballY;
        window.ballVY = -window.ballVY;
      }
    }

    ball.pos.x = window.ballX - BALL_SIZE / 2;
    ball.pos.y = window.ballY - BALL_SIZE / 2;
  });

})

go("main")
