package main

import "testing"

var testAdd = &Command{
	Usage: "add [-p path]",
	Short: "Add directory path into $_ZG_DATA. default is current path.",
	Long: `
Write the directory path into $_ZG_DATA. Default is current path.
If you specify a path use "-p" flag, It will write its path instead current path.
`,
	Run: runAdd,
}

func TestAdd(t *testing.T) {
	// var x interface{} = 7 // x has dynamic type int and value 7
	// i := x.(int)          // i has type int and value 7
	//
	// type I interface {
	// 	m()
	// }
	// var y I
	// s := y.(string)    // illegal: string does not implement I (missing method m)
	// r := y.(io.Reader) // r has type io.Reader and y must implement both I and io.Reader

	for i := 0; i < 1000; i++ {
		var f interface{} = getLastModifyFile("./")
		t.Log(f.(int64))
	}
	// expected := 30
	// if actual != expected {
	// 	t.Errorf("got %v\nwant %v", actual, expected)
	// }
}
