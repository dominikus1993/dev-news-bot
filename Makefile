output_file = dev-news-bot
webhookToken = ${DISCORD_WEBHOOK_TOKEN} 
webhookId = ${DISCORD_WEBHOOK_ID}
mongoConnection = ${MONGO_CONNECTION}
teamsWebhookUrl = ${TEAMS_WEBHOOK_URL}
build:
	go build -o $(output_file)

test:
	go test ./...

vet: 
	go vet ./...
	
buildandtest: build test

run:
	go run .

upgrade:
	go get -u
	
runbin: build
	./$(output_file) --discord-webhook-token "$(strip $(webhookToken))" --dicord-webhook-id "$(strip $(webhookId))" --mongo-connection-string "$(strip $(mongoConnection))"  -teams-webhook-url "$(strip $(teamsWebhookUrl))"