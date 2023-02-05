package main

import (
	"context"
	"encoding/json"
	"fmt"
	gpt "github.com/PullRequestInc/go-gpt3"
	"github.com/enescakir/emoji"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	ApiKey string `json:"api_key"`
}

func main() {
	configFile, err := os.Open("config.json")
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic(err)
		}
	}(configFile)

	var config Config

	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &config)

	if config.ApiKey == "" || config.ApiKey == "your-api-key" {
		panic("API_KEY is required.")
	}

	fmt.Print("  ___ _  _   _ _____ ___ ___ _____ \n / __| || | /_\\_   _/ __| _ \\_   _|\n| (__| __ |/ _ \\| || (_ |  _/ | |  \n \\___|_||_/_/ \\_\\_| \\___|_|   |_|  \n                                   \n")
	fmt.Print("A simple program for communicating with GPT through the terminal.\n")
	fmt.Print("Github Repository: https://github.com/j3yzz/chatgpt-go \n\n")

	ctx := context.Background()
	gptClient := gpt.NewClient(config.ApiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with OpenAI ChatGPT in terminal.",
		Run: func(cmd *cobra.Command, args []string) {
			quit := false
			for !quit {
				prompt := promptui.Prompt{
					Label: emoji.Man,
				}

				result, _ := prompt.Run()

				question := result
				questionParam, exit := validateQuestion(question)
				if exit {
					fmt.Println("Goodbye, See you soon!", emoji.WavingHand)
					break
				}
				switch questionParam {
				case "quit":
					quit = true
				case "":
					continue
				default:
					GetResponse(gptClient, ctx, question)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

func GetResponse(client gpt.Client, ctx context.Context, q string) {
	fmt.Printf("\n%v:\n", emoji.Robot)

	err := client.CompletionStreamWithEngine(ctx, gpt.TextDavinci003Engine, gpt.CompletionRequest{
		Prompt: []string{
			q,
		},
		MaxTokens:   gpt.IntPtr(3000),
		Temperature: gpt.Float32Ptr(0),
	}, func(response *gpt.CompletionResponse) {
		fmt.Print(response.Choices[0].Text)
	})

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("\n--------------------------------------------------\n")
}

func validateQuestion(question string) (string, bool) {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "cls", "exit", "goodbye"}
	exit := false
	for _, x := range keywords {
		if quest == x {
			exit = true
			return "", exit
		}
	}
	return quest, exit
}
