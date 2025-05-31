package agent

import (
	shell "GoodBash/pkg/Shell"
	"encoding/json"
	"fmt"
)

type STATE int

const (
	EXIT STATE = iota
	PLAN
	ACTION
	QUESTION
	COMMAND

	MEM_REGISTERS
	MEM_PLAN
	MEM_ACTION
)

func (ag *Agent) processState(state string) (string, error) {

	var result map[string]interface{}
	err := json.Unmarshal([]byte(state), &result)

	if err != nil {
		fmt.Printf("errorr in unmarshaling ,%v \n", err)
	}
	stateVal, ok := result["state"].(float64)
	if !ok {
		fmt.Println(state)
		fmt.Println(result)
		return "", fmt.Errorf("invalid response generated")
	}

	switch STATE(stateVal) {
	case PLAN:
		{
			fmt.Printf("here is the plan :")
			plan, ok := result["plan"].([]interface{})
			if !ok {
				return "", fmt.Errorf("no plan provided")
			}
			fmt.Println(plan)
			return "okk", nil
		}
	case ACTION:
		{
			fmt.Printf("gonna take the following action :")
			action, ok := result["action"].(string)
			if !ok {
				return "", fmt.Errorf("no action provided")
			}
			fmt.Println(action)
			return "okk", nil
		}
	case QUESTION:
		{
			query, ok := result["question"].(string)
			if !ok {
				return "", fmt.Errorf("no question asked")
			}
			fmt.Println(query)
			fmt.Print("Enter response: ")
			input, _ := ag.reader.ReadString('\n')
			return input, nil
		}
	case COMMAND:
		{
			command, ok := result["command"].(string)
			if !ok {
				return "invalid response", fmt.Errorf("no command provided")
			}
			resp, err := shell.Run_command(string(command))
			if err != nil {
				fmt.Println("error executing command")
			}
			if resp == "" {
				fmt.Println("resp ==  '' ")
				return "shell returned empty string", nil
			}
			return resp, nil
		}
	case EXIT:
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
