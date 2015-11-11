package main

import (
	"log"
	"os"
	"unsafe"
)

var cmdAdd = &Command{
	Usage: "add [-p path]",
	Short: "Add directory path into $_ZG_DATA. default is current path.",
	Long: `
Write the directory path into $_ZG_DATA. Default is current path.
If you specify a path use "-p" flag, It will write its path instead current path.
`,
	Run: runAdd,
}

var (
	flagPath    string
	currentPath string
)

func init() {
	cmdAdd.Flag.StringVar(&flagPath, "p", getCurrentPath(), "Add the specification path to the zg data")
}

func runAdd(cmd *Command, args []string) {
	f, _ := getDataFile()

	if len(flagPath) == 0 {
		currentPath = getCurrentPath()
	} else {
		currentPath = flagPath
	}
	cpvec := *(*[]byte)(unsafe.Pointer(&currentPath))

	fd, err := os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	if _, err = fd.Write(cpvec); err != nil {
		log.Fatal(err)
	}
}

func getCurrentPath() string {
	// TODO: Which is faster?
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, filename, _, _ := runtime.Caller(1)
	// f, err = os.Open(filepath.Join(filepath.Dir(filename), ""))
	// fmt.Printf("pc: %s\nfile:%s\nline:%s\nok:%s", pc, filename, line, ok)
	currentPath, _ := os.Getwd()

	return currentPath + "\n"
}
