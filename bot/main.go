package main

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	rand.Seed(time.Now().Unix())

	// TODO Derive the token at runtime
	gmHandler := NewGroupMeMessenger("c1d5a8381af0c8b21fe3bcb2b6", &url.URL{
		Scheme:   "https",
		Host:     "api.groupme.com",
		Path:     "v3/bots/post",
		RawQuery: "token=<Insert Real Token Value>",
	})

	annoyanceResponder := NewTriggerResponder("Filbot", "", []string{
		"хватит надоедать мне!",
		"БЫДЛО",
		"ДУРАК",
		"СВОЛОЧЬ",
		"ПОДОНОК",
		"Сволочь",
		"ДРЯНЬ",
	})

	ilyaResponder := NewRandomResponder(0.05, rand.Float64, []string{
		"When I talk English, I’m still computing.",
		"When I hear my voice, it sounds disgusting.",
		"Guys! A question! Is chicken meat? Or is it a bird?",
		"In Philly, win or lose, you’re still overpaid.",
		"Decepticons always have nice stuff." +
			" Autobots, they rebels. They weapons aren’t as nice.",
		"Whoa, whoa, whoa, whoa. I don’t have contract. I’m here as guest.",
		"Siberian Husky. She’s all white. Beautiful blue eyes." +
			" That’s basically blonde girl with blue eyes." +
			" Your dream, man. My husky, basically, she’s a hot girl, man.",
		"I think like, ‘And we have some problems here on the earth we worry" +
			" about? Compared to like…nothing. Just be happy. Don’t worry be happy" +
			" right now",
	})

	pingResponder := NewTriggerResponder("Ping", "", []string{"понг"})

	tradeResponder := NewTriggerResponder(
		"Trade", "You give best player, I give ", []string{
			"Bucket of Wood Chips",
			"Concrete Crumbles",
			"Fabergé egg",
			"Guy Fieri",
			"Putin's love",
			"2020 Presidential election",
			"Ukrainian prostitutes",
			"Pyser's kidney",
		})

	samHandler := NewSAMHandler(
		gmHandler, annoyanceResponder, ilyaResponder, pingResponder, tradeResponder)
	lambda.Start(samHandler.Handle)
}
