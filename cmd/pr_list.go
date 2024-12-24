package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// PullRequest represents a Bitbucket pull request.
type PullRequest struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Branch    string    `json:"branch"`
	CreatedAt time.Time `json:"created_at"`
}

// Dummy function to simulate fetching pull requests
func fetchPullRequests() []PullRequest {
	// Simulated data
	return []PullRequest{
		{ID: 1, Title: "Fix bug", Branch: "fix/bug", CreatedAt: time.Now().Add(-1 * time.Hour)},
		// Add 30+ items here for testing
	}
}

var jsonOutput bool

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests in the current repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch pull requests (replace this with actual API call)
		prs := fetchPullRequests()

		// Sort pull requests by creation date
		sort.Slice(prs, func(i, j int) bool {
			return prs[i].CreatedAt.After(prs[j].CreatedAt)
		})

		// Limit to 30 results
		if len(prs) > 30 {
			prs = prs[:30]
		}

		if jsonOutput {
			// Output as JSON for bb.nvim
			jsonData, err := json.Marshal(prs)
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output as a table for CLI users
			printPullRequestsTable(prs)
		}
	},
}

// Print pull requests as a table
func printPullRequestsTable(prs []PullRequest) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tBRANCH\tCREATED AT")
	fmt.Fprintln(w, "---\t-----\t------\t----------")

	for _, pr := range prs {
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\n",
			pr.ID,
			pr.Title,
			pr.Branch,
			pr.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	w.Flush()
}

func init() {
	prListCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	prCmd.AddCommand(prListCmd)
}
