package nwws

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/xmppo/go-xmpp"
)

type Config struct {
	Server   string
	Room     string
	User     string
	Pass     string
	Resource string
}

func (conf *Config) serverName() string {
	return strings.Split(conf.Server, ":")[0]
}

type NWWS struct {
	Config *Config
	client *xmpp.Client
	last   time.Time
}

type Message struct {
	Text       string
	ReceivedAt time.Time
}

func confirmConnection(client *xmpp.Client, config *Config) error {
	for {
		chat, err := client.Recv()
		if err != nil {
			return err
		}

		switch v := chat.(type) {
		case xmpp.Presence:
			if v.From == fmt.Sprintf("%s@%s/%s", config.User, config.serverName(), config.Resource) {
				return nil
			}
		}
	}
}

func New(config *Config) (*NWWS, error) {
	var err error

	xmpp.DefaultConfig = &tls.Config{
		ServerName:         config.serverName(),
		InsecureSkipVerify: false,
	}

	options := xmpp.Options{
		Host:        config.Server,
		User:        config.User + "@" + config.serverName(),
		Password:    config.Pass,
		Resource:    config.Resource,
		NoTLS:       true,
		StartTLS:    true,
		Debug:       false, // Set to true if you want to see debug information
		Session:     true,
		DialTimeout: 60 * time.Second,
	}

	client, err := options.NewClient()
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("Connected to %s", config.serverName()))

	client.SendOrg(fmt.Sprintf(`<presence xml:lang='en' from='%s@%s' to='%s@%s/%s'><x></x></presence>`, config.User, config.Server, config.Resource, config.Room, config.User))

	err = confirmConnection(client, config)
	if err != nil {
		client.Close()
		return nil, err
	}

	return &NWWS{
		Config: config,
		client: client,
	}, nil
}

func (nwws *NWWS) Start(feed chan *Message) {
	for {
		chat, err := nwws.client.Recv()
		if err != nil {
			log.Println(err)
			err = nwws.client.Close()
			if err != nil {
				log.Println(err)
			}
			break
		}

		switch v := chat.(type) {
		case xmpp.Chat:
			// fmt.Printf("%s\n%s\n", v.Subject, v.Text)
			for _, elem := range v.OtherElem {
				if elem.XMLName.Local == "x" {
					text := strings.ReplaceAll(elem.String(), "\n\n", "\n")
					feed <- &Message{
						Text:       text,
						ReceivedAt: time.Now(),
					}
					nwws.last = time.Now().UTC()
				}
			}
		}
	}
}

func (nwws *NWWS) LastReceived() time.Time {
	return nwws.last
}
