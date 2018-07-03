package main

import "testing"
import "errors"
import "os"
import "os/exec"

func TestCheck(t *testing.T){

	var err error
	check(err, "No messsage should print")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code id not panic")
			// t.Errorf("The warn function did not panic")
		}
	}()

	err = errors.New("Something bad happened")
	// Check out warn and see if it runs
	check(err, "Error message on panic.")
}

func TestWarn(t *testing.T){
	if os.Getenv("LOTS_OF_CRASHES") == "1" {
		// Check out warn and see if it runs
		warn("Error message on panic.")
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWarn")
	cmd.Env = append(os.Environ(), "LOTS_OF_CRASHES=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start; err != nil {
		t.Errorf("Code id not panic")
	}

  got := string(stdout)
  expected := "ERROR: Error message on panic."
  if !strings.HasSuffix(got[:len(got)-1], expected) {
    t.Fatalf("Unexpected warn message. Got %s but should contain %s", got[:len(got)-1], expected)
  }
}

// func TestGetHashString(t *testing.T) {
//
// 	cases := []struct{
// 		in, want string
// 	}{
// 		{ "commit a98d9e63902bad87b3d8", "a98d9e63902bad87b3d8" },
// 		{ "commit XXX", "XXX" },
// 	}
//
// 	for _, c := range cases {
// 		got := getHashString(c.in)
// 		if got != c.want {
// 			t.Errorf("getHashString(%q) == %q, want %q", c.in, got, c.want)
// 		}
// 	}
	// Output:
	// a98d9e63902bad87b3d8
	// XXX
// }
