package types

type PixelData [][][]int
type FontData [][]int

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type SpriteType int

const (
	AnimationSprite SpriteType = iota
	TextSprite
	StaticSprite
)

type DrawOptions struct {
	SpriteType  SpriteType
	Reverse     bool
	Loop        bool
	Scroll      bool
	ScrollSpeed int
	Direction   Direction
}
