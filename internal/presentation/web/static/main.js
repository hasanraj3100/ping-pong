// Start with default options
import kaplay from "./kaplay.mjs";

const SCREEN_W = 500;
const SCREEN_H = 300;
const PADDLE_W = 15;
const PADDLE_H = 80;
const BALL_SIZE = 10;

window.p1TargetY = (SCREEN_H / 2) - (PADDLE_H / 2);
window.p2TargetY = (SCREEN_H / 2) - (PADDLE_H / 2);
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


  player1.onUpdate(() => {
    const clampedTarget = Math.max(0, Math.min(SCREEN_H - PADDLE_H, window.p1TargetY));
    player1.pos.y = clampedTarget;
    window.p1TargetY = player1.pos.y;
  });

  player2.onUpdate(() => {
    const clampedTarget = Math.max(0, Math.min(SCREEN_H - PADDLE_H, window.p2TargetY));
    player2.pos.y = clampedTarget;
    window.p2TargetY = player2.pos.y;
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
