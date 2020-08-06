package cli

import (
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bejaneps/acroplia/internal/crud"
	"github.com/bejaneps/acroplia/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tebeka/selenium"
)

const loginRespFilename = "./config/.login-response.json"

var (
	// ErrLoginFailed used when login to Acroplia fails
	ErrLoginFailed = errors.New("login failed")

	// ErrLoginSpecify used when user doesn't supply any of credentials used by login
	ErrLoginSpecify = errors.New("you have to specify one of login methods, either by: phone, email, username or magic link")

	// ErrLoginPasswordRequired used when password is not supplied by user when performing login
	ErrLoginPasswordRequired = errors.New("password is required for login by email")
)

// cmdLogin is a command to make a login request to Acroplia using API
var cmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to Acroplia",
	Long: `Login to Acroplia.
	
You have to use it's subcommands: api or web to perform actual login.

Example: 
  ./acroplia login api
  ./acroplia login web
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
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

		// save login response to temporary file, so we can use it further for creating textpad
		f, err := os.Create(loginRespFilename)
		if err != nil {
			log.Logger.Debug().Msgf("os.Create: %v", err)

			return errors.New("couldn't create temporary file for login response")
		}
		defer f.Close()

		// write both to user specified output file and temporary file to save response
		wr := io.MultiWriter(f, outputFile)

		log.Logger.Info().Msgf("writing response message to: %s\n", conf.PathToOutputFile)
		enc := json.NewEncoder(wr)
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
		time.Sleep(2 * time.Second) // sleep for 2 seconds, so user can see end result

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func performLoginWeb(wd selenium.WebDriver) error {
	var err error

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

	return nil
}
