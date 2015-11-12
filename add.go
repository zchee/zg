package main

import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
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
	fd, err := os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	if len(flagPath) == 0 {
		currentPath = getCurrentPath()
	} else {
		currentPath = flagPath
	}
	cpvec := *(*[]byte)(unsafe.Pointer(&currentPath))

	if _, err = fd.Write(cpvec); err != nil {
		log.Fatal(err)
	}

	// TODO: Which is faster?
	// ioutil.ReadDir(getCurrentPath())
	// or
	// ioutil.ReadDir(string(cpvec))

	// fileAll, _ := ioutil.ReadDir(currentPath)

	// sort.Sort(sort.Reverse(fileTime))
	// fmt.Println(fileTime)

	// files, _ := ioutil.ReadDir(currentPath)
	// fmt.Println("files: ", files)
	// fmt.Println("cpvec: ", cpvec)
	// fmt.Println("string cpvec: ", string(cpvec))

	// files, _ := ioutil.ReadDir("/Users/zchee/go/src/github.com/zchee/zg")
	// r := frecent(fd, getModifyTime(fd).Unix())
	// fmt.Println(r)
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

	return currentPath + "/"
}

func getModifyTime(fd *os.File) time.Time {
	mt, _ := fd.Stat()

	return mt.ModTime()
}

func getFileList(currentPath string) {
	// path := os.Args[1]
	// files, _ := ioutil.ReadDir(path)
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// }
	cp := &currentPath
	fmt.Println(currentPath)
	files, _ := ioutil.ReadDir(*cp)
	// files, _ := ioutil.ReadDir("/Users/zchee/go/src/github.com/zchee/zg")
	fmt.Println(files)

	// ModifyTimeInt64 slice
	// mti := make([]int64, 0)
	// for _, f := range files {
	// 	fmt.Println("filename", f.Name())
	// 	fmt.Println("file time", f.ModTime().Unix())
	// 	mti = append(mti, f.ModTime().Unix())
	// }
	// fmt.Println("CurrentPath: ", cp)
	// fmt.Println("ModifyTimeInt64 slice", mti)
}

func frecent(fd *os.File, rank int64) int64 {
	dx := getModifyTime(fd).Unix() - time.Now().Unix()

	switch true {

	// Accessed less than an hour ago
	case dx < 3:
		fmt.Println("an hour ago")
		return rank * 4

	// Accessed less than an day ago
	case dx < 86400:
		fmt.Println("an day ago")
		return rank * 2

	// Accessed less than an week ago
	case dx < 604800:
		fmt.Println("an week ago")
		return rank / 2

	// Accessed more than an week ago
	default:
		fmt.Println("too old")
		return rank / 4
	}
}
