package main

import (
	"image"
	"log"
	"math/rand"

	_ "image/png"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/gl"
)

type (
	WindowSize struct {
		width  int
		height int
	}
	ColorFloat struct {
		red   float32
		green float32
		blue  float32
	}

	Game struct {
		windowSize WindowSize
		color      ColorFloat
		Pressed    bool
	}
)

const (
	squid       = iota
	squidWidth  = 527
	squidHeight = 589
)

var (
	splatoonColor = [][]byte{
		{0xed, 0x8b, 0x49},
		{0x7d, 0x8b, 0x0b},
		{0x26, 0x84, 0x6c},
		{0x45, 0x43, 0xd8},
		{0xe6, 0x8e, 0x27},
		{0x7d, 0x8b, 0x0e},
		{0x07, 0x7f, 0x73},
		{0x51, 0x24, 0xdd},
		{0xea, 0x8c, 0x0d},
		{0x7b, 0x91, 0x0f},
		{0x14, 0x7f, 0x86},
		{0x88, 0x4b, 0xd6},
		{0xe5, 0x92, 0x1e},
		{0x6c, 0x8e, 0x18},
		{0x07, 0x7a, 0x90},
		{0xcf, 0x60, 0xd3},
		{0xeb, 0x8f, 0x0b},
		{0x5a, 0x91, 0x0b},
		{0x0e, 0x70, 0xc5},
		{0xd7, 0x30, 0xdb},
		{0xe9, 0x94, 0x0b},
		{0x56, 0x89, 0x31},
		{0x0a, 0x4d, 0xd8},
		{0xe2, 0x4f, 0xd5},
		{0xd2, 0x92, 0x0a},
		{0x24, 0x83, 0x0d},
		{0x41, 0x64, 0xcf},
		{0xe6, 0x63, 0xd0},
		{0xcc, 0x93, 0x17},
		{0x70, 0xac, 0x62},
		{0x2b, 0x4d, 0xd6},
		{0xed, 0x74, 0xc9},
		{0xce, 0x92, 0x0f},
		{0x38, 0x84, 0x2e},
		{0x57, 0x62, 0xd1},
		{0xe7, 0x77, 0xbc},
		{0xc6, 0x92, 0x0b},
		{0x77, 0xab, 0x9a},
		{0x24, 0x2e, 0xda},
		{0xee, 0x76, 0xa9},
		{0x56, 0x92, 0x82},
		// フェス
		{0xcc, 0x4f, 0x25}, // 1
		{0x3d, 0x8e, 0x75},
		{0x7c, 0x2f, 0x4f}, // 2
		{0x06, 0x64, 0x3d},
		{0xe5, 0xcc, 0x9b}, // 3
		{0xd4, 0x7d, 0x23},
		{0x6f, 0x24, 0x82}, // 4
		{0x6d, 0x80, 0x20},
		{0x0c, 0x3f, 0xae}, // 5
		{0xe8, 0x6e, 0x19},
		{0x14, 0x87, 0xb5}, // 6
		{0xa6, 0x43, 0x52},
		{0xd6, 0x54, 0x7e}, // 7
		{0xe6, 0x8f, 0x1c},
		{0xf4, 0x6b, 0x13}, // 色覚
		{0x38, 0x2c, 0xb3},
	}
)

func NewGame(glctx *gl.Context) *Game {
	var g Game
	g.start(glctx)
	return &g
}

func (g *Game) start(glctx *gl.Context) {
	g.color = ColorFloat{red: 1.0, green: 1.0, blue: 1.0}

	(*glctx).Enable(gl.BLEND)
	(*glctx).BlendEquation(gl.FUNC_ADD)
	(*glctx).BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func (g *Game) update(glctx *gl.Context, sz size.Event) {
	g.windowSize.width = int(sz.WidthPt)
	g.windowSize.height = int(sz.HeightPt)
	if true == g.Pressed {
		g.Pressed = false
		color := splatoonColor[rand.Intn(len(splatoonColor))]
		g.color = ColorFloat{
			red:   1.0 / 255.0 * float32(color[0]),
			green: 1.0 / 255.0 * float32(color[1]),
			blue:  1.0 / 255.0 * float32(color[2]),
		}
	}

	(*glctx).ClearColor(g.color.red, g.color.green, g.color.blue, 1)
	(*glctx).Clear(gl.COLOR_BUFFER_BIT)
}

func (g *Game) exit() {

}

func (g *Game) Scene(eng sprite.Engine) *sprite.Node {
	texs := loadTextures(eng)

	scene := &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	newNode := func(fn arrangerFunc) {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		scene.AppendChild(n)
	}

	newNode(func(se sprite.Engine, n *sprite.Node, t clock.Time) {
		scaleFactor := getScaleFactor(g.windowSize)
		var width float32 = float32(squidWidth) * scaleFactor
		var height float32 = float32(squidHeight) * scaleFactor

		a := f32.Affine{
			{width, 0, (float32(g.windowSize.width) - width) / 2},
			{0, height, (float32(g.windowSize.height) - height) / 2},
		}

		se.SetSubTex(n, texs[squid])
		se.SetTransform(n, a)
	})

	return scene
}

func getScaleFactor(ws WindowSize) float32 {
	if squidHeight/float32(ws.height) > squidWidth/float32(ws.width) {
		return float32(ws.height) / squidHeight
	}
	return float32(ws.width) / squidWidth
}

func loadTextures(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("splatoon_mask.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	img, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}

	return []sprite.SubTex{
		squid: sprite.SubTex{t, image.Rect(0, 0, squidWidth, squidHeight)},
	}
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
