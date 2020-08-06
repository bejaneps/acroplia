package cli

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	// ErrMessageFailed used when sending message to a user fails
	ErrMessageFailed = errors.New("sending message failed")

	// ErrMessageReceiver used when user didn't supply fullname and username
	ErrMessageReceiver = errors.New("fullname and username of receiver can't be empty")
)

var cmdMessage = &cobra.Command{
	Use:   "message",
	Short: "Message a user in Acroplia",
	Long: `Message a user in Acroplia.

You have to use it's subcommands: api or web to send actual message.

Example: 
	./acroplia message api
	./acroplia message web
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var messageWeb = &cobra.Command{
	Use:   "web",
	Short: "Use Acroplia Web interface (Selenium) to send a message",
	Long: `Use Acroplia Web interface (Selenium) to send a message.

You have to supply your email or phone and password to perform login, after that creating textpad is possible.
`,
	PreRunE: buildPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		if conf.MessageFullname == "" && conf.MessageUsername == "" {
			return ErrMessageReceiver
		}

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
			wd.Close()
			wd.Quit()

			// TODO: work on exit better
			os.Exit(0)
		}()

		// check how user wants to login and perform a login
		log.Logger.Debug().Msg("performing a login ...")
		if err := performLoginWeb(wd); err != nil {
			return err
		}
		log.Logger.Debug().Msg("login was done successfully")

		// send a message to user
		log.Logger.Debug().Msg("sending message ...")
		err = crud.SendMessage(conf.MessageFullname, conf.MessageUsername, conf.MessageText, wd)
		if errors.Is(err, crud.ErrUserNotFound) {
			return crud.ErrUserNotFound
		} else if err != nil {
			log.Logger.Debug().Msgf("crud.SendMessage: %v", err)

			return ErrMessageFailed
		}
		log.Logger.Debug().Msg("sending message was done successfully")

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var messageAPI = &cobra.Command{
	Use:   "api",
	Short: "Use Acroplia API to send a message",
	Long: `User Acroplia API to send a message.

You have to perform login before using this command as it needs user information from login.
`,
	PreRunE: buildPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		if conf.MessageFullname == "" && conf.MessageUsername == "" {
			return ErrMessageReceiver
		}

		// open login response file for decoding it
		log.Logger.Debug().Msgf("opening response file: %v ...", loginRespFilename)
		f, err := os.Open(loginRespFilename)
		if err != nil {
			log.Logger.Debug().Msgf("os.Open: %v", err)

			return errors.New("couldn't open login response file, make sure you have performed login")
		}
		log.Logger.Debug().Msgf("opening response file was done successfully")

		// decode contents of login response into struct
		log.Logger.Debug().Msg("decoding response file into struct")
		authResponse := &crud.ResponseAuthResponse{}
		err = json.NewDecoder(f).Decode(authResponse)
		if err != nil {
			log.Logger.Debug().Msgf("json.Decode: %v", err)

			return errors.New("couldn't decode login response to struct")
		}
		log.Logger.Debug().Msg("decoding response file was done successfully")

		// sending a message
		log.Logger.Debug().Msg("sending message ...")
		msg := crud.NewMessage(conf.MessageText, authResponse.Data.User)
		resp, err := msg.SendMessage(conf.MessageChatUUID, authResponse.Data.AccessToken)
		if errors.Is(err, crud.ErrMessageInternal) {
			return crud.ErrMessageInternal
		} else if err != nil {
			log.Logger.Debug().Msgf("crud.SendMessage: %v", err)

			return ErrMessageFailed
		}
		log.Logger.Debug().Msg("sending message was done successfully")

		// store response message to specified output file
		log.Logger.Info().Msgf("writing response message to: %s\n", conf.PathToOutputFile)
		enc := json.NewEncoder(outputFile)
		enc.SetIndent("", "  ")
		err = enc.Encode(resp)
		if err != nil {
			log.Logger.Debug().Msgf("json.NewEncoder: %v", err)

			return errors.New("couldn't encode response")
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
