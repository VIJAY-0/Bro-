package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"

	"GoodBash/pkg/shell"
)

// type Response struct {
// 	respType  int `type`
// 	map[string]interface{}
// }

func main() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "AIzaSyCt4xwCe7rolx-7oPh7JOeOHPTT4IquifI",
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Fatal(err)
	}

	BasePrompt := `You are an assistant designed to convert user tasks into a series of executable steps and shell commands.
Your responsibilities:
1. Understand the user's input task.
2. Generate a step-by-step plan.
3. Ask concise questions if more information is needed — only one question per response.
4. Convert steps into executable shell commands.
5. Indicate task completion.

Each response must be in following structure:

------------------------------------------------------------
RESPONSE FORMAT (as string):

{ "type": 0, "plan": ["step1", "step2", ...], "req": ["requirement1", "requirement2", ...] }
  -> Type 0: Planning response

{ "type": 1, "action": "Describe a non-command action" }
  -> Type 1: Action step

{ "type": 2, "question": "Ask a follow-up question" }
  -> Type 2: Clarification question

{ "type": 3, "command": "Shell command to execute" }
  -> Type 3: Executable shell command

{ "type": 4 }
  -> Type 4: Task completion

------------------------------------------------------------
RULES:
- Only ask one question at a time using type 2.
- Wait for user input before continuing.
- Respond with type 4 to end the task.

------------------------------------------------------------
EXAMPLE:

User: list all directories here

Response 1:
{ "type": 0, "plan": ["List all directories using the ls command"], "req": ["Access to the ls command", "Permission to execute shell commands"] }

Response 2:
{ "type": 3, "command": "ls -d */" }

------------------------------------------------------------

Respond strictly in a plain string that follows this format(not a json object). Always use double quotes for keys and values.
`

	history := []*genai.Content{
		genai.NewContentFromText(BasePrompt, genai.RoleUser),
	}

	var prompt string
	// fmt.Scan(&command)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter task: ")
	input, _ := reader.ReadString('\n')
	prompt = strings.TrimSpace(input)
	// prompt := "create a k8s cluster with 2 master and 5 worker nodes ."

	chat, _ := client.Chats.Create(ctx, "gemini-2.0-flash", nil, history)

	for i := 0; i >= 0; i++ {

		// fmt.Scan(&command)
		res, _ := chat.SendMessage(ctx, genai.Part{Text: prompt})

		if len(res.Candidates) > 0 {
			text := res.Candidates[0].Content.Parts[0].Text
			text = strings.TrimPrefix(text, "```json")
			text = strings.TrimPrefix(text, "```")
			text = strings.TrimSuffix(text, "```")
			text = strings.TrimSpace(text)

			fmt.Println(text)

			var result map[string]interface{}

			err = json.Unmarshal([]byte(text), &result)
			if err != nil {
				fmt.Printf("errorr in unmarshaling ,%v \n", err)
			}

			//question
			if val, ok := result["type"].(float64); ok && int(val) == 2 {
				fmt.Print("Enter text: ")
				input, _ := reader.ReadString('\n')
				prompt = strings.TrimSpace(input)
				fmt.Printf("You entered: %s\n", input)

				//exit
			} else if val, ok := result["type"].(float64); ok && int(val) == 4 {
				fmt.Println("safe exiting")
				break

				//command
			} else if val, ok := result["type"].(float64); ok && int(val) == 3 {

				command, ok := result["command"].(string)

				if !ok {

				}
				resp, err := shell.Run_command(string(command))
				// print(resp)

				if err != nil {
					fmt.Println("error executing command")
				}

				fmt.Println("command line response  ----|")
				fmt.Println("command line response      V")
				fmt.Printf("%s", resp)
				if resp == "" {
					fmt.Println("resp ==  '' ")
					prompt = "done"
				}

				fmt.Println("----------------------")
				// prompt = resp
			}

		} else {
			break
		}

	}
}
