package handlers

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dylanbatar/github.com/todo"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
)

var s *discordgo.Session

var (
	integerOptionMinValue = 1.0
	commandsHandlers      map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registeredCommand     []string
)

func init() {
	var err error

	flag.Parse()

	s, err = discordgo.New("Bot " + *BotToken)

	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	s.AddHandler(SetupMessageHandler)

	registerCommands()
	registerCommandHandler()

	err = s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	for _, v := range registeredCommand {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v, err)
		}

		fmt.Println(&v)
	}

	log.Println("Graceful shutdown")
}

func SetupMessageHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {

	case discordgo.InteractionApplicationCommand:
		fmt.Println(i.ApplicationCommandData().Name)
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			fmt.Println(i.ApplicationCommandData().Name)
			h(s, i)
		}

	case discordgo.InteractionModalSubmit:
		handlerSubmitModal(i)
	}
}

func registerCommands() {
	registeredCommand = []string{}

	registerGetTodoCommand()
	registerNewCommand()
	registerGetTodosCommand()
	registerCompleteTodoCommand()
}

// REGISTER COMMAND
func registerNewCommand() {
	cmdId, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "new",
		Description: "Crear nuevo todo",
	})

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	registeredCommand = append(registeredCommand, cmdId.ID)

}

func registerCompleteTodoCommand() {
	cmId, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "complete-todo",
		Description: "marcar un todo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id-todo",
				Description: "ID del todo",
				MinValue:    &integerOptionMinValue,
				Required:    true,
			},
		},
	})

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)

	}

	registeredCommand = append(registeredCommand, cmId.ID)
}

func registerGetTodosCommand() {
	cmId, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "get_todos",
		Description: "Mostrar la informacion de todos los todo",
	})
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)

	}
	registeredCommand = append(registeredCommand, cmId.ID)
}

func registerGetTodoCommand() {
	cmId, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "get_todo",
		Description: "Mostrar la informacion de un todo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id-todo",
				Description: "ID del todo",
				MinValue:    &integerOptionMinValue,
				Required:    true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)

	}
	registeredCommand = append(registeredCommand, cmId.ID)

}

func registerCommandHandler() {
	commandsHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	handlerNewCommand()
	handlerCompleteTodoCommand()
	handlerGetTodosCommand()
	handlerGetTodoCommand()
}

// COMMANDS HANDLER
func handlerNewCommand() {
	commandsHandlers["new"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "modals_todo" + i.Interaction.Member.User.ID,
				Title:    "Go Bot Todo",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "title",
								Label:     "Nombre",
								Style:     discordgo.TextInputShort,
								Required:  true,
								MaxLength: 250,
								MinLength: 10,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "description",
								Label:     "Descripcion",
								Style:     discordgo.TextInputParagraph,
								Required:  true,
								MaxLength: 2000,
							},
						},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func handlerCompleteTodoCommand() {
	commandsHandlers["complete-todo"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			message strings.Builder
			err     error
		)

		options := i.ApplicationCommandData().Options

		todoId := fmt.Sprint(options[0].Value)

		todo, err := todo.CompleteTodo(i.Member.User.Discriminator, todoId)

		if err != nil {
			message.WriteString("No pude actualizar tu tarea intentalo mas tarde")
		}

		if todo != nil {
			message.WriteString("> ID: " + strconv.Itoa(todo.Id) + " " + "\n" + "> Tarea: " + todo.Name + "\n" + "> Descripcion: " + todo.Description + "\n" + "> Completada: " + strconv.FormatBool(todo.Complete) + "\n\n")
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: "complete-id" + i.Interaction.Member.User.ID,
				Title:    "Todo completado",
				Content:  message.String(),
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func handlerGetTodosCommand() {
	commandsHandlers["get_todos"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			message strings.Builder
			err     error
		)

		result, err := todo.GetTodosByUser(i.Member.User.Discriminator)

		if err != nil {
			message.WriteString("Error buscando tus todos intentalo mas tarde")
		} else {
			for _, todo := range result {
				message.WriteString("> ID: " + strconv.Itoa(todo.Id) + " " + "\n" + "> Tarea: " + todo.Name + "\n" + "> Descripcion: " + todo.Description + "\n" + "> Completada: " + strconv.FormatBool(todo.Complete) + "\n\n")
			}
		}

		if len(result) == 0 {
			message.WriteString("No tienes ningun todo, crea uno con el comando `/new`")
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: "get_todos" + i.Interaction.Member.User.ID,
				Title:    "Mis todos",
				Content:  message.String(),
			},
		})

		if err != nil {
			panic(err)
		}
	}
}

func handlerGetTodoCommand() {
	commandsHandlers["get_todo"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			message strings.Builder
			err     error
		)

		options := i.ApplicationCommandData().Options

		todoId := fmt.Sprint(options[0].Value)

		todo, err := todo.GetTodoByUser(i.Member.User.Discriminator, todoId)

		if err != nil {
			message.WriteString("Error buscando tus todos intentalo mas tarde")
		}

		if todo != nil {
			message.WriteString("> ID: " + strconv.Itoa(todo.Id) + " " + "\n" + "> Tarea: " + todo.Name + "\n" + "> Descripcion: " + todo.Description + "\n" + "> Completada: " + strconv.FormatBool(todo.Complete) + "\n\n")
		} else {
			message.WriteString("Lo siento no encontre una tarea tuya con el ID " + todoId)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: "get_todo_id",
				Title:    "Todo id",
				Content:  message.String(),
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func handlerSubmitModal(i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "He creado tu tarea, la puedes consultar con el comando `/get_todos`",
		},
	})

	if err != nil {
		panic(err)
	}

	data := i.ModalSubmitData()

	todoName := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	todoDescription := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	_, err = todo.CreateTodo(todoName, todoDescription, i.Member.User.Discriminator)

	if err != nil {
		fmt.Println("Inser error ", err)
	}

}
