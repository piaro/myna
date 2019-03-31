package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/c-bata/go-prompt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"

	client "github.com/myna/githubclient"
)

const myEmoji = "🐦"

var _suggestions = []prompt.Suggest{}

func init() {
	loadEnv()
}

/**
 * - Accept user input and handle function
 */
func main() {
	app := cli.NewApp() // using urfave/cli
	app.Name = "myna, github connection tool for TCC."
	app.Usage = "show description for specified repository"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "new issue",
			Aliases: []string{"ni"},
			Usage:   "add new issue to the Repository(please type Repository's Name)",
			Action: func(context *cli.Context) error {
				createIssue(context)
				return nil
			},
		},
		{
			Name:    "filepath",
			Aliases: []string{"fp"},
			Usage:   "echo file path of cache data",
			Action: func(context *cli.Context) error {
				fmt.Println(os.Getenv("GOPATH"))
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "branch, b",
			Usage: "create branch with issue number",
		},
	}

	app.Action = func(context *cli.Context) error {
		fmt.Printf("c.NArg()        : %+v\n", context.NArg())
		fmt.Printf("c.Args()        : %+v\n", context.Args())
		fmt.Printf("c.Args()        : %+v\n", context.Args())
		fmt.Printf("c.FlagNames     : %+v\n", context.FlagNames())

		showRepository() // Default Action
		return nil
	}

	app.Before = func(c *cli.Context) error {
		initSuggestions()
		return nil
	}

	app.After = func(c *cli.Context) error {
		// nothing to do
		return nil
	}

	app.Run(os.Args)
}

// UPDATE _suggestions
func initSuggestions() {
	repos := client.FetchRepositories()
	_suggestions = []prompt.Suggest{}
	for _, repository := range repos {
		_suggestions = append(_suggestions, prompt.Suggest{
			Text:        *repository.Name,
			Description: "",
		})
	}
}

/* createIssue
 */
func createIssue(context *cli.Context) {
	// if len(context.Args()) < 1 {
	// 	inputError("Please input repository's name after `new issue` (ex. `myna ni my-repo`)")
	// 	return
	// }
	// repo := context.Args()[0]
	repo := suggestRepos()
	issue := client.CreateIssue(repo)
	showMyna("CREATED!!")
	fmt.Println(*issue.HTMLURL)
}

/**
 * Read and Show Repository's Description
 * - read stding
 * - call fetch method
 * - show result
 **/
func showRepository() {
	result := client.FetchRepository(suggestRepos())

	showMyna("あったよ！")
	fmt.Println("Name : " + result.Name)
	fmt.Println("Desc : " + result.Description)
	fmt.Println("ID   : " + strconv.FormatInt(result.ID, 10))
}

func inputError(message string) {
	fmt.Println(myEmoji + " Input Error!")
	fmt.Println("message : " + message)
}

func stdinWithMessage(message string) string {
	fmt.Println(message)
	fmt.Printf(" > ")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	return text
}

func showMyna(message string) {
	mynaAA := `
   .-.
  /'v'\ < ` + message + `
 (/   \)
='="="===<
   |_|
`
	fmt.Println(mynaAA)
}

func suggestRepos() string {
	return prompt.Input(">>> ", completer, prompt.OptionTitle("Repo-Selector"))
}

func completer(in prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(_suggestions, in.GetWordBeforeCursor(), true)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error has occured at loading dotenv")
	}
}
