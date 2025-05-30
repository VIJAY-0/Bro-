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

	BasePrompt := `You will be given some task to generate shell e
	you will be given user input and you need to generate a plan and create actions accordingly to complete the task, you can ask question in following format. you need to give response in josn format, no extra text needed.
	ask 1 question per prompt.
	return plain json string nothing else.
	after asking question you wiil be primpted with whatever user input was given.
	after completeion return type exit
			
	The format wiil be of following type.
					{'type':'typeofresponse', .... }

					type can be of following 0->plan,
											 1->action, 
											 2->question,
											 3->command,
											 4->exit.

					{'type':0,'plan': ['do this and this and this' , 'fllowing with this'] , 'req':[req1 , req2 , req3]} }}
					{'type':1,'action': " "}
					{'type':2,'question': "question statement here"}
					{'type':3,'command':"<give a shell executteable command in string here>"}
					example {'type':3,'command':"ls -d /directory_xyz" }
					{'type':4 } 
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
			if strings.HasPrefix(text, "```json") {
				text = strings.TrimPrefix(text, "```json")
				text = strings.TrimSuffix(text, "```")
				text = strings.TrimSpace(text)
			}

			fmt.Println(text)

			var result map[string]interface{}
			err = json.Unmarshal([]byte(text), &result)
			if err != nil {
				fmt.Printf("errorr in unmarshaling ,%v\ns", err)
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

				fmt.Printf("command line response ::  %s \n", resp)
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
