// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	nats "github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
)

// Message is the subject to sub/pub to
var Message string

// Subject is the subject to sub/pub to
var Subject string

// sayCmd represents the say command
var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "say will send a message to the immutability channel",
	Long: `say will send a short message to the immutability channel. Use the
--keybase option to encrypt the message.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		opts := []nats.Option{nats.Name("Immutability Publisher")}
		opts = append(opts, nats.UserCredentials("tinytalk.creds"))
		nc, err := nats.Connect("connect.ngs.global", opts...)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()

		message := []byte(Message)
		subject := DefaultSubject
		if Keybase != "" {
			_, message, err = keybaseEncrypt(Keybase, message)
			if err != nil {
				log.Fatal(err)
			}
			if Subject != "" {
				subject = PublishChannel + Subject + "." + Keybase
			} else {
				subject = PublishChannel + Keybase
			}
		} else if Subject != "" {
			subject = PublishChannel + Subject
		}
		nc.Publish(subject, message)
		nc.Flush()

		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Sent [%s] : '%s'\n", subject, Message)
		}
	},
}

func init() {
	sayCmd.PersistentFlags().StringVarP(&Keybase, "keybase", "k", "", "Keybase identity to encrypt the message with")
	sayCmd.PersistentFlags().StringVarP(&Message, "message", "m", "", "Message to publish")
	sayCmd.PersistentFlags().StringVarP(&Subject, "subject", "s", "", "Subject of message")
	sayCmd.MarkPersistentFlagRequired("message")
	rootCmd.AddCommand(sayCmd)

}
