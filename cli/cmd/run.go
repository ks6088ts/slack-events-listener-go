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

// Package cmd ...
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ks6088ts/slack-events-listener-go/internal"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/spf13/cobra"
)

type runOptions struct {
	slackSigningSecret string
	slackBotToken      string
	port               int
	verbose            bool
	credentialsPath    string
	sheetId            string
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
		signingSecret := o.slackSigningSecret
		sheetClient, err := internal.NewSheetClient(o.credentialsPath, o.sheetId)
		if err != nil {
			log.Println(err)
		}

		http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if _, err := sv.Write(body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := sv.Ensure(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if o.verbose {
				log.Printf("%v\n", eventsAPIEvent)
			}

			switch eventsAPIEvent.Type {
			case slackevents.URLVerification:
				var r *slackevents.ChallengeResponse
				err := json.Unmarshal([]byte(body), &r)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text")
				if _, err := w.Write([]byte(r.Challenge)); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				if o.verbose {
					log.Printf("%v\n", innerEvent)
				}
				switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					if o.verbose {
						log.Println(ev)
					}
					if _, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false)); err != nil {
						log.Println(err)
					}
				case *slackevents.ReactionAddedEvent:
					if o.verbose {
						log.Println(ev)
					}
					log.Printf("do something with channel=%v, reaction=%v\n", ev.Item.Channel, ev.Reaction)
				case *slackevents.MessageEvent:
					if o.verbose {
						log.Println(ev)
					}
					// filter other subtypes: https://api.slack.com/events/message#subtypes
					if ev.SubType != "" {
						return
					}

					// convert timestamp
					ts, err := internal.GetTimeFromSlackTimeStamp(ev.TimeStamp)
					if err != nil {
						log.Println(err)
						return
					}

					// update sheet
					err = sheetClient.AppendValues(ev.Text, ts.Format("2006年01月02日 15時04分05秒"))
					if err != nil {
						log.Println(err)
						return
					}
					if o.verbose {
						log.Printf("successfully appended %v, %v\n", ev.Text, ts.Format("2006年01月02日 15時04分05秒"))
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
	runCmd.Flags().BoolVarP(&o.verbose, "verbose", "v", false, "verbosity")
	runCmd.Flags().StringVarP(&o.credentialsPath, "credentials", "c", "default", "path to credentials file")
	runCmd.Flags().StringVarP(&o.sheetId, "sheetId", "i", "default", "sheet ID")

	for _, requiredCmd := range []string{
		"secret",
		"token",
	} {
		err := runCmd.MarkFlagRequired(requiredCmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}
