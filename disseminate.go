package main

// import "os"
import "flag"
// import "io/ioutil"
import "os/exec"
import "fmt"
import "regexp"
import "strings"
import "strconv"

func main(){

	// var number string
	var message string
	var hash []string

	// Regular Expression filters
	var authTag = regexp.MustCompile(`Author: .*\n`)
	var dateTag = regexp.MustCompile(`Date: .*\n`)
	var commTag = regexp.MustCompile(`commit ([a-z0-9]*)\n`)
	var mer1Tag = regexp.MustCompile(`Merge pull request .*\n`)
	var mer2Tag = regexp.MustCompile(`Merge: .*\n`)

	// wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "an int")
	flag.Parse()

	// exec.Command requires a single string for third arg.  Combine strings
	s := []string{ "git log -n", strconv.Itoa(*numbPtr) }

	gitlogCmd := exec.Command("bash", "-c", strings.Join(s, "  "))
	gitlogOut, err := gitlogCmd.CombinedOutput()

	if err != nil {
		// log := spew.Sdump(gitlogCmd)
		// fmt.Println("ERROR:", log)
		fmt.Println("ERROR", gitlogOut)
		panic(err)
	}

	message = string(gitlogOut)
	// First test is if the most recent update actually has a merge message
	if message == "" {
		fmt.Println("No message was obtained from git log")
	}

	hash = commTag.FindStringSubmatch(message)
	// Store the message as a hash document
	if hash[1] == "" {
		fmt.Println(message)
		panic("Was not able to get the hash value for the git commit message. ")
	}


	is_merge := mer1Tag.MatchString(message)

	if ! is_merge {
		fmt.Println("Last commit was not a merge")
	} else {

		// Now we remove the author, date, and commit from the message
		message = authTag.ReplaceAllString(message, "")
		message = dateTag.ReplaceAllString(message, "")
		message = commTag.ReplaceAllString(message, "")
		message = mer1Tag.ReplaceAllString(message, "")
		message = mer2Tag.ReplaceAllString(message, "")
		message = strings.TrimSpace(message)

		if message == "" {
			panic("Commit message is empty.")
		} else  {
			fmt.Println("SUCCESS:", hash[1])
			fmt.Println(message)
		}
	}
}
