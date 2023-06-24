package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type NukeAccount struct {
	Session *discordgo.Session
	Config  NukeConfig
}

var logger = logrus.New()

func main() {
	nukeAccount := NukeAccount{}
	contents, err := os.ReadFile("nukeconf.toml")
	if err != nil {
		log.Fatalln(err)
	}
	toml.Unmarshal(contents, &nukeAccount.Config)
	logger.Infoln("initializing...")
	session, _ := discordgo.New(nukeAccount.Config.Token)
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentGuildMembers
	nukeAccount.Session = session
	if !nukeAccount.Config.FeatureConfig.AutoNukeConfig.Enabled {
		fmt.Println("Initialized with auto nuke off. Proceed with nuke (press enter)?")
		fmt.Scanln()
	}
	err = nukeAccount.BeginNuke()
	if err != nil {
		logger.Infoln(err)
	}
}
