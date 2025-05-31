package agent

import (
	"bufio"
	"os"
	"path/filepath"
)

type Memory struct {
	plan      []string
	actions   []string
	registers Registers
}

type DBs struct {
	planDB      db
	actionsDB   db
	contextDB   db
	registersDB db
}

func NewDBset(DBdir string) DBs {
	err := os.MkdirAll(DBdir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	plandbpath := filepath.Join(DBdir, "planDB.dbfile")
	registerdbpath := filepath.Join(DBdir, "registerDB.dbfile")
	return DBs{
		planDB:      NewplanDB(plandbpath),
		registersDB: NewplanDB(registerdbpath),
	}
}

type db interface {
	read(mem *[]string)
	write(text string)
	append(text string)
}

type planDB struct {
	filepath string
}

func NewplanDB(filepath string) *planDB {

	return &planDB{
		filepath: filepath,
	}
}

func (pdb *planDB) read(mem *[]string) {

	f, err := os.Open(pdb.filepath)
	defer f.Close()

	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	*mem = nil
	for scanner.Scan() {
		*mem = append(*mem, scanner.Text())
	}

}
func (pdb *planDB) write(data string) {
	err := os.WriteFile(pdb.filepath, []byte(data), 0644)
	if err != nil {
		panic(err)
	}

}

func (pdb *planDB) append(text string) {
	//TODO appending logic
}
