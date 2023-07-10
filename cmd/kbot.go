package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

var (
	// Map to store links, their last known response codes, and content
	linkResponses = make(map[string]responseInfo)

	// Map to store user IDs and their requested links
	userLinks = make(map[int64][]string)
)

type responseInfo struct {
	StatusCode int
	Content    string
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			// Add this log statement to check the received payload
			log.Printf("Received payload: %s", m.Message().Payload)

			payload := m.Message().Payload

			if payload == "" {
				// Send the usage guide to the user
				err := m.Send("Welcome to kbot!\n\nIn order to start monitoring a link, you can send the following message:\n/s <LINK URL>\n\nTo stop receiving alerts, you can use the command:\n/s stop")
				if err != nil {
					log.Printf("Failed to send response: %s", err)
					return err
				}
				return nil
			}

			userID := m.Sender().ID

			if payload == "stop" {
				// Remove all links associated with the user ID
				delete(userLinks, userID)

				// Send a confirmation message to the user
				err := m.Send("You have stopped receiving alerts.")
				if err != nil {
					log.Printf("Failed to send response: %s", err)
					return err
				}
			} else {
				if isURL(payload) {
					// Store the link for the user
					userLinks[userID] = []string{payload}

					if lastRespInfo, ok := linkResponses[payload]; ok {
						// Check the link's current response info
						respInfo, err := checkResponseInfo(payload)
						if err != nil {
							log.Printf("Failed to check response info for %s: %s", payload, err)
							return err
						}

						if respInfo.StatusCode != lastRespInfo.StatusCode || respInfo.Content != lastRespInfo.Content {
							// Notify the user about the change
							notifyUser(userID, payload, respInfo.StatusCode)

							// Update the last known response info for the link
							linkResponses[payload] = respInfo
						}
					} else {
						// Store the link's initial response info
						respInfo, err := checkResponseInfo(payload)
						if err != nil {
							log.Printf("Failed to check response info for %s: %s", payload, err)
							return err
						}

						linkResponses[payload] = respInfo

						// Notify the user that the link is being monitored
						notifyUser(userID, payload, respInfo.StatusCode)
					}
				} else {
					// Handle other text messages (e.g., "hello" message)
					switch payload {
					case "hello":
						err := m.Send(fmt.Sprintf("Hello, I'm Kbot %s!", appVersion))
						if err != nil {
							log.Printf("Failed to send response: %s", err)
							return err
						}
					}
				}
			}

			return nil
		})

		go func() {
			for {
				// Sleep for 5 seconds
				time.Sleep(5 * time.Second)

				// Check all links for changes
				for link := range linkResponses {
					respInfo, err := checkResponseInfo(link)
					if err != nil {
						log.Printf("Failed to check response info for %s: %s", link, err)
						continue
					}

					lastRespInfo := linkResponses[link]
					if respInfo.StatusCode != lastRespInfo.StatusCode || respInfo.Content != lastRespInfo.Content {
						// Notify the user about the change
						for _, userID := range getUsersForLink(link) {
							notifyUser(userID, link, respInfo.StatusCode)
						}

						// Update the last known response info for the link
						linkResponses[link] = respInfo
					}
				}
			}
		}()

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}

// Function to check the HTTP response info (status code and content) of a URL
func checkResponseInfo(url string) (responseInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return responseInfo{}, err
	}
	defer resp.Body.Close()

	// Read the content of the response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseInfo{}, err
	}

	return responseInfo{
		StatusCode: resp.StatusCode,
		Content:    string(content),
	}, nil
}

// Function to check if a string is a valid URL
func isURL(str string) bool {
	_, err := http.Get(str)
	return err == nil
}

// Function to get the user IDs associated with a link
func getUsersForLink(link string) []int64 {
	var users []int64
	for userID, links := range userLinks {
		for _, l := range links {
			if l == link {
				users = append(users, userID)
			}
		}
	}
	return users
}

// Function to notify a user about the link content change
func notifyUser(userID int64, link string, statusCode int) {
	// Retrieve the bot instance
	kbot, err := telebot.NewBot(telebot.Settings{
		URL:    "",
		Token:  TeleToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Printf("Failed to create bot instance: %s", err)
		return
	}

	// Send notification to the user
	_, err = kbot.Send(&telebot.User{ID: userID}, fmt.Sprintf("The content of the link %s has changed! New status code: %d", link, statusCode))
	if err != nil {
		log.Printf("Failed to send notification for link %s to user %d: %s", link, userID, err)
	}
}
