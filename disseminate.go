package main

// import "os"
import "flag"
// import "io/ioutil"
import "os/exec"
// import "strings"
import "fmt"
import "regexp"
import "strings"
// import "github.com/brockers/gitlogutil"
// import "github.com/davecgh/go-spew/spew"
import "strconv"

func main(){

	// var number string

	var authTag = regexp.MustCompile(`Author: .*\n`)
	var dateTag = regexp.MustCompile(`Date: .*\n`)
	var commTag = regexp.MustCompile(`commit [a-z0-9]*\n`)
	var message string
	// wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "an int")

	flag.Parse()
	count := *numbPtr

	// message := "Revese argument example"
	cmdmsg := "bash"
	s := []string{ "git log -n", strconv.Itoa(count), "--grep merge" }
	gitmsg := strings.Join(s, "  ")
	// number = "1"

	// fmt.Println("PATH: ", os.Getenv("PATH"))
	// fmt.Println(message, gitlogutil.Reverse(message))
	// number = strconv.Itoa(*numbPtr)
	// fmt.Println("filter: ", *wordPtr)
	// fmt.Println("count: ", number)
	// fmt.Println("arguments: ", flag.Args())
	// , gitlogutil.Reverse(strings.Join(argsWithProg) ))
	// --merges
	// if wordPtr != nil {
	gitlogCmd := exec.Command(cmdmsg, "-c", gitmsg)
	gitlogOut, err := gitlogCmd.CombinedOutput()

	if err != nil {
		// log := spew.Sdump(gitlogCmd)
		// fmt.Println("ERROR:", log)
		fmt.Println("ERROR", gitlogOut)
		panic(err)
	} else {
		message = string(gitlogOut)
	}

  message = authTag.ReplaceAllString(message, "")
  message = dateTag.ReplaceAllString(message, "")
  message = commTag.ReplaceAllString(message, "")
  message = strings.TrimSpace(message)

	fmt.Println(message)
}
