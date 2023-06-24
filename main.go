package main

import (
	"fmt"
	"log"
	"math/rand"
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

func (na NukeAccount) getGuildMemberIDs() []string {
	ids := []string{}
	collected := make(chan bool)
	na.Session.AddHandlerOnce(func(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
		for _, member := range chunk.Members {
			ids = append(ids, member.User.ID)
		}
		collected <- true
	})
	na.Session.RequestGuildMembers(na.Config.GuildID, "", 0, fmt.Sprint(rand.Intn(10000000)), false)
	<-collected
	return ids
}
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
	if !nukeAccount.Config.FeatureConfig.AutoNuke {
		fmt.Println("Initialized with auto nuke off. Proceed with nuke (press enter)?")
		fmt.Scanln()
	}
	err = nukeAccount.BeginNuke()
	if err != nil {
		logger.Infoln(err)
	}
}
