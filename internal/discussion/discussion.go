package discussion

import (
	"fmt"

	"github.com/antgrutta/gh-discussions/internal/api"
	"github.com/shurcooL/githubv4"
)

type Discussion struct {
	Id        string
	Number    string
	Author    string
	Title     string
	Body      string
	Category  Category
	Comments  []Comment
	UpdatedAt string
}

type Category struct {
	Id          string
	Name        string
	Slug        string
	Description string
}

type Comment struct {
	Id       string
	Author   string
	Body     string
	IsAnswer bool
}

// Create a new discussion object to be used to hold data
func NewDiscussion() Discussion {
	return Discussion{}
}

// Create an slice of discussions from GitHub data
func GetDiscussions(repoName string, repoOwner string) []Discussion {
	var discussions []Discussion

	var query struct {
		Repository struct {
			Discussions struct {
				PageInfo struct {
					StartCursor     githubv4.String
					EndCursor       githubv4.String
					HasNextPage     githubv4.Boolean
					HasPreviousPage githubv4.Boolean
				}
				Nodes []struct {
					Id       githubv4.ID
					Number   githubv4.Int
					Author   struct{ Login githubv4.String }
					Title    githubv4.String
					Body     githubv4.String
					Category struct {
						Id          githubv4.String
						Name        githubv4.String
						Slug        githubv4.String
						Description githubv4.String
					}
					UpdatedAt githubv4.DateTime
				}
			} `graphql:"discussions(first: $first, after: $after)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":  githubv4.String(repoName),
		"owner": githubv4.String(repoOwner),
		"first": githubv4.Int(25),
		"after": (*githubv4.String)(nil),
	}

	for {
		err := api.Query(&query, variables)
		if err != nil {
			panic(err)
		}

		// Populate the discussions slice with data from the query
		for _, d := range query.Repository.Discussions.Nodes {
			// TODO: Get comments for each discussion
			discussions = append(discussions, Discussion{
				Id:     fmt.Sprint(d.Id),
				Number: fmt.Sprint(d.Number),
				Author: string(d.Author.Login),
				Title:  string(d.Title),
				Body:   string(d.Body),
				Category: Category{
					Id:          string(d.Category.Id),
					Name:        string(d.Category.Name),
					Slug:        string(d.Category.Slug),
					Description: string(d.Category.Description),
				},
			})
		}

		if !query.Repository.Discussions.PageInfo.HasNextPage {
			break
		} else {
			variables["after"] = githubv4.NewString(query.Repository.Discussions.PageInfo.EndCursor)
		}
	}
	return discussions
}

func (d Discussion) ImportDiscussion(repositoryId string) {
	// Create the discussion
	d.Id = createDiscussion(repositoryId, d)

	// Create the comments
	for _, c := range d.Comments {
		c.Id = createComment(d.Id, c)

		// Mark the comment as answered if it is
		if c.IsAnswer {
			markCommentAsAnswered(c.Id)
		}
	}
}

func createDiscussion(repositoryId string, discussion Discussion) string {
	var mutation struct {
		CreateDiscussion struct {
			Discussion struct {
				Id     githubv4.ID
				Number githubv4.Int
			}
		} `graphql:"createDiscussion(input: $input)"`
	}

	input := githubv4.CreateDiscussionInput{
		RepositoryID: githubv4.String(repositoryId),
		Title:        githubv4.String(discussion.Title),
		Body:         githubv4.String(discussion.Body),
		CategoryID:   githubv4.String(discussion.Category.Id),
	}

	err := api.Mutation(&mutation, input)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}

	return fmt.Sprint(mutation.CreateDiscussion.Discussion.Id)
}

func createComment(discussionId string, comment Comment) string {
	var mutation struct {
		AddDiscussionComment struct {
			Comment struct {
				Id githubv4.ID
			}
		} `graphql:"addDiscussionComment(input: $input)"`
	}

	input := githubv4.AddDiscussionCommentInput{
		DiscussionID: githubv4.String(discussionId),
		Body:         githubv4.String(comment.Body),
	}

	err := api.Mutation(&mutation, input)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}

	return fmt.Sprint(mutation.AddDiscussionComment.Comment.Id)
}

func markCommentAsAnswered(commentId string) {
	var mutation struct {
		MarkDiscussionCommentAsAnswer struct {
			Discussion struct {
				Id githubv4.ID
			}
		} `graphql:"markDiscussionCommentAsAnswer(input: $input)"`
	}

	input := githubv4.MarkDiscussionCommentAsAnswerInput{
		ID: githubv4.String(commentId),
	}

	err := api.Mutation(&mutation, input)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}
}
