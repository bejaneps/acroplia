package cli

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	// ErrTextpadFailed used when creating textpad in Acroplia fails
	ErrTextpadFailed = errors.New("creating textpad failed")
)

var cmdTextpad = &cobra.Command{
	Use:   "textpad",
	Short: "Create a Textpad in Acroplia",
	Long: `Create a Textpad in Acroplia.

You have to use it's subcommands: api or web to create actual textpad.

Don't forget to perform login through API, before using this command !

Example: 
  ./acroplia textpad api
  ./acroplia textpad web
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var textpadAPI = &cobra.Command{
	Use:   "api",
	Short: "Use Acroplia API to create texpad",
	Long: `Use Acroplia API to create textpad.

You have to perform login before using this command as it needs user information from login.
`,
	PreRunE: buildPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if user supplied textpad title
		if conf.TextpadTitle == "" {
			return errors.New("textpad title can't be empty for this command")
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

		// create an actual textpad
		log.Logger.Debug().Msgf("creating a textpad with title %s and subtitle %s", conf.TextpadTitle, conf.TextpadSubtitle)
		textpad := crud.NewTextpad(conf.TextpadTitle, conf.TextpadSubtitle, authResponse.Data.User)
		textpadResp, err := textpad.Create(authResponse.Data.AccessToken)
		if err != nil {
			log.Logger.Debug().Msgf("crud.NewTextpad: %v", err)

			return ErrTextpadFailed
		}
		log.Logger.Debug().Msg("creating a textpad was done successfully")

		// store response message to specified output file
		log.Logger.Info().Msgf("writing response message to: %s\n", conf.PathToOutputFile)
		enc := json.NewEncoder(outputFile)
		enc.SetIndent("", "  ")
		err = enc.Encode(textpadResp)
		if err != nil {
			log.Logger.Debug().Msgf("json.NewEncoder: %v", err)

			return errors.New("couldn't encode response")
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var textpadWeb = &cobra.Command{
	Use:  "web",
	Long: "Use Acroplia Web to create texpad",
	Short: `Use Acroplia Web to create texpad.

You have to supply your email or phone and password to perform login, after that creating textpad is possible.
	`,
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

		// create a textpad
		log.Logger.Debug().Msg("creating a textpad ...")
		err = crud.TextpadCreate(conf.TextpadTitle, conf.TextpadSubtitle, wd)
		if err != nil {
			log.Logger.Debug().Msgf("crud.TextpadCreate: %v", err)
			time.Sleep(5 * time.Minute)

			return ErrTextpadFailed
		}
		log.Logger.Debug().Msg("creating textpad was done successfully")

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
