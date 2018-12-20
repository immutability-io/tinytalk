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
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
)

func isKeybaseMessage(subject string) bool {
	return strings.Contains(strings.ToLower(subject), "keybase:")
}

func printMsg(m *nats.Msg, i int) {
	if isKeybaseMessage(m.Subject) {
		filename := fmt.Sprintf("%s.%d.enc", m.Subject, i)
		ioutil.WriteFile(filename, m.Data, 0777)
		log.Printf("[#%d] Received (ENCRYPTED) on [%s]: 'keybase pgp decrypt -i %s'", i, m.Subject, filename)
	} else {
		log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
	}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatal("Exiting, no servers available")
	}))
	return opts
}

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen will listen on the immutability channel",
	Long:  `Keybase encrypted messages are written to the file system.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Connect Options.
		opts := []nats.Option{nats.Name("Immutability Subscriber")}
		opts = setupConnOptions(opts)
		opts = append(opts, nats.UserCredentials("tinytalk.creds"))

		// Connect to NATS
		nc, err := nats.Connect("connect.ngs.global", opts...)
		if err != nil {
			log.Fatal(err)
		}

		i := 0

		nc.Subscribe(SubscribeChannel, func(msg *nats.Msg) {
			i++
			printMsg(msg, i)
		})
		nc.Flush()

		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on [%s]", SubscribeChannel)

		runtime.Goexit()
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}
