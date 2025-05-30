package agent

import (
	"GoodBash/pkg/shell"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

type Agent struct {
	ctx     context.Context
	client  *genai.Client
	history []*genai.Content
	chat    *genai.Chat
	reader  *bufio.Reader
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
		// log.Fatalf("error during prompting \n %v",err)
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

func (ag *Agent) processState(state string) (string, error) {

	// fmt.Println("PROCESSIGN STATE")
	// fmt.Println(state)
	var result map[string]interface{}

	err := json.Unmarshal([]byte(state), &result)

	if err != nil {
		fmt.Printf("errorr in unmarshaling ,%v \n", err)
	}

	stateVal, ok := result["type"].(float64)
	if !ok {
		fmt.Println(state)
		fmt.Println(result)
		return "", fmt.Errorf("invalid response generated\n")
	}

	switch stateVal {
	case 0:
		{
			fmt.Printf("here is the plan :")
			plan, ok := result["plan"].([]interface{})
			if !ok {
				return "", fmt.Errorf("no plan provided")
			}

			fmt.Println(plan)
			return "okk", nil
		}
	case 1:
		{
			fmt.Printf("gonna take the following action :")
			action, ok := result["action"].(string)
			if !ok {
				return "", fmt.Errorf("no action provided")
			}
			fmt.Println(action)
			return "okk", nil
		}
	case 2:
		{
			query, ok := result["question"].(string)
			if !ok {
				return "", fmt.Errorf("no question asked")
			}
			fmt.Println(query)
			fmt.Print("Enter response: ")
			input, _ := ag.reader.ReadString('\n')
			// fmt.Printf("You entered: %s\n", input)
			return input, nil
		}
	case 3:
		{
			command, ok := result["command"].(string)
			if !ok {
				return "invalid response", fmt.Errorf("no command provided")
			}
			resp, err := shell.Run_command(string(command))
			if err != nil {
				fmt.Println("error executing command")
			}

			// fmt.Println("command line response  ----|")
			// fmt.Printf("%s", resp)
			if resp == "" {
				fmt.Println("resp ==  '' ")
				return "shell returned empty string", nil
			}
			return resp, nil
		}
	case 4:
		{
			response, ok := result["finalResponse"].(string)
			if !ok {
				return "invalid response", fmt.Errorf("envalid response")
			}
			fmt.Println("------------")
			fmt.Println(response)
			fmt.Println("------------")
			fmt.Println("safe exiting")

		}
	default:
		{
			return "invalid state", nil
		}
	}

	return "exit", nil
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
