package main

import "os"
import "flag"
// import "io/ioutil"
import "os/exec"
// import "strings"
import "fmt"
import	"github.com/brockers/gitlogutil"

func main(){

	wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "an int")

	flag.Parse()

	message := "Revese argument example"
	cmdmsg := "bash"
	gitmsg := "-c git log"

	fmt.Println("PATH: ", os.Getenv("PATH"))
	fmt.Println(message, gitlogutil.Reverse(message))
	fmt.Println("filter: ", *wordPtr)
	fmt.Println("count: ", *numbPtr)
	fmt.Println("arguments: ", flag.Args())
	// , gitlogutil.Reverse(strings.Join(argsWithProg) ))
	// --merges
	if wordPtr != nil {
		fmt.Println(cmdmsg, gitmsg, "-n", *numbPtr, "--grep", string(*wordPtr))
		gitlogCmd := exec.Command(cmdmsg, gitmsg, "-n", string(*numbPtr), "--grep", string(*wordPtr))
		gitlogOut, err := gitlogCmd.Output()

		if err != nil {
			panic(err)
		}

		fmt.Println(string(gitlogOut))
	} else {
		gitlogCmd := exec.Command(cmdmsg, gitmsg, "-n", string(*numbPtr))
		gitlogOut, err := gitlogCmd.Output()
		if err != nil {
			panic(err)
		}

		fmt.Println(string(gitlogOut))
	}


}
