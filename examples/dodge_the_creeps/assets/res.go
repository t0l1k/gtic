package assets

import (
	_ "embed"
	_ "image/png"
)

const XoloniumRegular = "XoloniumRegular"

var (
	//go:embed fonts/Xolonium-Regular.ttf
	XoloniumRegular_ttf []byte
)

var (
	//go:embed art/playerGrey_walk1.png
	PlayerGrey_walk1_png []byte
	//go:embed art/playerGrey_walk2.png
	PlayerGrey_walk2_png []byte
	//go:embed art/playerGrey_up1.png
	PlayerGrey_up1_png []byte
	//go:embed art/playerGrey_up2.png
	PlayerGrey_up2_png []byte
)
var (
	//go:embed art/enemyFlyingAlt_1.png
	EnemyFlyingAlt_1_png []byte
	//go:embed art/enemyFlyingAlt_2.png
	EnemyFlyingAlt_2_png []byte
	//go:embed art/enemySwimming_1.png
	EnemySwimming_1_png []byte
	//go:embed art/enemySwimming_2.png
	EnemySwimming_2_png []byte
	//go:embed art/enemyWalking_1.png
	EnemyWalking_1_png []byte
	//go:embed art/enemyWalking_2.png
	EnemyWalking_2_png []byte
)

var (
	//go:embed art/gameover.wav
	Gameover_wav []byte
	//go:embed art/HouseInAForestLoop.ogg
	HouseInAForestLoop_ogg []byte
)
