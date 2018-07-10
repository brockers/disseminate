package main

import "os"
import "os/exec"
import "flag"
import "time"
import "fmt"
import "regexp"
import "strings"
import "strconv"
import "encoding/json"
import "io/ioutil"
import "github.com/dghubble/oauth1"
import "gopkg.in/russross/blackfriday.v2"
// import "github.com/russross/blackfriday"
// import "github.com/microcosm-cc/bluemonday"

// Regular Expression filters
var authTag = regexp.MustCompile(`Author: .*\n`)
var dateTag = regexp.MustCompile(`Date:\s+(.*)\n`)
var mer1Tag = regexp.MustCompile(`Merge pull request .*\n`)
var mer2Tag = regexp.MustCompile(`Merge: .*\n`)
var commTag = regexp.MustCompile(`commit ([a-z0-9]*)\n`)
var leadSpace = regexp.MustCompile(`(?m)^[\t ]{2,}`)

type Response struct {
	Commit  string `json:"commit"`
	Message string `json:"message"`
	Date    string `json:"timestamp"`
	Time    string `json:"buildtime"`
	Disseminate PackageDisseminate `json:"package"`
}

type WpPost struct {
	Title string `json:"title"`
	Status string `json:"status"`
	Catagory int `json:"categories"`
	Content string `json:"content"`
	Comment string `json:"comment_status"`
	Format string `json:"format"`
}

type PackageDisseminate struct {
	Product string `json:"product"`
	Website string `json:"website"`
}

type PackageJSON struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	License     string `json:"license"`
	Private     bool   `json:"private"`
	Disseminate PackageDisseminate `json:"disseminate"`
}

// Standarize our error message
func check(e error, s string) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

// No panic but we want a exit with a message
// os.Exit is hard to test. By assigning it a var name I can mock it in tests
var warn = func(s string) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", s)
	os.Exit(1)
}

// Get the package file and convert it from json
func getPackage(f string) PackageJSON {
	raw, err := ioutil.ReadFile(f)
	check(err,"Unable to open configuration file")

	var c PackageJSON
	json.Unmarshal(raw, &c)

	if c == (PackageJSON{}) {
		warn("Bad JSON document or non-npm style json file type.")
	}

	if c.Disseminate == (PackageDisseminate{}) {
		warn("No disseminate tag in configuration file. { disseminate: {} } is required.")
	}

	return c
}

// Log messages from git
func getGitlogMessage(n string) string {

	// exec.Command requires a single string for third arg.  Combine strings
	s := []string{ "git log -n", n }

	gitlogCmd := exec.Command("bash", "-c", strings.Join(s, "  "))
	gitlogOut, err := gitlogCmd.CombinedOutput()
	check(err, string(gitlogOut))

	m := string(gitlogOut)
	// First test is if the most recent update actually has a merge message
	if m == "" {
		warn("No message was obtained from git log")
	}

	return m
}

// Grab has from our commit message,
func getHashString(m string) string {

	var hsh []string

	hsh = commTag.FindStringSubmatch(m)
	// Store the message as a hsh document
	if hsh[1] == "" {
		warn("Was not able to get the hash value for the git commit message. ")
	}

	return hsh[1]
}

// Get time from commit message
func getCommitTime(m string) string {

	var cTime []string

	cTime = dateTag.FindStringSubmatch(m)
	// If we do not have one get very upset
	if cTime[1] == "" {
		warn("Was not able to get the date/time value for the git commit message. ")
	}

	return cTime[1]
}

// See if it is a merge, and it a merge is required.
func checkMergeRequirement(m string, is_nMrg bool) {

	var has_merge bool

	has_merge = mer1Tag.MatchString(m)

	if ! has_merge && ! is_nMrg {
		warn("Last commit was not a merge")
	}
}

// Convert our json into an actual string
func (p PackageDisseminate) toString() string {
	return toJson(p)
}

// Decode json document into native go structure
func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	check(err, "Failed to decide json document.")

	return string(bytes)
}

func main(){

	// Message Conditions
	var maxLenth int = 280
	var minLenth int = 30
	var is_noMerge bool = false
	var is_print bool = false
	var is_markdown bool = false
	var configFile string = "./package.json"
	var saveFile string = "./disseminate.json"

	// Wordpress Post Options
	var is_post bool = false
	var wp_status string = "publish"
	var wp_category int = 11
	var wp_comment string = "closed"
	var wp_format string = "standard"

	// var number string
	var message string
	var hash string
	var commitTime string

	// oAuth1 values
	postClientKey := os.Getenv("D_OAUTH_CLIENT_KEY")
	postClientSecret := os.Getenv("D_OAUTH_CLIENT_SECRET")
	postToken := os.Getenv("D_OAUTH_TOKEN")
	postTokenSecret := os.Getenv("D_OAUTH_TOKEN_SECRET")
	postUrl := os.Getenv("D_POST_URL")

	// Setup OAuth1
	config := oauth1.NewConfig(postClientKey,postClientSecret)
	token := oauth1.NewToken(postToken,postTokenSecret)

	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, token)

	// Build our timestamp
	now := time.Now()
	timestamp := []string{
		strconv.Itoa(now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		strconv.Itoa(now.Day()),
		strconv.Itoa(now.Hour()),
		strconv.Itoa(now.Minute()) }

	// Less simpler printing
	p := fmt.Println

	// wordPtr := flag.String("filter", "merge", "a string")
	numbPtr := flag.Int("count", 1, "Number of previous commit messages to return.")
	maxLenthPtr := flag.Int("max", maxLenth, "Maximum allowable length of the commit message.")
	minLenthPtr := flag.Int("min", minLenth, "Minimum allowable length of the commit message.")
	is_noMergePtr := flag.Bool("nomerge", is_noMerge, "Include non-merge commits in results.")
	is_printPtr := flag.Bool("print", is_print, "Print output to terminal instead of to a file.")
	is_postPtr := flag.Bool("post", is_post, "Post to a RESTful Endpoint. You will need a number of ENV setup to do this.")
	is_markdownPtr := flag.Bool("markdown", is_markdown, "Support Markdown formatting in the gitlog message and convert to HTML for Wordpress.")
	configFilePtr := flag.String("config", configFile, "Disseminate configuration file in JSON format.")
	saveFilePtr := flag.String("save", saveFile, "Save output to a file.  Cannot be used with -print.")
	flag.Parse()
	// Reset our defaults to new imports
	is_post = *is_postPtr
	minLenth = *minLenthPtr
	maxLenth = *maxLenthPtr
	is_print = *is_printPtr
	is_markdown = *is_markdownPtr

	// Get out package information
	pkgs := getPackage(*configFilePtr)
	// Get our package messages
	message = getGitlogMessage(strconv.Itoa(*numbPtr))

	// Hash values for uniq storage
	hash = getHashString(message)
	// Need different build and commit times
	commitTime = getCommitTime(message)
	// If we have  a merge requirement let us know
	checkMergeRequirement(message, *is_noMergePtr)

	// Now we remove the author, date, and commit from the message
	message = authTag.ReplaceAllString(message, "")
	message = dateTag.ReplaceAllString(message, "")
	message = commTag.ReplaceAllString(message, "")
	message = mer1Tag.ReplaceAllString(message, "")
	message = mer2Tag.ReplaceAllString(message, "")
	message = strings.TrimSpace(message)
	message = leadSpace.ReplaceAllString(message, "")

	if is_markdown {
		// p("====Before Markdown====")
		// p(message)
		// unsafe := blackfriday.Run([]byte(message))
		// message = string(blackfriday.MarkdownBasic([]byte(message)))
		message = string(blackfriday.Run([]byte(message), blackfriday.WithNoExtensions()))
		// message = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		// p("====After Markdown====")
		// p(message)
	}

	// Message tests, to make sure they have a certain "quality"
	if message == "" {
		warn("Commit message is empty.")
	} else if len(message) <= minLenth {
		warn("Commit message does not meet minimum allowable length.")
	} else if len(message) >= maxLenth {
		warn("Commit message exceeds the maximum allowable length.")
	} else  {
		res := &Response{
			Commit: hash,
			Date: commitTime,
			Time: strings.Join(timestamp, ""),
			Message: message,
			Disseminate: pkgs.Disseminate}
		resJson, _ := json.Marshal(res)

		// Post to a RESTFUL API
		if is_post {

			post := &WpPost{
				Title: "New Update to " + pkgs.Disseminate.Product,
				Status: wp_status,
				Catagory:wp_category,
				Comment: wp_comment,
				Format: wp_format,
				Content: message}
			postJson, _ := json.Marshal(post)

			resp, err := httpClient.Post(postUrl,  "application/json", strings.NewReader(string(postJson)))
			defer resp.Body.Close()

			// Check for bad post
			check(err, "Unable to make a successful HTTP POST.")
			p("Wordpress Post Successful")

			body, _ := ioutil.ReadAll(resp.Body)
			p(string(body))

		}

		// Print or write.
		if is_print {
			p(string(resJson))
		} else {
			err := ioutil.WriteFile(*saveFilePtr, resJson, 0644)
			check(err, "Unable to write to save file")
			p((*saveFilePtr), "has been updated")
		}
	}

}
