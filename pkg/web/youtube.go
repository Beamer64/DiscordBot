package web

import (
	"context"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var StopPlaying chan bool
var IsPlaying bool
var MpFileQueue []string

func GetYtAudioLink(s *discordgo.Session, m *discordgo.Message, link string) (mpFileLink string, fileName string, err error) {
	replacer := strings.NewReplacer("m.", "", "youtube", "youtubex2")
	url := replacer.Replace(link)

	ctx, cancel := chromedp.NewContext(context.Background()) // options: chromedp.WithDebugf(log.Printf)
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	var res string
	var ok *bool

	msg, err := s.ChannelMessageEdit(m.ChannelID, m.ID, "Prepping vidya...20% [##        ]")
	if err != nil {
		return "", "", err
	}

	// navigate to url and get redirect url
	NavTasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Location(&res),
	}
	// run navigate task list
	err = chromedp.Run(ctx, NavTasks)
	if err != nil {
		return "", "", err
	}

	// navigate to redirect and click button
	button := "/html/body/div[1]/main/section[2]/div[2]/div/div[2]/div/div[2]/div/a"
	clickTasks := chromedp.Tasks{
		chromedp.Navigate(res),
		chromedp.Click(button),
	}
	// run clickTask list
	err = chromedp.Run(ctx, clickTasks)
	if err != nil {
		return "", "", err
	}

	msg, err = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Prepping vidya...40% [####      ]")
	if err != nil {
		return "", "", err
	}

	// wait for page to load and get button redirect url
	searchElem := "/html/body/div/main/section[1]/div/div/div[5]/div/div[1]/div"
	waitTasks := chromedp.Tasks{
		chromedp.WaitNotPresent(searchElem),
		chromedp.Location(&res),
	}

	// run waitTasks list
	err = chromedp.Run(ctx, waitTasks)
	if err != nil {
		return "", "", err
	}

	msg, err = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Prepping vidya...50% [#####     ]")
	if err != nil {
		return "", "", err
	}

	// navigate to button redirect and get download link
	button = "/html/body/div[1]/main/section/div/div[2]/div/div[2]/div[1]/div[3]/a[1]"
	resURL := res
	navTasks := chromedp.Tasks{
		chromedp.Navigate(resURL),
		chromedp.AttributeValue(button, "href", &res, ok),
	}

	// run navTasks list
	err = chromedp.Run(ctx, navTasks)
	if err != nil {
		return "", "", err
	}

	msg, err = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Prepping vidya...70% [#######   ]")
	if err != nil {
		return "", "", err
	}

	// navigate to download link to parse network response
	getLinkTasks := chromedp.Tasks{
		chromedp.Navigate(res),
	}

	// listen for response containing mp3 link
	mpLink := ""
	chromedp.ListenTarget(
		ctx, func(ev interface{}) {
			if ev, ok := ev.(*network.EventResponseReceived); ok {
				if strings.Contains(ev.Response.URL, ".mp3") {
					mpLink = ev.Response.URL
					//fmt.Println("closing alert:", ev.Response)
				}
			}
		},
	)

	// run getLinkTasks list
	err = chromedp.Run(ctx, getLinkTasks)
	if err != nil {
		if !strings.Contains(err.Error(), "net::ERR_ABORTED") {
			return "", "", err
		}
	}

	msg, err = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Prepping vidya...90% [######### ]")
	if err != nil {
		return "", "", err
	}

	time.AfterFunc(
		2*time.Second, func() {
			_ = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		},
	)

	fileName = strings.SplitAfterN(mpLink, "/", 12)[7]

	return mpLink, fileName, nil
}

func DownloadMpFile(m *discordgo.MessageCreate, link string, fileName string) error {
	// Get the data
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)
	if err != nil {
		return err
	}

	// Create the dir
	dir := fmt.Sprintf("%s/Audio", m.GuildID)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		// does not exist
		err = os.MkdirAll(dir, 0777)
		fmt.Println(fmt.Sprintf("Dir created: %s", dir))
	}
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filepath.Join(dir, filepath.Base(fileName)))
	if err != nil {
		return err
	}
	fmt.Println("Created File")

	defer func(out *os.File) {
		err = out.Close()
	}(out)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func PlayAudioFile(dgv *discordgo.VoiceConnection, fileName string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	dir := fmt.Sprintf("%s/Audio", m.GuildID)

	cleanFileName := ""
	var err error
	if fileName != "" {
		cleanFileName, err = FormatAudioFileName(fileName)
		if err != nil {
			return err
		}

		if !IsPlaying {
			if fileName != "" {
				MpFileQueue = append(MpFileQueue, filepath.Join(dir, filepath.Base(fileName)))
			}

			IsPlaying = true
			for _, v := range MpFileQueue {
				fmt.Println("PlayAudioFile: ", v)

				_, err = s.ChannelMessageSend(m.ChannelID, "Now playing: "+cleanFileName)
				if err != nil {
					return err
				}

				dgvoice.PlayAudioFile(dgv, v, StopPlaying)
			}
			//remove file from queue
			MpFileQueue = nil
			//MpFileQueue = append(MpFileQueue[:i], MpFileQueue[i+1:]...)

			if dgv != nil {
				err = dgv.Disconnect()
				if err != nil {
					return err
				}
			}

			err = MpFileCleanUp(dir)
			if err != nil {
				return err
			}

			/*IsPlaying = false
			if len(MpFileQueue) > 0 {
				err := PlayAudioFile(dgv, "", m, s)
				if err != nil {
					return err
				}

			} else {
				if dgv != nil {
					err := dgv.Disconnect()
					if err != nil {
						return err
					}
				}*/

			/*err := MpFileCleanUp(dir)
				if err != nil {
					return err
				}
			}*/

		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added to queue: %s", cleanFileName))
			if err != nil {
				return err
			}

			MpFileQueue = append(MpFileQueue, filepath.Join(dir, filepath.Base(fileName)))
		}
	}

	return nil
}

// FormatAudioFileName formats audio file name to look better
func FormatAudioFileName(fileName string) (string, error) {
	//split at "/"
	splitName := strings.SplitAfterN(fileName, "\\", 3)
	fileName = splitName[2]

	//replace characters
	replacer := strings.NewReplacer("/", "", "_", " ", "-", "", ".mp3", "")
	fileName = replacer.Replace(fileName)

	//remove numbers
	numRegex, err := regexp.Compile("[0-9]")
	fileName = numRegex.ReplaceAllString(fileName, "")
	if err != nil {
		return "", err
	}

	//capitalize first letters
	caser := cases.Title(language.AmericanEnglish)
	fileName = caser.String(fileName)

	return fileName, nil
}

// MpFileCleanUp clear out Audio directory
func MpFileCleanUp(dir string) error {
	MpFileQueue = nil

	fmt.Println("\nRunning Cleanup")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if strings.Contains(filepath.Join(dir, filepath.Base(f.Name())), ".mp3") {
			err = os.Remove(filepath.Join(dir, filepath.Base(f.Name())))
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Cleanup Finished")
	return nil
}
