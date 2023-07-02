package main

import (
	"context"
	"log"

	"github.com/BimaAdi/WebsocketScaler/core"
	"github.com/BimaAdi/WebsocketScaler/scalergoredis"
	"github.com/BimaAdi/WebsocketScaler/wsclientfiber"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/redis/go-redis/v9"
)

type Event struct {
}

func (e *Event) OnConnect(s core.ScalerContract, socket_id string, params core.Params) {
	s.SendToSingleUser(socket_id, "Welcome")
	s.SendToAll("sommeone connected")
}

func (e *Event) OnMessage(s core.ScalerContract, socket_id string, payload string) {
	s.SendToAll(payload)
}

func (e *Event) OnDisconnect(s core.ScalerContract, socket_id string) {

}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	scl := scalergoredis.NewRedisScaler(rdb, ctx, "ws_channel")
	ws_router := wsclientfiber.NewFiberWebsocket()
	go scl.Subscribe(ws_router)
	event := Event{}
	router := ws_router.CreateWebsocketRoute(&event, scl)

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/", router)

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Render("index", fiber.Map{})
	})

	log.Fatal(app.Listen(":3000"))
}
