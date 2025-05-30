package agent

func (ag *Agent) read() (string, error) {

	// fmt.Print("Enter task: ")
	input, err := ag.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return input, err
}
