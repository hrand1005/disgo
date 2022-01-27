package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var LuckyWheelCommand = &discordgo.ApplicationCommand{
	Name:        "lucky-wheel",
	Description: "You know what this does.",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "items",
			Description:  "Lucky wheel participants to be chosen at random",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
	},
}

func LuckyWheelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// will have to check Options[0] == choice
				Content: fmt.Sprintf("SPIN DA WHEEL!\n%q", data.Options[0].StringValue()),
				//Content: data.Options[0].StringValue(),
			},
		})
		onlineMembers, err := onlineMembersForGuild(s, i.Interaction.GuildID)
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: onlineMembers[0].Mention(),
		})
		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Encountered error :(",
			})
			log.Println("Encountered error spinning lucky-wheel: ", err)
			return
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		onlineMembers, err := onlineMembersForGuild(s, i.Interaction.GuildID)
		log.Printf("Online users: %v\n", onlineMembers)
		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Encountered error :(",
			})
			log.Println("Encountered error retrieving online users: ", err)
			return
		}
		choices := []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "all users",
				Value: "adds everyone online to the lucky wheel",
				//Value: "all users",
			},
			{
				Name:  "channel users",
				Value: "adds everyone active in the specified voice channel to the lucky wheel",
				//Value: channelUsers,
			},
		}

		if data.Options[0].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[0].StringValue(),
				Value: data.Options[0].StringValue(),
			})
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})

		if err != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Encountered error :(",
			})
			log.Println("Encountered error with lucky-wheel autocomplete: ", err)
			return
		}
	}
}

func onlineMembersForGuild(s *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	g, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	log.Printf("GuildID: %v\n", guildID)
	log.Printf("Guild Info: %v\n", g)
	log.Printf("Guild Members: %v\n", g.Members)
	log.Printf("A member mentioned: %s\n", g.Members[0].Mention())
	log.Printf("Guild Presences: %v\n", g.Presences)

	return g.Members, nil
}
