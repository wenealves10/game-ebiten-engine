package main

import (
	"image/color"
	"log"

	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	scale        = 30.0 // 30 pixels por metro
)

type Game struct {
	world  *box2d.B2World
	player *box2d.B2Body
	ground *box2d.B2Body
}

// NewGame inicializa o mundo físico, o chão e o jogador.
func NewGame() *Game {
	// Cria o mundo com gravidade (no Box2D, eixo Y cresce para cima, então gravidade negativa faz cair)
	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	world := box2d.MakeB2World(gravity)

	// --- Criação do chão ---
	// Definindo o corpo estático para o chão.
	groundBodyDef := box2d.MakeB2BodyDef()
	// Posicionamos o corpo no centro horizontal e em y = -5 (metros)
	groundBodyDef.Position.Set(0.0, -5.0)
	ground := world.CreateBody(&groundBodyDef)

	// Usaremos um retângulo fino para representar uma “linha” sólida.
	// A largura será a largura da tela convertida para metros (800/scale) e a altura, por exemplo, 10 pixels (10/scale).
	halfWidth := (800.0 / scale) / 2.0 // metade da largura em metros
	halfHeight := (10.0 / scale) / 2.0 // metade da altura em metros
	groundBox := box2d.MakeB2PolygonShape()
	groundBox.SetAsBox(halfWidth, halfHeight)
	// Fixture do chão (densidade zero em corpos estáticos)
	ground.CreateFixture(&groundBox, 0.0)

	// --- Criação do jogador ---
	// Corpo dinâmico
	playerBodyDef := box2d.MakeB2BodyDef()
	playerBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	// Posiciona o jogador em (0, 0) (centro do mundo)
	playerBodyDef.Position.Set(0.0, 0.0)
	player := world.CreateBody(&playerBodyDef)

	// Cria uma caixa para o jogador com 40x40 pixels (convertendo para metros)
	halfPlayer := (40.0 / scale) / 2.0 // metade do tamanho em metros
	playerBox := box2d.MakeB2PolygonShape()
	playerBox.SetAsBox(halfPlayer, halfPlayer)

	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &playerBox
	fixtureDef.Density = 1.0
	fixtureDef.Friction = 0.3
	player.CreateFixtureFromDef(&fixtureDef)

	return &Game{
		world:  &world,
		player: player,
		ground: ground,
	}
}

// worldToScreen converte as coordenadas do Box2D (metros) para pixels na tela.
// Consideramos (0,0) do mundo como o centro da tela.
func worldToScreen(pos box2d.B2Vec2) (x, y float64) {
	x = float64(screenWidth)/2.0 + pos.X*scale
	y = float64(screenHeight)/2.0 - pos.Y*scale
	return
}

func (g *Game) Update() error {
	// Avança a simulação física
	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3
	g.world.Step(timeStep, velocityIterations, positionIterations)

	// Movimentação lateral do jogador
	vel := g.player.GetLinearVelocity()
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.SetLinearVelocity(box2d.B2Vec2{X: -5.0, Y: vel.Y})
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.SetLinearVelocity(box2d.B2Vec2{X: 5.0, Y: vel.Y})
	}

	// Pulo: condição simples – se a velocidade vertical estiver quase zero, consideramos que está no chão.
	if ebiten.IsKeyPressed(ebiten.KeySpace) && (vel.Y > -0.1 && vel.Y < 0.1) {
		// Para pular, damos uma velocidade vertical positiva (o que faz subir, já que a gravidade é negativa)
		g.player.SetLinearVelocity(box2d.B2Vec2{X: vel.X, Y: 5.0})
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fundo
	screen.Fill(color.RGBA{30, 30, 30, 255})

	// --- Desenha o chão ---
	// O corpo do chão está em (0, -5) e o retângulo tem 2*halfWidth x 2*halfHeight metros.
	groundPos := g.ground.GetPosition()
	halfWidth := (800.0 / scale) / 2.0
	halfHeight := (10.0 / scale) / 2.0

	// Calcula a posição do canto superior esquerdo do chão
	left := groundPos.X - halfWidth
	top := groundPos.Y + halfHeight
	sx, sy := worldToScreen(box2d.B2Vec2{X: left, Y: top})
	widthPixels := 2 * halfWidth * scale
	heightPixels := 2 * halfHeight * scale
	ebitenutil.DrawRect(screen, sx, sy, widthPixels, heightPixels, color.RGBA{0, 255, 0, 255})

	// --- Desenha o jogador ---
	playerPos := g.player.GetPosition()
	px, py := worldToScreen(playerPos)
	// O jogador tem 40x40 pixels
	ebitenutil.DrawRect(screen, px-20, py-20, 40, 40, color.RGBA{255, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Exemplo: Ebiten + Box2D")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
