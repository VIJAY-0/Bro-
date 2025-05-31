package agent

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

type Agent struct {
	ctx       context.Context
	client    *genai.Client
	history   []*genai.Content
	chat      *genai.Chat
	reader    *bufio.Reader
	memory    Memory
	dbs       DBs
	registers Registers
}

func (ag *Agent) Activate() {

	data, err := os.ReadFile("Prompts/BasePrompt.pmt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	BasePrompt := string(data)
	ag.addHistory([]string{BasePrompt})
	ag.reader = bufio.NewReader(os.Stdin)
	ag.init("AIzaSyCt4xwCe7rolx-7oPh7JOeOHPTT4IquifI")
	ag.createChat()

	err = ag.startTask()
	if err != nil {
		panic(err)
	}

}

func (ag *Agent) InitDBs(DBroot string) {
	ag.dbs = NewDBset(DBroot)
}

func (ag *Agent) addHistory(strs []string) {
	for _, str := range strs {
		ag.history = append(ag.history, genai.NewContentFromText(str, genai.RoleUser))
	}
}

func (ag *Agent) init(APIkeys string) {

	var err error
	ag.ctx = context.Background()
	ag.client, err = genai.NewClient(ag.ctx, &genai.ClientConfig{
		APIKey:  APIkeys,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (ag *Agent) createChat() error {
	var err error
	ag.chat, err = ag.client.Chats.Create(ag.ctx, "gemini-2.0-flash", nil, ag.history)
	if err != nil {
		return err
	}
	return err
}

func (ag *Agent) prompt(prompt string) (string, error) {
	resp, err := ag.chat.SendMessage(ag.ctx, genai.Part{Text: prompt})
	if err != nil {
		return "", err
	}
	response, err := ag.getText(resp)
	return response, err
}

func (ag *Agent) getText(response *genai.GenerateContentResponse) (string, error) {
	if len(response.Candidates) > 0 {
		text := response.Candidates[0].Content.Parts[0].Text
		return text, nil
	} else {
		return "", fmt.Errorf("no text returned by LLM")
	}
}

func (ag *Agent) unmarshable(jsonString string) string {

	text := strings.TrimPrefix(jsonString, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSuffix(text, "```\n")
	text = strings.TrimSpace(text)

	return text
}

func (ag *Agent) startTask() error {

	fmt.Print("Enter task: ")
	input, _ := ag.reader.ReadString('\n')

	prompt := string(input)

	for i := 0; i >= 0; i++ {

		response, err := ag.prompt(prompt)
		if err != nil {
			return err
		}

		state := ag.unmarshable(response)
		prompt, err = ag.processState(state)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if prompt == "exit" {
			break
		}

	}
	return nil
}
