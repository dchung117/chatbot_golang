package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

func HandleSlack(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	// get message from slack client
	query := request.Param("message")

	// parse question into JSON from witai
	question := ParseWitAI(query)

	// ask wolfram alpha
	answer := AskWolfram(question)

	// send answer back to slack client
	response.Reply(answer)
}

func ParseWitAI(query string) string {
	// create witai client
	client := witai.NewClient(os.Getenv("WITAI_SERVER_ACCESS_TOKEN"))

	// parse question from message using witai
	msgResponse, _ := client.Parse(&witai.MessageRequest{
		Query: query,
	})

	// marshal and re-format the JSON query w/ added whitespace (no prefix, indent by 4)
	data, _ := json.MarshalIndent(msgResponse, "", "    ")

	// convert JSON byte-object to string
	rough := string(data[:])

	// get the question (in field "value") from the JSON-string
	value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")

	return value.String()
}

func AskWolfram(question string) string {
	// create wolfram client
	client := &wolfram.Client{AppID: os.Getenv("WOLFRAM_ALPHA_APP_ID")}

	// ask WolframAlpha
	response, err := client.GetSpokentAnswerQuery(question, wolfram.Metric, 1000)
	if err != nil {
		fmt.Println("There is an error w/ Ask Wolfram.")
	}

	return response
}
