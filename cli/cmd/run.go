/*
Copyright © 2021 ks6088ts <ks6088ts@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/spf13/cobra"
)

type runOptions struct {
	slackSigningSecret string
	slackBotToken      string
	port               int
}

var (
	o = &runOptions{}
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a slack event listener",
	Long:  `Run a slack event listener`,
	Run: func(cmd *cobra.Command, args []string) {
		api := slack.New(o.slackBotToken)

		http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
			verifier, err := slack.NewSecretsVerifier(r.Header, o.slackSigningSecret)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bodyReader := io.TeeReader(r.Body, &verifier)
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := verifier.Ensure(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			switch eventsAPIEvent.Type {
			case slackevents.URLVerification:
				var res *slackevents.ChallengeResponse
				if err := json.Unmarshal(body, &res); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/plain")
				if _, err := w.Write([]byte(res.Challenge)); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch event := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					message := strings.Split(event.Text, " ")
					if len(message) < 2 {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					command := message[1]
					switch command {
					case "ping":
						if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("pong", false)); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
				}
			}
		})

		log.Println("[INFO] Server listening")
		if err := http.ListenAndServe(":"+strconv.Itoa(o.port), nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&o.port, "port", "p", 8080, "port number")
	runCmd.Flags().StringVarP(&o.slackSigningSecret, "secret", "s", "default", "slack signing secret")
	runCmd.Flags().StringVarP(&o.slackBotToken, "token", "t", "default", "slack bot token")

	err := runCmd.MarkFlagRequired("secret")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = runCmd.MarkFlagRequired("token")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
