package main

import (
	"bytes"
	"io"
	"strings"

	serversmtp "github.com/emersion/go-smtp"
	clientsmtp "net/smtp"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log                 *zap.SugaredLogger
	destination_backups []string
	destination_primary string
)

type Session struct {
	from string
	to   []string
	data bytes.Buffer
}

// Discard currently processed message.
func (s *Session) Reset() {}

// Free all resources associated with session.
func (s *Session) Logout() error {
	return s.relayMessageWithError()
}

// Set return path for currently processed message.
func (s *Session) Mail(from string, _ *serversmtp.MailOptions) error {
	s.from = from
	return nil
}

// Add recipient for currently processed message.
func (s *Session) Rcpt(to string, _ *serversmtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

// Set currently processed message contents and send it.
// r must be consumed before Data returns.
func (s *Session) Data(r io.Reader) error {
	s.data.Reset()
	_, err := s.data.ReadFrom(r)
	return err
}

func (s *Session) relayMessageWithError() error {
	from := s.from
	to := s.to
	msg := s.data.Bytes()
	// Premier envoi (serveur principal)
	primary := destination_primary
	host := strings.Split(primary, ":")[0]
	auth := clientsmtp.PlainAuth("", "", "", host)

	log.Infof("Relaying to primary %s", primary)
	err := clientsmtp.SendMail(primary, auth, from, to, msg)
	if err != nil {
		log.Errorf("Primary relay failed (%s): %v", primary, err)
		return err
	}

	log.Infof("Primary relay to %s succeeded", primary)

	// Envoi secondaire (asynchrone)
	if len(destination_backups) > 0 {
		go func() {
			for _, dest := range destination_backups {
				host := dest
				if i := strings.Index(dest, ":"); i != -1 {
					host = dest[:i]
				}

				auth := clientsmtp.PlainAuth("", "", "", host) // no auth used here
				err := clientsmtp.SendMail(dest, auth, from, to, msg)
				if err != nil {
					log.Errorf("Failed to send to %s: %v", dest, err)
				} else {
					log.Infof("Message relayed to %s", dest)
				}
			}
		}()
	}
	return nil
}

// Backend implements SMTP server backend.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *serversmtp.Conn) (serversmtp.Session, error) {
	var s serversmtp.Session = &Session{}
	return s, nil
}

func initLogger() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/smtp-dup/")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	destination_backups = viper.GetStringSlice("relay.destination_backups")
	destination_primary = viper.GetString("relay.destination_primary")
}

func main() {
	initLogger()
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("zap log sync error %s", err)
		}
	}(log)

	loadConfig()

	backend := &Backend{}
	s := serversmtp.NewServer(backend)

	s.Addr = viper.GetString("smtp.listen")
	s.Domain = viper.GetString("smtp.domain")
	s.AllowInsecureAuth = true
	//s.AuthDisabled = true
	s.EnableSMTPUTF8 = true
	s.MaxMessageBytes = 10 << 20
	s.MaxRecipients = 100

	log.Infof("SMTP duplicator starting on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
