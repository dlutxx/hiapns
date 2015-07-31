package hiapns

import (
	"encoding/json"
	"fmt"
	"github.com/dlutxx/apns"
	"io/ioutil"
	"log"
	"strings"
)

type Hub struct {
	clients     map[string]*apns.Client
	cnt         Counter
	notifIdBase uint32
}

func initClientAndFeedback(name string, cfg map[string]string) (*apns.Client, error) {
	gw := apns.ProductionGateway
	fbgw := apns.ProductionFeedbackGateway
	if cfg["env"] == "sandbox" {
		gw = apns.SandboxGateway
		fbgw = apns.SandboxFeedbackGateway
	}
	cert, ok := cfg["cert"]
	if !ok {
		return nil, ErrCertMissing
	}
	key, ok := cfg["key"]
	if !ok {
		return nil, ErrKeyMissing
	}

	client, err := apns.NewClient(gw, cert, key, name)
	if err != nil {
		return nil, err
	}

	fb, err := apns.NewFeedback(fbgw, cert, key, name)
	if err != nil {
		return nil, err
	}
	go fb.LogFeedbacks()

	return &client, nil
}

/*
{
	"YCIS-DEV": {
		"cert": "./a.crt",
		"key": "./a.key",
		"env": "sandbox"
	},
	"YCIS-DIST": {
		"cert": "./d.crt",
		"key": "./d.key",
		"env": "production"
	}
}
*/

func NewHub(cfg map[string]map[string]string) *Hub {
	h := Hub{
		clients: map[string]*apns.Client{},
		cnt:     NewCounter(),
	}
	for name, conf := range cfg {
		client, err := initClientAndFeedback(name, conf)
		if err != nil {
			log.Fatalln(err)
		}
		h.clients[name] = client
	}
	return &h
}

func (hb Hub) Send(name string, n *apns.Notification) error {
	client, ok := hb.clients[name]
	if !ok {
		return ErrNoClient
	}
	return client.Send(*n)
}

func parseConfigFile(filepath string) (map[string]map[string]string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	cfg := map[string]map[string]string{}

	if err = json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("Invalid Config: %v", err)
	}
	for name, c := range cfg {
		cert := c["cert"]
		if strings.HasPrefix(cert, "/") || strings.HasPrefix(cert, "./") {
			bs, err := ioutil.ReadFile(cert)
			if err != nil {
				return nil, fmt.Errorf("%v cert %v", name, err)
			}
			cfg[name]["cert"] = string(bs)
		}
		key := c["key"]
		if strings.HasPrefix(key, "/") || strings.HasPrefix(key, "./") {
			bs, err := ioutil.ReadFile(key)
			if err != nil {
				return nil, fmt.Errorf("%v cert %v", name, err)
			}
			cfg[name]["key"] = string(bs)
		}
	}
	return cfg, nil
}

func NewHubFromCfgFile(filepath string) (*Hub, error) {
	cfg, err := parseConfigFile(filepath)
	if err != nil {
		return nil, err
	}
	return NewHub(cfg), nil
}
