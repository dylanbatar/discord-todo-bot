## GO TODO BOT
Hey this project was inspired by this [blog](https://dev.to/aurelievache/learning-go-by-examples-part-4-create-a-bot-for-discord-in-go-43cf) post but arrived to the next level

This bot implement a discord slash command and component options such a modal and autocomplete

## HOW IS WORK
1. The most important part, you need create discord bot application created in [discord application]("https://discord.com/developers/applications")
2. Before start the project you need set the `.env` file. You need setup these variables.
```
APPLICATION_ID=
GUILD_ID=
BOT_TOKEN=

HOST_DATABASE=
USER_DATABASE=
PASSWORD_DATABASE=
PORT_DATABASE=
```
3. This example use mysql database to work, for that you need a database running
4. Invite The discord bot to your channel using this invitation link 

      Replace the <CLIENT-ID> for you own
`https://discord.com/api/oauth2/authorize?client_id=<CLIENT-ID>&permissions=8&scope=bot`
5. Run the project `go run main.go` or compiling the app