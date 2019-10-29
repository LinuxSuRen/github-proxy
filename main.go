package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var token string
var port int

func init() {
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "token for target repo")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for http server")
	if err := rootCmd.MarkFlagRequired("token"); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "github",
	Short: "github proxy",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/api/webhook", webhookHandler)

		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal(err)
		}
	},
}

func main()  {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// curl -X POST -H "Accept: application/vnd.github.everest-preview+json" -H "Authorization: token ${TOKEN}" -i "https://api.github.com/repos/jenkins-zh/jenkins-zh/dispatches" -d '{"event_type":"repository_dispatch"}'
func webhookHandler(writer http.ResponseWriter, request *http.Request)  {
	payload, err := json.Marshal(map[string]string{
		"event_type": "repository_dispatch",
	})
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST",
		"https://api.github.com/repos/jenkins-zh/jenkins-zh/dispatches",
		bytes.NewBuffer(payload))
	if err == nil {
		req.Header.Add("Accept", "application/vnd.github.everest-preview+json")
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
		client := http.Client{}

		var response *http.Response
		response, err = client.Do(req)
		if err != nil {
			log.Print(err)
		} else if response.StatusCode != 204 {
			fmt.Println(response)
		}
	}
}
