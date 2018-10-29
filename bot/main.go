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
		"Shut Up Jake!",
		"Fuck Off!",
		"I like cats",
		"Who's askin?",
		"Fuck your couch",
		"Just let me live!",
		"I fart in your general direction!",
		"Your mother was a hamster and your father smelt of elder berries!",
	})

	randomResponder := NewRandomResponder(0.05, rand.Float64, []string{
		"It me, Filbot!",
		"You can build a throne with bayonets, but you can't sit on it for long.",
		"The Vegas Golden Knights are poorly constructed.",
		"All hail Trollbot!",
		"Yahoo Sports was made by the son of a barmaid and a donkey",
		"Doesn't matter how good Temple's stadium will be, will never compare to Luzhniki",
		"I wish I was a puppy",
	})

	pingResponder := NewTriggerResponder("Ping", "", []string{"понг"})

	tradeResponder := NewTriggerResponder(
		"Trade", "You give best player, I give ", []string{
			"Lebron James",
			"Meir Engel",
			"Evgeni Malkin",
			"Magadan Oblast",
			"Rob Gronkowski",
			"Pool full of Borscht.",
		})

	samHandler := NewSAMHandler(
		gmHandler, annoyanceResponder, randomResponder, pingResponder, tradeResponder)
	lambda.Start(samHandler.Handle)
}
