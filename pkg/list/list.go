package list

import (
	"fmt"

	"github.com/antgrutta/gh-discussions/internal/repository"
	"github.com/pterm/pterm"
)

func ListDiscussions(repoName string) {
	// Get repository data from GitHub
	repo := repository.NewRepository(repoName)

	// Display number
	display := 25
	if len(repo.Discussions) <= 25 {
		display = len(repo.Discussions)
	}

	// Print repository data
	pterm.Println() //spacer
	pterm.Println("Showing " + fmt.Sprint(display) + " of " + fmt.Sprint(len(repo.Discussions)) + " discussions in " + repo.Owner + "/" + repo.Name)
	pterm.Println() //spacer

	// Live ouput using area
	area, _ := pterm.DefaultArea.WithCenter().Start()

	area.Update(pterm.DefaultTable.
		WithHasHeader(false).
		WithSeparator("    ").
		WithData(repo.DiscussionsToStrings()).
		Render())

	area.Stop()
}
