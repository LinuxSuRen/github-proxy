package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/linuxsuren/github-proxy/pkg/module"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"github.com/linuxsuren/github-proxy/pkg"
	"os/exec"
	"strings"
)

var token string
var port int
var task string

func init() {
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "token for target repo")
	rootCmd.Flags().StringVarP(&task, "task", "", "", "task path")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for http server")
	if err := rootCmd.MarkFlagRequired("token"); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "github",
	Short: "github proxy",
	Run: func(cmd *cobra.Command, args []string) {
		task, err := pkg.ParsePipelineTask(task)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(task)

		http.HandleFunc("/api/webhook", func (writer http.ResponseWriter, request *http.Request){
			if task != nil && handleWebHookFromTask(task, request) != nil {
				webhookHandler(writer, request)
			}
		})

		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal(err)
		}
	},
}

func handleWebHookFromTask(task *module.PipelineTask, request *http.Request) error {
	if len(task.Triggers) == 0 {
		return fmt.Errorf("no triggers")
	}

	fmt.Println(request)
	for _, trigger := range task.Triggers {
		if len(trigger.Headers) == 0 {
			continue
		}

		for _, header := range trigger.Headers {
			val := request.Header.Get(header.Name)
			if val == header.Value {
				fmt.Println("found task", task.Name)

				array := strings.Split(task.Script, " ")
				cmd := exec.Command(array[0], array[1:]...)
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}

				return nil
			}
		}
	}

	return fmt.Errorf("no matched triggerd")
}

func main()  {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// curl -X POST -H "Accept: application/vnd.github.everest-preview+json" -H "Authorization: token ${TOKEN}" -i "https://api.github.com/repos/jenkins-zh/jenkins-zh/dispatches" -d '{"event_type":"repository_dispatch"}'
func webhookHandler(writer http.ResponseWriter, request *http.Request) {
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
