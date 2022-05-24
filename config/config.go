package config

import (
	"context"
	"io/ioutil"

	"github.com/fox-one/mixin-sdk-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseUrl string `yaml:"base_url"`
	User    User   `yaml:"user"`
	Mixin   Mixin  `yaml:"mixin"`
}

type User struct {
	UserID string `yaml:"user_id"`
	Token  string `yaml:"token"`
}

type Mixin struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	SessionID    string `yaml:"session_id"`
	PinCode      string `yaml:"pin_code"`
	PinToken     string `yaml:"pin_token"`
	PrivateKey   string `yaml:"private_key"`
}

var client *mixin.Client

func (m *Mixin) Client() *mixin.Client {
	if client != nil {
		return client
	}

	ks := mixin.Keystore{
		ClientID:   m.ClientID,
		SessionID:  m.SessionID,
		PinToken:   m.PinToken,
		PrivateKey: m.PrivateKey,
	}

	client, err := mixin.NewFromKeystore(&ks)
	if err != nil {
		panic(err)
	}

	if err := client.VerifyPin(context.Background(), m.PinCode); err != nil {
		panic(err)
	}

	return client
}

func Load(path string, c *Config) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		panic(err)
	}
}
