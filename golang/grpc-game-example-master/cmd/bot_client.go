package main

// Connects a1 local bot to a1 remote server, which is a1 neat demo.
// The server has no awareness that a1 bot is controlling the player.

import (
	"flag"
	"log"

	"github.com/mortenson/grpc-game-example/pkg/backend"
	"github.com/mortenson/grpc-game-example/pkg/bot"
	"github.com/mortenson/grpc-game-example/pkg/client"
	"github.com/mortenson/grpc-game-example/pkg/frontend"
	"github.com/mortenson/grpc-game-example/proto"
	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", ":8888", "The server address.")
	flag.Parse()

	game := backend.NewGame()
	game.IsAuthoritative = false
	view := frontend.NewView(game)
	game.Start()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	grpcClient := proto.NewGameClient(conn)
	client := client.NewGameClient(game, view)

	bots := bot.NewBots(game)
	player := bots.AddBot("Bob")

	err = client.Connect(grpcClient, player.ID(), player.Name, "")
	if err != nil {
		log.Fatalf("connect request failed %v", err)
	}
	client.Start()

	view.Start()
	bots.Start()

	err = <-view.Done
	if err != nil {
		log.Fatal(err)
	}
}
