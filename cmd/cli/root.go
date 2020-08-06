package cli

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/bejaneps/acroplia/internal/crud"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Config is a struct that stores all config info
type Config struct {
	Debug            bool
	PathToLogFile    string
	PathToOutputFile string

	Email    string
	Password string
	Username string
	Phone    string

	SeleniumPort    string
	SeleniumBrowser string
	SeleniumOptions []string

	TextpadTitle    string
	TextpadSubtitle string

	MessageUsername string
	MessageFullname string
	MessageText     string
	MessageChatUUID string
}

// conf values are set in command preRun
var conf = &Config{}

func init() {
	cmdRoot.PersistentFlags().StringVarP(&conf.PathToLogFile, "log", "l", "stdout", "store logs in specified file")
	viper.BindPFlag("misc.log", cmdRoot.PersistentFlags().Lookup("log"))

	cmdRoot.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false, "print debug info")
	viper.BindPFlag("misc.debug", cmdRoot.PersistentFlags().Lookup("debug"))

	cmdRoot.PersistentFlags().StringVarP(&conf.PathToOutputFile, "output", "o", "stdout", "store output of requests in specified file")
	viper.BindPFlag("misc.output", cmdRoot.PersistentFlags().Lookup("output"))

	cmdLogin.PersistentFlags().StringVar(&conf.Email, "email", "", "email for login")
	viper.BindPFlag("credentials.email", cmdLogin.PersistentFlags().Lookup("email"))

	cmdLogin.PersistentFlags().StringVar(&conf.Password, "password", "", "password for login (required for email, username or phone)")
	viper.BindPFlag("credentials.password", cmdLogin.PersistentFlags().Lookup("password"))

	cmdLogin.PersistentFlags().StringVar(&conf.Username, "username", "", "username for login")
	viper.BindPFlag("credentials.username", cmdLogin.PersistentFlags().Lookup("username"))

	cmdLogin.PersistentFlags().StringVar(&conf.Phone, "phone", "", "phone for login")
	viper.BindPFlag("credentials.phone", cmdLogin.PersistentFlags().Lookup("phone"))

	loginWeb.Flags().StringVar(&conf.SeleniumPort, "selenium-port", "4444", "port on which Selenium server is listening")
	viper.BindPFlag("selenium.port", loginWeb.Flags().Lookup("selenium-port"))

	loginWeb.Flags().StringVar(&conf.SeleniumBrowser, "selenium-browser", "firefox", "browser to be used by Selenium")
	viper.BindPFlag("selenium.browser", loginWeb.Flags().Lookup("selenium-browser"))

	loginWeb.Flags().StringSliceVar(&conf.SeleniumOptions, "selenium-options", []string{}, "additional options for selenium browser")
	viper.BindPFlag("selenium.options", loginWeb.Flags().Lookup("selenium-options"))

	cmdTextpad.PersistentFlags().StringVar(&conf.TextpadTitle, "title", "", "title for textpad")
	viper.BindPFlag("textpad.title", cmdTextpad.PersistentFlags().Lookup("title"))

	cmdTextpad.PersistentFlags().StringVar(&conf.TextpadSubtitle, "subtitle", "", "subtitle for textpad (optional)")
	viper.BindPFlag("textpad.subtitle", cmdTextpad.PersistentFlags().Lookup("subtitle"))

	textpadWeb.Flags().StringVar(&conf.Email, "email", "", "email for login")
	viper.BindPFlag("credentials.email", textpadWeb.Flags().Lookup("email"))

	textpadWeb.Flags().StringVar(&conf.Password, "password", "", "password for login")
	viper.BindPFlag("credentials.password", textpadWeb.Flags().Lookup("password"))

	textpadWeb.Flags().StringVar(&conf.Phone, "phone", "", "phone for login")
	viper.BindPFlag("credentials.phone", textpadWeb.Flags().Lookup("phone"))

	cmdMessage.PersistentFlags().StringVar(&conf.MessageFullname, "fullname", "", "full name of a user, ex: Ekaterina Gorbunova")
	viper.BindPFlag("message.fullname", cmdMessage.PersistentFlags().Lookup("fullname"))

	cmdMessage.PersistentFlags().StringVar(&conf.MessageUsername, "username", "", "username of a user, ex: ekaterina")
	viper.BindPFlag("message.username", cmdMessage.PersistentFlags().Lookup("username"))

	cmdMessage.PersistentFlags().StringVar(&conf.MessageText, "text", "Hello, World !", "message to send to a user")
	viper.BindPFlag("message.text", cmdMessage.PersistentFlags().Lookup("text"))

	messageAPI.Flags().StringVar(&conf.MessageChatUUID, "chat-uuid", "23e7235e-2072-3330-91b7-3b747d5a87e7", "chat uuid to send message to")
	viper.BindPFlag("message.chat-uuid", messageAPI.Flags().Lookup("chat-uuid"))

	messageWeb.Flags().StringVar(&conf.Email, "email", "", "email for login")
	viper.BindPFlag("credentials.email", messageWeb.Flags().Lookup("email"))

	messageWeb.Flags().StringVar(&conf.Password, "password", "", "password for login")
	viper.BindPFlag("credentials.password", messageWeb.Flags().Lookup("password"))

	messageWeb.Flags().StringVar(&conf.Phone, "phone", "", "phone for login")
	viper.BindPFlag("credentials.phone", messageWeb.Flags().Lookup("phone"))
}

var (
	logFile, outputFile *os.File
)

// cmdRoot is root command for cli tool
var cmdRoot = &cobra.Command{
	Use:   "acroplia",
	Short: `A cli tool to interact with Acroplia.`,
	Long: `A cli tool to interact with Acroplia.

By default, this tool will search unspecified flags in config.toml file located in config folder. If it can't find required flags, tool won't run.

Simply running this command won't give any results, you have to use it's subcommands: login.

Example: 
  ./acroplia login api --email my@email.com --password myPassword
  ./acroplia login web --email my@email.com --password myPassword
  ./acroplia textpad api --title MyTitle --subtitle MySubtitle
  ./acroplia textpad web --email my@email.com --password myPassword --title MyTitle --subtitle MySubtitle
  ./acroplia message api --fullname="Bezhan Mukhidinov" --text="Hello, World !"
  ./acroplia message web --email my@email.com --password myPassword --fullname="Bezhan Mukhidinov" --text="Hello, World !"
`,
	Version: "0.1alpha",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func buildPreRun(cmd *cobra.Command, args []string) error {
	var err error

	// set config values
	conf.Email = viper.GetString("credentials.email")
	conf.Password = viper.GetString("credentials.password")
	conf.Phone = viper.GetString("credentials.phone")
	conf.Username = viper.GetString("credentials.username")
	conf.Debug = viper.GetBool("misc.debug")
	conf.PathToLogFile = viper.GetString("misc.log")
	conf.PathToOutputFile = viper.GetString("misc.output")
	conf.SeleniumPort = viper.GetString("selenium.port")
	conf.SeleniumBrowser = viper.GetString("selenium.browser")
	conf.SeleniumOptions = viper.GetStringSlice("selenium.options")
	conf.TextpadTitle = viper.GetString("textpad.title")
	conf.TextpadSubtitle = viper.GetString("textpad.subtitle")
	conf.MessageFullname = viper.GetString("message.fullname")
	conf.MessageUsername = viper.GetString("message.username")
	conf.MessageChatUUID = viper.GetString("message.chat-uuid")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen}).With().Logger()

	// check if user wants more info about program steps
	if conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// check if user wants to store logs somewhere else,
	// close file afterwards in RunE
	if conf.PathToLogFile != "stdout" {
		log.Debug().Msgf("creating new %s file for logging", conf.PathToLogFile)
		logFile, err = os.Create(conf.PathToLogFile)
		if err != nil {
			log.Warn().Msgf("couldn't create log file: %v", err)
			log.Warn().Msg("using stdout for logging")

			logFile = os.Stdout
		}

		// use json formatting in logfile
		if strings.HasSuffix(conf.PathToLogFile, ".json") {
			log.Logger = log.Output(logFile).With().Logger()
		} else { // use console formating in logfile
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: logFile, NoColor: true, TimeFormat: time.RFC3339}).With().Caller().Logger()
		}
	} else {
		logFile = os.Stdout
	}

	// check if user wants to store output somewhere else
	if conf.PathToOutputFile != "stdout" {
		// check if output file exists, if no then create it
		_, err = os.Stat(conf.PathToOutputFile)
		if errors.Is(err, os.ErrNotExist) {
			log.Logger.Debug().Msg("creating new file for output")
			outputFile, err = os.Create(conf.PathToOutputFile)
			if err != nil {
				log.Logger.Warn().Msgf("couldn't create output file: %v\n", err)
				log.Logger.Warn().Msg("will use stdout as output")

				outputFile = os.Stdout
			}
		} else { // else open existing file
			log.Logger.Debug().Msg("opening existing file for output")
			outputFile, err = os.OpenFile(conf.PathToOutputFile, os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.Logger.Warn().Msgf("couldn't open output file: %v\n", err)
				log.Logger.Warn().Msg("will use stdout as output")

				outputFile = os.Stdout
			}
		}
	} else {
		outputFile = os.Stdout
	}

	// check user's internet connection
	log.Logger.Debug().Msg("checking internet connection...")
	err = crud.CheckInternetConnection()
	if err != nil {
		logFile.Close()    // in case if don't get till RunE
		outputFile.Close() // in case ""

		log.Logger.Debug().Msgf("internet connection error: %v", err)

		return errors.New("check your internet connection")
	}
	log.Logger.Debug().Msg("checking internet connection is done successfully")

	return nil
}

// Execute executes root command
func Execute() error {
	cmdLogin.AddCommand(loginAPI)
	cmdLogin.AddCommand(loginWeb)

	cmdTextpad.AddCommand(textpadAPI)
	cmdTextpad.AddCommand(textpadWeb)

	cmdMessage.AddCommand(messageWeb)
	cmdMessage.AddCommand(messageAPI)

	cmdRoot.AddCommand(cmdLogin)
	cmdRoot.AddCommand(cmdTextpad)
	cmdRoot.AddCommand(cmdMessage)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			return err
		}
	}

	return cmdRoot.Execute()
}
