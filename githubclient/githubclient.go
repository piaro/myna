package githubclient

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var isOrg = false
var client *github.Client

func init() {
	loadEnv()
	client = auth()
	fmt.Println(client)
	initUser()
}

// Repository type of Github Repository
// used in GraphQL and returning value
type Repository struct {
	ID          int64
	Name        string
	Description string
}

// Issue type of Github Issue
type Issue struct {
	author    string
	body      string
	createdAt string
}

type Input interface{}

type AddIssueInput struct {
	RepositoryID string `json:"repositoryId"`
	Body         string `json:"body"`
	Title        string `json:"title"`
}

// FetchRepository fetch Repository info from github
func FetchRepository(name string) Repository {
	repoName := "tokyocameraclub/" + name
	result, _, err := client.Search.Repositories(context.Background(), repoName, nil)
	if err != nil {
		fmt.Println(err)
	}

	targetRepo := result.Repositories[0]
	answerRepo := Repository{
		ID:          *targetRepo.ID,
		Name:        *targetRepo.Name,
		Description: *targetRepo.Description,
	}

	return answerRepo
}

// FetchRepositories list repositories owned by specified user
func FetchRepositories() []*github.Repository {
	if isOrg {
		return fetchRepositoriesByOrg()
	}
	return fetchRepositoriesByUser()
}

func fetchRepositoriesByOrg() []*github.Repository {
	result, _, err := client.Repositories.ListByOrg(context.Background(), getGithubOwner(), nil)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func fetchRepositoriesByUser() []*github.Repository {
	result, _, err := client.Repositories.List(context.Background(), getGithubOwner(), nil)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func getCurrentUser() *github.User {
	user, _, err := client.Users.Get(context.Background(), getGithubOwner())
	if err != nil {
		fmt.Println(err)
	}
	return user
}

// CreateIssue create issue
func CreateIssue(repoName string) *github.Issue {
	client := auth()
	reqIssue := &github.IssueRequest{
		Title: github.String("test issue"),
		Body:  github.String("- test **body**"),
	}
	issue, result, err := client.Issues.Create(
		context.Background(),
		getGithubOwner(),
		repoName,
		reqIssue,
	)
	fmt.Println(result)
	fmt.Println(err)
	return issue
}

func getGithubOwner() string {
	return os.Getenv("GITHUB_OWNER")
}

func auth() *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func initUser() {
	user := getCurrentUser()
	if *user.Type == "Organization" {
		isOrg = true
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error has occured at loading dotenv")
	}
}

func scaffold() {
	// check & make dir

	// check & make dir

	// fetch repository & save(cache)

	// load repositories
}

func createCacheDir() {
	dir := os.Getenv("GOPATH") + "/.myna/cache"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}
}

func putCacheFile() {

}

func loadRepoFromCache() {

}
