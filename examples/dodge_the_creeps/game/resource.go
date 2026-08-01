package game

import (
	"etic"
	"etic/examples/dodge_the_creeps/assets"
)

const (
	playerSprWalk1 etic.SpriteID = "player.walk.1"
	playerSprWalk2 etic.SpriteID = "player.walk.2"
	playerSprUp1   etic.SpriteID = "player.up.1"
	playerSprUp2   etic.SpriteID = "player.up.2"
	enemyFly1      etic.SpriteID = "enemy.fly.1"
	enemyFly2      etic.SpriteID = "enemy.fly.2"
	enemySwim1     etic.SpriteID = "enemy.swim.1"
	enemySwim2     etic.SpriteID = "enemy.swim.2"
	enemyWalk1     etic.SpriteID = "enemy.walk.1"
	enemyWalk2     etic.SpriteID = "enemy.walk.2"
)

func LoadResources(t *etic.Console) {
	t.LoadSprite(playerSprWalk1, etic.LoadImage(assets.PlayerGrey_walk1_png))
	t.LoadSprite(playerSprWalk2, etic.LoadImage(assets.PlayerGrey_walk2_png))
	t.LoadSprite(playerSprUp1, etic.LoadImage(assets.PlayerGrey_up1_png))
	t.LoadSprite(playerSprUp2, etic.LoadImage(assets.PlayerGrey_up2_png))

	t.LoadSprite(enemyFly1, etic.LoadImage(assets.EnemyFlyingAlt_1_png))
	t.LoadSprite(enemyFly2, etic.LoadImage(assets.EnemyFlyingAlt_2_png))
	t.LoadSprite(enemySwim1, etic.LoadImage(assets.EnemySwimming_1_png))
	t.LoadSprite(enemySwim2, etic.LoadImage(assets.EnemySwimming_2_png))
	t.LoadSprite(enemyWalk1, etic.LoadImage(assets.EnemyWalking_1_png))
	t.LoadSprite(enemyWalk2, etic.LoadImage(assets.EnemyWalking_2_png))
}
