package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type GuildJoinLeaveHandler struct{}
type ReactionHandler struct{}
type ReadyHandler struct{}

func NewGuildJoinLeaveHandler() *GuildJoinLeaveHandler {
	return &GuildJoinLeaveHandler{}
}

func NewReactionHandler() *ReactionHandler {
	return &ReactionHandler{}
}

func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{}
}

func (h *ReadyHandler) ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("Bot Session is ready!")
	fmt.Printf("Logged in as %s\n", e.User.String())
}

func (r *ReactionHandler) ReactHandlerAdd(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
	if mr.MessageReaction.Emoji.Name == "lmgtfy" {
		msg, _ := s.ChannelMessage(mr.ChannelID, mr.MessageID)

		err := r.sendLmgtfy(s, msg)
		if err != nil {
			fmt.Printf("%+v", errors.WithStack(err))
		}
	}
}

func (d *GuildJoinLeaveHandler) GuildJoinHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		fmt.Println("Failed getting guild object: ", err)
		fmt.Printf("%+v", errors.WithStack(err))
		return
	}

	fmt.Printf("Hey! Look at this goofy goober! %s joined our %s server!\n", e.Member.User.String(), guild.Name)
}

func (d *GuildJoinLeaveHandler) GuildLeaveHandler(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		fmt.Println("Failed getting guild object: ", err)
		fmt.Printf("%+v", errors.WithStack(err))
		return
	}

	fmt.Printf("%s left the server %s\n Seacrest OUT..", e.Member.User.String(), guild.Name)
}
