package main

import agent "GoodBash/pkg/Agent"

// type Response struct {
// 	respType  int `type`
// 	map[string]interface{}
// }

func main() {

	Agent := agent.Agent{}
	Agent.InitDBs("DatabaseROOT")

	Agent.Activate()

}
