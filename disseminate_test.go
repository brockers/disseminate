package main

import "testing"
import "errors"
import "os"
import "os/exec"
// import "strings"
import "fmt"
import "io/ioutil"

func TestCheck(t *testing.T){

	var err error
	check(err, "No Check Error Message should print")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code id not panic")
			// t.Errorf("The warn function did not panic")
		}
	}()

	err = errors.New("Something bad happened")
	// Check out warn and see if it runs
	check(err, "Check Error Message")
}

func TestWarn(t *testing.T){
	if os.Getenv("LOTS_OF_CRASHES") == "1" {
		// Check out warn and see if it runs
		fmt.Println("Run Lots of Crashes")
		warn("Panic Error Message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWarn")
	cmd.Env = append(os.Environ(), "LOTS_OF_CRASHES=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Check Exit Message
	gotBytes, _ := ioutil.ReadAll(stdout)
  got := string(gotBytes)
  expected := "ERROR: Panic Error Message\n"
  if got != expected {
  	t.Fatalf("Unexpected warn message. Got %s but should contain %s", got[:len(got)-1], expected)
  }

	// Check Exit status
  err := cmd.Wait()
  if e, ok := err.(*exec.ExitError); !ok || e.Success() {
    t.Fatalf("Process ran with err %v, want exit status 1", err)
  }
}

func TestGetPackage(t *testing.T){

	got := getPackage("./test/package_good.json")

	if got != (PackageJSON{}){
		t.Fatalf("A valid package.json file was not unmarshaled")
	}

	// TODO: Need to find a way to test our warn functions.
}

// func TestGetGitlogMeessage(t *testing.T){
// 	TODO: getGitlogMessages Testing
// }

func TestGetHashString(t *testing.T) {

	cases := []struct{
		in, want string
	}{
		{ "commit a98d9e63902bad87b3d8", "a98d9e63902bad87b3d8" },
		{ "commit XXX", "XXX" },
	}

	for _, c := range cases {
		got := getHashString(c.in)
		if got != c.want {
			t.Errorf("getHashString(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

