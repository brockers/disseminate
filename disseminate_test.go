package main

import "testing"

func TestgetHashString(t *testing.T) {

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
