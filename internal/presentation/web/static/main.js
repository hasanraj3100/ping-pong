// Start with default options
import kaplay from "https://unpkg.com/kaplay@3001.0.19/dist/kaplay.mjs";

const SCREEN_W = 500;
const SCREEN_H = 300;
const PADDLE_W = 15;
const PADDLE_H = 80;
const PADDLE_SPEED = 150;

window.p2TargetY = (SCREEN_H / 2) - (PADDLE_H / 2);

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


  onKeyDown("up", () => {
    player1.move(0, -PADDLE_SPEED);
    if (player1.pos.y < 0) {
      player1.pos.y = 0;
    }
  });

  onKeyDown("down", () => {
    player1.move(0, PADDLE_SPEED);
    if (player1.pos.y + PADDLE_H > SCREEN_H) {
      player1.pos.y = SCREEN_H - PADDLE_H;
    }
  });

  player2.onUpdate(() => {
    const clampedTarget = Math.max(0, Math.min(SCREEN_H - PADDLE_H, window.p2TargetY));
    player2.pos.y = clampedTarget;
  });

})

go("main")
