package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

// concurrent goroutine for printing observed command events by slack clients
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Sprintln("Timestamp", event.Timestamp)
		fmt.Sprintln("Command", event.Command)
		fmt.Sprintln("Parameters", event.Parameters)
		fmt.Sprintln("Event", event.Event)
		fmt.Sprintln()
	}

}
func main() {
	// load .env file
	godotenv.Load(".env")

	// create new slack client (slack bot)
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// concurrent gooutine to print bot command events
	go printCommandEvents(bot.CommandEvents())

	// create new context, defer cancel after main process exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// slack bot listens for events
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
