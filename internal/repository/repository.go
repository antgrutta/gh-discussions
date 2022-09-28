package repository

import (
	"fmt"
	"strings"

	"github.com/antgrutta/gh-discussions/internal/api"
	"github.com/antgrutta/gh-discussions/internal/discussion"
	"github.com/pterm/pterm"
	"github.com/shurcooL/githubv4"
)

type Repository struct {
	Id                   string
	Name                 string
	Owner                string
	DiscussionCategories []discussion.Category
	Discussions          []discussion.Discussion
}

// Create a new repository object to be used to hold data
func NewRepository(repoName string) Repository {
	repo := Repository{
		Name:  strings.Split(repoName, "/")[1],
		Owner: strings.Split(repoName, "/")[0],
	}

	repo.getData()

	return repo
}

// Get Category Id from Category Slug
func (r *Repository) GetCategoryId(slug string) string {
	var result string
	for _, category := range r.DiscussionCategories {
		if category.Slug == slug {
			result = category.Id
		}
	}

	return result
}

// Update the repository object with data from GitHub
func (r *Repository) getData() {
	var query struct {
		Repository struct {
			Id    githubv4.ID
			Name  githubv4.String
			Owner struct {
				Login githubv4.String
			}
			DiscussionCategories struct {
				Nodes []struct {
					Id          githubv4.ID
					Name        githubv4.String
					Slug        githubv4.String
					Description githubv4.String
				}
			} `graphql:"discussionCategories(first: $first)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":  githubv4.String(r.Name),
		"owner": githubv4.String(r.Owner),
		"first": githubv4.Int(100),
	}

	err := api.Query(&query, variables)
	if err != nil {
		panic(err)
	}

	// Syncing data from GitHub to the repository object
	r.Id = fmt.Sprint(query.Repository.Id)
	r.Name = string(query.Repository.Name)
	r.Owner = string(query.Repository.Owner.Login)

	for _, c := range query.Repository.DiscussionCategories.Nodes {
		r.DiscussionCategories = append(r.DiscussionCategories, discussion.Category{
			Id:          fmt.Sprint(c.Id),
			Name:        string(c.Name),
			Slug:        string(c.Slug),
			Description: string(c.Description),
		})
	}

	r.Discussions = discussion.GetDiscussions(r.Name, r.Owner)
}

// Get all discussions as a slice of string slices
func (r *Repository) DiscussionsToStrings() [][]string {
	var result [][]string
	for _, d := range r.Discussions {
		result = append(result, []string{
			pterm.FgGreen.Sprint("#" + fmt.Sprint(d.Number)),
			d.Title,
			pterm.FgBlue.Sprint(d.Category.Name),
		})
	}

	return result
}
