package cli

import (
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"

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

	XAuthToken string

	SeleniumPort    string
	SeleniumBrowser string
	SeleniumOptions []string
}

// conf values are set in command preRun
var conf = &Config{}

func init() {
	cmdRoot.PersistentFlags().StringVarP(&conf.PathToLogFile, "log", "l", "stdout", "store logs in specified file")
	viper.BindPFlag("misc.log", cmdRoot.PersistentFlags().Lookup("log"))

	cmdRoot.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false, "print debug info")
	viper.BindPFlag("misc.debug", cmdRoot.PersistentFlags().Lookup("debug"))

	cmdLogin.PersistentFlags().StringVar(&conf.Email, "email", "", "email for login")
	viper.BindPFlag("credentials.email", cmdLogin.PersistentFlags().Lookup("email"))

	cmdLogin.PersistentFlags().StringVar(&conf.Password, "password", "", "password for login (required for email, username or phone)")
	viper.BindPFlag("credentials.password", cmdLogin.PersistentFlags().Lookup("password"))

	cmdLogin.PersistentFlags().StringVar(&conf.Username, "username", "", "username for login")
	viper.BindPFlag("credentials.username", cmdLogin.PersistentFlags().Lookup("username"))

	cmdLogin.PersistentFlags().StringVar(&conf.Phone, "phone", "", "phone for login")
	viper.BindPFlag("credentials.phone", cmdLogin.PersistentFlags().Lookup("phone"))

	loginAPI.Flags().StringVar(&conf.XAuthToken, "x-auth-token", "x", "token for making requests (textpad, messages)")
	viper.BindPFlag("credentials.x-auth-token", loginAPI.Flags().Lookup("x-auth-token"))

	loginAPI.Flags().StringVarP(&conf.PathToOutputFile, "output", "o", "stdout", "store output of requests in specified file")
	viper.BindPFlag("misc.output", loginAPI.Flags().Lookup("output"))

	loginWeb.Flags().StringVar(&conf.SeleniumPort, "selenium-port", "4444", "port on which Selenium server is listening")
	viper.BindPFlag("selenium.port", loginWeb.Flags().Lookup("selenium-port"))

	loginWeb.Flags().StringVar(&conf.SeleniumBrowser, "selenium-browser", "firefox", "browser to be used by Selenium")
	viper.BindPFlag("selenium.browser", loginWeb.Flags().Lookup("selenium-browser"))

	loginWeb.Flags().StringSliceVar(&conf.SeleniumOptions, "selenium-options", []string{}, "additional options for selenium browser")
	viper.BindPFlag("selenium.options", loginWeb.Flags().Lookup("selenium-options"))
}

var (
	logFile, outputFile *os.File
)

var (
	// ErrLoginFailed used when login to Acroplia fails
	ErrLoginFailed = errors.New("login failed")

	// ErrLoginSpecify used when user doesn't supply any of credentials used by login
	ErrLoginSpecify = errors.New("you have to specify one of login methods, either by: phone, email, username or magic link")

	// ErrLoginPasswordRequired used when password is not supplied by user when performing login
	ErrLoginPasswordRequired = errors.New("password is required for login by email")
)

// cmdRoot is root command for cli tool
var cmdRoot = &cobra.Command{
	Use:   "acroplia",
	Short: `A cli tool to interact with Acroplia.`,
	Long: `A cli tool to interact with Acroplia.

By default, this tool will search unspecified flags in config.toml file located in config folder. If it can't find required flags, tool won't run.

Simply running this command won't give any results, you have to use it's subcommands: login.

Example: ./acroplia login api --email my@email.com --password myPassword
	 ./acroplia login web --email my@email.com --password myPassword
`,
	Version: "0.1alpha",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// cmdLogin is a command to make a login request to Acroplia using API
var cmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to Acroplia",
	Long: `Login to Acroplia.
	
You have to use it's subcommands: api or web to perform actual login.

Example: ./acroplia login api
	 ./acroplia login web
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// loginAPI ...
var loginAPI = &cobra.Command{
	Use:     "api",
	Short:   "Use Acroplia API for login",
	Long:    `Use Acroplia API for login`,
	PreRunE: buildPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		// catch os signals, like Ctrlc + C and clean up opened resources
		go func() {
			log.Logger.Debug().Msg("catching os signals")
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			log.Logger.Debug().Msg("received SIGINT(Ctrl + C) signal, cleaning up opened resources...")

			// close opened resources
			logFile.Close()
			outputFile.Close()

			// TODO: work on exit better
			os.Exit(0)
		}()

		var err error

		// check how user wants to login and perform a login
		log.Logger.Debug().Msg("performing a login ...")
		resp := &crud.ResponseAuthResponse{}
		if conf.Email != "" {
			if conf.Password == "" {
				return ErrLoginPasswordRequired
			}

			s := crud.NewSignInEmailRequest(conf.Email, conf.Password)

			resp, err = s.LoginByEmail()
			if errors.Cause(err) == crud.ErrInvalidCredentials {
				return crud.ErrInvalidCredentials
			} else if err != nil {
				log.Logger.Debug().Msgf("loginByEmail: %v", err)

				return ErrLoginFailed
			}
		} else if conf.Phone != "" {
			if conf.Password == "" {
				return ErrLoginPasswordRequired
			}

			s := crud.NewSignInPhoneRequest(conf.Phone, conf.Password)

			resp, err = s.LoginByPhone()
			if err != nil {
				log.Logger.Debug().Msgf("loginByPhone: %v", err)

				return ErrLoginFailed
			}
		} else if conf.Username != "" {
			if conf.Password == "" {
				return ErrLoginPasswordRequired
			}

			s := crud.NewSignInUsernameRequest(conf.Username, conf.Password)

			resp, err = s.LoginByUsername()
			if err != nil {
				log.Logger.Debug().Msgf("loginByUsername: %v", err)

				return ErrLoginFailed
			}
		} else {
			return ErrLoginSpecify
		}
		log.Logger.Debug().Msg("login was done successfully")

		log.Logger.Info().Msgf("writing response message to: %s\n", conf.PathToOutputFile)
		err = json.NewEncoder(outputFile).Encode(resp)
		if err != nil {
			log.Logger.Debug().Msgf("json.NewEncoder: %v", err)

			return errors.New("couldn't encode response")
		}

		return nil
	},
}

// loginWeb ...
var loginWeb = &cobra.Command{
	Use:     "web",
	Short:   "Use Acroplia Web interface (Selenium) for login",
	Long:    `Use Acroplia Web interface (Selenium) for login`,
	PreRunE: buildPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// instantiate selenium browser
		log.Logger.Debug().Msgf("opening instance of %s browser...", conf.SeleniumBrowser)
		wd, err := services.NewSelenium(conf.SeleniumPort, conf.SeleniumBrowser, conf.SeleniumOptions...)
		if err != nil {
			log.Logger.Debug().Msgf("newSelenium: %v", err)

			return errors.New("failed to start Selenium")
		}
		defer wd.Close()
		log.Logger.Debug().Msgf("%s browser opened successfully", conf.SeleniumBrowser)

		// catch os signals, like Ctrlc + C and clean up opened resources
		go func() {
			log.Logger.Debug().Msg("catching os signals")
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			log.Logger.Debug().Msg("received SIGINT(Ctrl + C) signal, cleaning up opened resources...")

			// close opened resources
			logFile.Close()
			outputFile.Close()
			wd.Close()
			wd.Quit()

			// TODO: work on exit better
			os.Exit(0)
		}()

		// check how user wants to login and perform a login
		log.Logger.Debug().Msg("performing a login ...")
		if conf.Email != "" { // login by email
			if conf.Password == "" {
				return ErrLoginPasswordRequired
			}

			err = crud.LoginByEmail(conf.Email, conf.Password, wd)
			if errors.Is(err, crud.ErrInvalidEmail) || errors.Is(err, crud.ErrInvalidPassword) {
				return err
			} else if err != nil {
				log.Logger.Debug().Msgf("LoginByEmail: %v", err)

				return ErrLoginFailed
			}
		} else if conf.Phone != "" { // login by phone
			if conf.Password == "" {
				return ErrLoginPasswordRequired
			}

			err = crud.LoginByPhone(conf.Phone, conf.Password, wd)
			if errors.Is(err, crud.ErrInvalidPhone) || errors.Is(err, crud.ErrInvalidPassword) {
				return err
			} else if err != nil {
				log.Logger.Debug().Msgf("LoginByPhone: %v", err)

				return ErrLoginFailed
			}
		} else {
			return ErrLoginSpecify
		}
		log.Logger.Debug().Msg("login was done successfully")
		time.Sleep(2 * time.Second) // sleep for 2 seconds, so user can see end result

		return nil
	},
}

func buildPreRun(cmd *cobra.Command, args []string) error {
	var err error

	// set config values
	conf.Email = viper.GetString("credentials.email")
	conf.Password = viper.GetString("credentials.password")
	conf.Phone = viper.GetString("credentials.phone")
	conf.Username = viper.GetString("credentials.username")
	conf.XAuthToken = viper.GetString("credentials.x-auth-token")
	conf.Debug = viper.GetBool("misc.debug")
	conf.PathToLogFile = viper.GetString("misc.log")
	conf.PathToOutputFile = viper.GetString("misc.output")
	conf.SeleniumPort = viper.GetString("selenium.port")
	conf.SeleniumBrowser = viper.GetString("selenium.browser")
	conf.SeleniumOptions = viper.GetStringSlice("selenium.options")

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

	cmdRoot.AddCommand(cmdLogin)

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
