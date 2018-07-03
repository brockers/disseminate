package main

import "testing"
import "errors"

// func TestWarn(t *testing.T){
// 	defer func() {
// 		if r := recover(); r == nil {
// 			fmt.Println("Roveredin f", r)
// 			// t.Errorf("The warn function did not panic")
// 		}
// 	}()
// 	// Check out warn and see if it runs
// 	warn("Error message on panic.")
// 	t.Errorf("Code id not panic")
// 	// Output:
// 	// ERROR: Error message on panic.
// }

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
	// Output:
	// ERROR: Error message on panic.
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
