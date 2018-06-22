package main

import "os"
import "flag"
// import "io/ioutil"
import "os/exec"
import "fmt"
import "regexp"
import "strings"
import "strconv"
import "encoding/json"

type response struct {
		Commit  string `json:"commit"`
		Message string `json:"message"`
		Date    string `json:"timestamp"`
}

func main(){

	// var number string
	var message string
	var hash []string
	var time []string

	// Regular Expression filters
	var authTag = regexp.MustCompile(`Author: .*\n`)
	var dateTag = regexp.MustCompile(`Date:\s+(.*)\n`)
	var commTag = regexp.MustCompile(`commit ([a-z0-9]*)\n`)
	var mer1Tag = regexp.MustCompile(`Merge pull request .*\n`)
	var mer2Tag = regexp.MustCompile(`Merge: .*\n`)

	// Less simpler printing
	p := fmt.Println
	// wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "An integer identifying the number of previous commits to check.")
	flag.Parse()

	// exec.Command requires a single string for third arg.  Combine strings
	s := []string{ "git log -n", strconv.Itoa(*numbPtr) }

	gitlogCmd := exec.Command("bash", "-c", strings.Join(s, "  "))
	gitlogOut, err := gitlogCmd.CombinedOutput()

	if err != nil {
		// log := spew.Sdump(gitlogCmd)
		// p("ERROR:", log)
		p("ERROR", gitlogOut)
		panic(err)
	}

	message = string(gitlogOut)
	// First test is if the most recent update actually has a merge message
	if message == "" {
		p("ERROR: No message was obtained from git log")
		os.Exit(1)
	}

	hash = commTag.FindStringSubmatch(message)
	time = dateTag.FindStringSubmatch(message)

	// Store the message as a hash document
	if hash[1] == "" {
		p(message)
		panic("Was not able to get the hash value for the git commit message. ")
	}

	// Store the datestamp
	if time[1] == "" {
		p(message)
		panic("Was not able to get the date/time value for the git commit message. ")
	}

	is_merge := mer1Tag.MatchString(message)

	if ! is_merge {
		p("ERROR: Last commit was not a merge")
		os.Exit(1)
	} else {

		// Now we remove the author, date, and commit from the message
		message = authTag.ReplaceAllString(message, "")
		message = dateTag.ReplaceAllString(message, "")
		message = commTag.ReplaceAllString(message, "")
		message = mer1Tag.ReplaceAllString(message, "")
		message = mer2Tag.ReplaceAllString(message, "")
		message = strings.TrimSpace(message)

		if message == "" {
			p("ERROR: Commit message is empty.")
			os.Exit(1)
		} else  {
			// p("SUCCESS:", hash[1])
			// p(message)
			res := &response{
				Commit: hash[1],
				Date: time[1],
				Message: message}
			resJson, _ := json.Marshal(res)
			p(string(resJson))
		}
	}
}
