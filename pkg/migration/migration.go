package migration

import (
	"github.com/antgrutta/gh-discussions/internal/data"
	"github.com/antgrutta/gh-discussions/internal/discussion"
	"github.com/antgrutta/gh-discussions/internal/repository"
	"github.com/pterm/pterm"
)

func ImportData(dataFile string, repoName string, importType string) {
	// Get repository data from GitHub
	repo := repository.NewRepository(repoName)

	// Get data based on import type
	var migrationData []data.MigrationDiscussion
	if importType == "stackoverflow" {
		// Get the import data
		migrationData = data.GetSODiscussions(dataFile)
	}

	// Create progressbar
	p, _ := pterm.DefaultProgressbar.WithTotal(len(migrationData)).WithTitle("Creating discussions").Start()

	// Create the discussion and comments
	for _, d := range migrationData {

		disc := discussion.Discussion{
			Title: d.Title,
			Body:  d.Body,
			Category: discussion.Category{
				Id:   repo.GetCategoryId(d.CategorySlug),
				Slug: d.CategorySlug,
			},
		}

		for _, c := range d.Comments {
			disc.Comments = append(disc.Comments, discussion.Comment{
				Body:     c.Body,
				IsAnswer: c.AcceptedAnswer,
			})
		}

		// Import the discussion to GitHub
		disc.ImportDiscussion(repo.Id)

		// Update progressbar
		p.Increment()
	}

}
