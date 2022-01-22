package events

import (
	"github.com/beamer64/discordBot/pkg/games"
	"github.com/beamer64/discordBot/pkg/gcp"
	"github.com/beamer64/discordBot/pkg/ssh"
	"github.com/beamer64/discordBot/pkg/web_scrape"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func (d *MessageCreateHandler) testMethod(s *discordgo.Session, m *discordgo.MessageCreate, param string) error {
	/*err := d.playYoutubeLink(s, m, param)
	if err != nil {
		return err
	}*/

	return nil
}

func (d *MessageCreateHandler) sendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var cmds []string
	if d.memberHasRole(s, m, d.cfg.Configs.Settings.BotAdminRole) { // bot mod

		//the weird order is, so it looks better until all are converted (listed $ then /)
		cmds = append(cmds, d.cfg.Cmd.Desc.Tuuck, d.cfg.Cmd.Desc.LMGTFY)

		if d.cfg.Configs.Server.MachineIP != "" {
			cmds = append(cmds, d.cfg.Cmd.Desc.ServerStatus, d.cfg.Cmd.Desc.StartServer, d.cfg.Cmd.Desc.StopServer)
		}

		cmds = append(
			cmds, d.cfg.Cmd.Desc.Horoscope, d.cfg.Cmd.Desc.Version, d.cfg.Cmd.Desc.CoinFlip, d.cfg.Cmd.Desc.Play,
			d.cfg.Cmd.Desc.Stop, d.cfg.Cmd.Desc.Queue,
		)

		if d.cfg.Configs.Keys.InsultAPI != "" {
			cmds = append(cmds, d.cfg.Cmd.Desc.Insult)
		}

	} else {
		//the weird order is, so it looks better until all are converted (listed $ then /)
		cmds = append(cmds, d.cfg.Cmd.Desc.Tuuck, d.cfg.Cmd.Desc.LMGTFY)

		if d.cfg.Configs.Server.MachineIP != "" {
			cmds = append(cmds, d.cfg.Cmd.Desc.ServerStatus)
		}

		cmds = append(
			cmds, d.cfg.Cmd.Desc.Horoscope, d.cfg.Cmd.Desc.CoinFlip, d.cfg.Cmd.Desc.Play, d.cfg.Cmd.Desc.Stop,
			d.cfg.Cmd.Desc.Queue,
		)

		if d.cfg.Configs.Keys.InsultAPI != "" {
			cmds = append(cmds, d.cfg.Cmd.Desc.Insult)
		}
	}

	cmdDesc := ""
	for _, command := range cmds {
		cmdDesc = cmdDesc + "\n" + command
	}

	_, err := s.ChannelMessageSend(m.ChannelID, cmdDesc)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReactionHandler) sendLmgtfy(s *discordgo.Session, m *discordgo.Message) error {
	lmgtfyURL := CreateLmgtfyURL(m.Content)

	lmgtfyShortURL, err := ShortenURL(lmgtfyURL)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "\""+m.Content+"\""+"\n"+lmgtfyShortURL)
	if err != nil {
		return err
	}

	return nil
}

func (d *MessageCreateHandler) displayHoroscope(s *discordgo.Session, m *discordgo.MessageCreate, param string) error {
	horoscope, err := web_scrape.ScrapeSign(param)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, horoscope)
	if err != nil {
		return err
	}

	return nil
}

func (d *MessageCreateHandler) playNIM(s *discordgo.Session, m *discordgo.MessageCreate, param string) error {
	if strings.HasPrefix(param, "<@") {
		err := games.StartNim(s, m, param, true)
		if err != nil {
			return err
		}

	} else {
		if param == "" {
			err := games.StartNim(s, m, param, false)
			if err != nil {
				return err
			}

		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.Invalid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *MessageCreateHandler) sendStartUpMessages(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// sleep for 1 minute while saying funny things and to wait for instance to start up
	sm := 0
	for i := 1; i < 5; i++ {
		loadingMessage := getRandomLoadingMessage(d.cfg.LoadingMessages)
		time.Sleep(3 * time.Second)

		_, err := s.ChannelMessageSend(m.ChannelID, loadingMessage)
		if err != nil {
			return err
		}

		sm += i
	}
	time.Sleep(3 * time.Second)
	return nil
}

func (d *MessageCreateHandler) startServer(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if d.cfg.Configs.Server.MachineIP != "" { // check if Minecraft server is set up
		client, err := gcp.NewGCPClient("config/auth.json", d.cfg.Configs.Server.Project_ID, d.cfg.Configs.Server.Zone)
		if err != nil {
			return err
		}

		err = client.StartMachine("instance-2-minecraft")
		if err != nil {
			return err
		}

		sshClient, err := ssh.NewSSHClient(d.cfg.Configs.Server.SSHKeyBody, d.cfg.Configs.Server.MachineIP)
		if err != nil {
			return err
		}

		status, serverUp := sshClient.CheckServerStatus(sshClient)
		if serverUp {
			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.ServerUP+status)
			if err != nil {
				return err
			}

		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.WindUp)
			if err != nil {
				return err
			}

			_, err = sshClient.RunCommand("docker container start 06ae729f5c2b")
			if err != nil {
				return err
			}

			err = d.sendStartUpMessages(s, m)
			if err != nil {
				return err
			}

			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.FinishOpperation)
			if err != nil {
				return err
			}
		}

	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.MCServerError)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *MessageCreateHandler) stopServer(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if d.cfg.Configs.Server.MachineIP != "" { // check if Minecraft server is set up
		sshClient, err := ssh.NewSSHClient(d.cfg.Configs.Server.SSHKeyBody, d.cfg.Configs.Server.MachineIP)
		if err != nil {
			return err
		}

		status, serverUp := sshClient.CheckServerStatus(sshClient)
		if serverUp {
			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.WindDown)

			_, err = sshClient.RunCommand("docker container stop 06ae729f5c2b")
			if err != nil {
				return err
			}

			client, errr := gcp.NewGCPClient("config/auth.json", d.cfg.Configs.Server.Project_ID, d.cfg.Configs.Server.Zone)
			if errr != nil {
				return err
			}

			err = client.StopMachine("instance-2-minecraft")
			if err != nil {
				return err
			}

			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.FinishOpperation)
			if err != nil {
				return err
			}

		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.ServerDOWN+status)
			if err != nil {
				return err
			}
		}

	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.MCServerError)
		if err != nil {
			return err
		}
	}
	return nil
}

// d.sendServerStatusAsMessage Sends the current server status as a message in discord
func (d *MessageCreateHandler) sendServerStatusAsMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	client, err := gcp.NewGCPClient("config/auth.json", d.cfg.Configs.Server.Project_ID, d.cfg.Configs.Server.Zone)
	if err != nil {
		return err
	}

	err = client.StartMachine("instance-2-minecraft")
	if err != nil {
		return err
	}

	sshClient, err := ssh.NewSSHClient(d.cfg.Configs.Server.SSHKeyBody, d.cfg.Configs.Server.MachineIP)
	if err != nil {
		return err
	}

	status, serverUp := sshClient.CheckServerStatus(sshClient)
	if serverUp {
		_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.CheckStatusUp+status)
		if err != nil {
			return err
		}

	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, d.cfg.Cmd.Msg.CheckStatusDown+status)
		if err != nil {
			return err
		}
	}
	return nil
}
