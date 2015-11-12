package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"
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
	flagPath string
	cDir     string
)

// for sort ModTime().Unix() format
type byName []os.FileInfo

func (f byName) Len() int {
	return len(f)
}
func (f byName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byName) Less(i, j int) bool {
	return f[j].ModTime().Unix() < f[i].ModTime().Unix()
}

func init() {
	cmdAdd.Flag.StringVar(&flagPath, "p", getCurrentDir(), "Add the specification path to the zg data")
}

func runAdd(cmd *Command, args []string) {
	// z.sh original format sample:
	// path|rank|TimeUnix
	// e.g.: /Users/zchee/go/src/github.com/zchee/zg|1252.96|1447342208

	// Parse flag, get current directory path
	if len(flagPath) == 0 {
		cDir = getCurrentDir()
	} else {
		cDir = flagPath
	}

	f, _ := getDataFile()
	fd, err := os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	lf := getLastModifyFile(cDir)

	r := frecent(lf, 4290)
	fmt.Println(r)

	// StringToSlicePointerByte
	// Non memory copy when cast string to []byte
	// cpvec := *(*[]byte)(unsafe.Pointer(&cDir))

	// Write unsafe []byte
	// if _, err = fd.Write(cpvec); err != nil {
	// 	log.Fatal(err)
	// }
}

func ReadDirSortByUnixTime(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Sort(byName(list))
	return list, nil
}

func getCurrentDir() string {
	// TODO: Which is faster?
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, filename, _, _ := runtime.Caller(1)
	// f, err = os.Open(filepath.Join(filepath.Dir(filename), ""))
	// fmt.Printf("pc: %s\nfile:%s\nline:%s\nok:%s", pc, filename, line, ok)
	cDir, _ := os.Getwd()

	return cDir + "/"
}

func getModifyTime(fd *os.File) time.Time {
	mt, _ := fd.Stat()

	return mt.ModTime()
}

func getFile(cDir string) {
	// path := os.Args[1]
	// files, _ := ioutil.ReadDir(path)
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// }
	cp := &cDir
	fmt.Println(cDir)
	files, _ := ReadDirSortByUnixTime(*cp)
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

func frecent(fd int64, rank int64) int64 {
	dx := fd - time.Now().Unix()

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

func getLastModifyFile(cDir string) int64 {
	var err error

	lastModFile := make([]os.FileInfo, 1)
	lastModFile, err = ReadDirSortByUnixTime(cDir)
	if err != nil {
		log.Fatal(err)
	}

	return lastModFile[0].ModTime().Unix()
}
