package agent

import (
	"fmt"
	"strconv"
	"strings"
)

type Registers struct {
	planCounter   int
	actionCounter int
}

func (rg *Registers) getRegisters() []int {
	var arr []int
	arr = append(arr, rg.planCounter)
	arr = append(arr, rg.actionCounter)
	return arr
}

func (rg *Registers) setRegisters(arr []int) {
	rg.planCounter = arr[0]
	rg.actionCounter = arr[1]
}

func (ag *Agent) loadRegisters() error {
	var temp []string
	var temp1 []int
	ag.dbs.registersDB.read(&temp)

	for _, str := range temp {
		str = strings.TrimSpace(str)
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return fmt.Errorf("error loading registers from db")
		}
		temp1 = append(temp1, int(val))
	}
	ag.registers.setRegisters(temp1)
	return nil
}

func (ag *Agent) dumpRegisters() error {
	temp := ""
	arr := ag.registers.getRegisters()

	for _, val := range arr {
		temp += strconv.Itoa(val)
		temp += "/n"
	}

	temp = strings.TrimSuffix(temp, "\n")
	ag.dbs.registersDB.write(temp)

	return nil
}
