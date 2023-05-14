package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	telegraph "github.com/anonyindian/telegraph-go"
	discord "github.com/bwmarrin/discordgo"
)

// Определение параметров запуска и их получение
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed — bot registers commands globally")
	BotToken       = flag.String("token", "", "Authorization token to the Discord API for bots")
	RemoveCommands = flag.Bool("clean", true, "Remove all commands after shutdown or not")
)

func init() {
	flag.Parse()
}

// Создание и подключение сессии к ППИ Дискорда
var session *discord.Session

func init() {
	var creationError error
	session, creationError = discord.New("Bot " + *BotToken)
	if creationError != nil {
		log.Fatalf("The token is invalid: %v", creationError)
	}
}

// commands — упорядоченный список комманд класса «ApplicationCommands»,
var commands = []*discord.ApplicationCommand{
	{ // Команда регистрации в Телеграфе
		Name: "register",
		NameLocalizations: &map[discord.Locale]string{
			discord.Russian: "зарегистрировать",
		},
		Description: "Register an account",
		DescriptionLocalizations: &map[discord.Locale]string{
			discord.Russian: "Зарегестрировать аккаунт на Телеграфе",
		},
		Options: []*discord.ApplicationCommandOption{
			{
				Required: true,
				Type:     discord.ApplicationCommandOptionString,
				Name:     "login",
				NameLocalizations: map[discord.Locale]string{
					discord.Russian: "логин",
				},
				Description: "Login name to be logged as into your account",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.Russian: "Регистрационное имя пользователя для создания «автора» и публикаций",
				},
			},
		},
	},
}

// commandHandlers — словарь «название команды — функция (колбэк) используящая сессию и создание интеракции»
var commandHandlers = map[string]func(session *discord.Session, interaction *discord.InteractionCreate){
	"register": func(session *discord.Session, interaction *discord.InteractionCreate) {
		shortname := interaction.ApplicationCommandData().Options[0].StringValue()
		account, _ := telegraph.CreateAccount(shortname, nil)
		switch interaction.Locale {
		default:
			session.InteractionRespond(interaction.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Embeds: []*discord.MessageEmbed{
						{
							Title:       "User \"" + account.ShortName + "\" successfully created",
							Description: "The token is `" + account.AccessToken + "`.",
						},
					},
					Flags: discord.MessageFlagsEphemeral,
					Components: []discord.MessageComponent{
						discord.ActionsRow{
							Components: []discord.MessageComponent{
								discord.Button{
									Style: discord.LinkButton,
									URL:   account.AuthUrl,
									Label: "Write an article in the web-editor",
								},
							},
						},
					},
				},
			})
		case discord.Russian:
			session.InteractionRespond(interaction.Interaction, &discord.InteractionResponse{
				Type: discord.InteractionResponseChannelMessageWithSource,
				Data: &discord.InteractionResponseData{
					Embeds: []*discord.MessageEmbed{
						{
							Title:       "Пользователь " + account.ShortName + " создан",
							Description: "Токен доступа — `" + account.AccessToken + "`.",
						},
					},
					Flags: discord.MessageFlagsEphemeral,
					Components: []discord.MessageComponent{
						discord.ActionsRow{
							Components: []discord.MessageComponent{
								discord.Button{
									Style: discord.LinkButton,
									URL:   account.AuthUrl,
									Label: "Перейти в веб-редактор",
								},
							},
						},
					},
				},
			})
		}

	},
}

// Со строки 77 идёт регистрация комманд, после каждой успешной регистрации
// название команды привязывается к её колбэку.
func init() {
	session.AddHandler(func(session *discord.Session, commandCreation *discord.InteractionCreate) {
		if callback, ok := commandHandlers[commandCreation.ApplicationCommandData().Name]; ok {
			callback(session, commandCreation)
		}
	})
}

func main() {
	// Ожидание запуска бота и лог результатов запуска
	session.AddHandler(func(session *discord.Session, ready *discord.Ready) {
		log.Printf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)
	})

	// Запуск бота
	connectionError := session.Open()
	if connectionError != nil {
		log.Fatalf("Cannot open the session: %v", connectionError)
	}

	// Регистрация интеракций
	log.Println("Adding commands...")
	registeredCommands := make([]*discord.ApplicationCommand, len(commands))
	for count, command := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *GuildID, command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}
		registeredCommands[count] = cmd
	}

	updateError := session.UpdateStatusComplex(discord.UpdateStatusData{Status: "dsa", AFK: false, IdleSince: nil})
	if updateError != nil {
		log.Panic("Пиздец")
	} else {
		log.Println("dsadas")
	}

	// Начало канала, ожидающего сигналы системы равного одному, т.е. завершения.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	// Алгорит далее начнётся только после того как «stop» будет реализован.
	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	// По завершению всего остального бот выключится (defer)
	log.Println("Gracefully shutting down.")
	defer session.Close()
}
