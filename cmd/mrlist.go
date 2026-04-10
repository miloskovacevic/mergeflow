package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/miloskovacevic/mergeflow/internal/domain"
	"time"

	//"github.com/miloskovacevic/mergeflow/internal/infrastructure/gitlab"
	"github.com/spf13/cobra"
)

const blockedUser = "renovate"

var (
	title  = color.New(color.FgCyan, color.Bold).SprintFunc()
	ok     = color.New(color.FgGreen, color.Bold).SprintFunc()
	warn   = color.New(color.FgYellow).SprintFunc()
	accent = color.New(color.FgMagenta).SprintFunc()
)

// mrlistCmd represents the hosts command
var mrlistCmd = &cobra.Command{
	Use:   "mrlist",
	Short: "List all open merge requests on your gitlab projects",
	Long: `Helpful to get an overview of all open merge requests on your gitlab projects. It will list the project name, merge request title and the author of the merge request. Also will show how long is the mr open.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		return ListMergeRequestsHandler(ctx, appInstance.Gitlab)
	},
}

func ListMergeRequestsHandler(ctx context.Context, provider domain.GitProvider) error {
	gitlabProjectIds := appInstance.Config.GitLab.Repos

	var totalOpenMRs int
	var olderThan2Days int
	var olderThan7Days int
	var oldestMRAgeDays int
	var oldestMRID string
	var hasOldest bool

	for _, projectId := range gitlabProjectIds {

		collections, err := provider.GetMergeRequests(ctx, domain.MergeRequestFilter{
			ProjectId: projectId,
			Status:    "opened",
		})
		if err != nil {
			panic(err)
		}

		var openMrCount int = 0
		name := companies[projectId]
		if name == "" {
			name = "unknown"
		}

		fmt.Println(accent("📦 Project:"), name)
		fmt.Println(accent("🆔 Project Id:"), projectId)
		fmt.Println(title("🚀 Merge Requests"))

		for _, mr := range collections {
			if mr.Author != blockedUser {
				ageDays := int(time.Since(mr.CreatedAt).Hours() / 24)
				totalOpenMRs++

				if ageDays > 2 {
					olderThan2Days++
				}
				if ageDays > 7 {
					olderThan7Days++
				}
				if !hasOldest || ageDays > oldestMRAgeDays {
					hasOldest = true
					oldestMRAgeDays = ageDays
					oldestMRID = mr.ID
				}

				fmt.Println(accent("🆔 MR Id:"), mr.ID)
				fmt.Println(accent("📌 Title:"), mr.Title)
				fmt.Println(accent("👤 Author:"), mr.Author)
				fmt.Println(accent("📅 Created:"), mr.CreatedAt.Format("2006-01-02 15:04"))
				fmt.Println(warn("📊 Status:"), mr.Status)
				fmt.Println("────────────────────────")
				openMrCount++
			}
		}
		fmt.Println(ok("✅ Open merge requests count:"), openMrCount)
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	fmt.Println(title("📊 MergeFlow Summary"))
	fmt.Println("")
	fmt.Println(accent("Total Open MRs:"), totalOpenMRs)
	fmt.Println(warn("⚠️ Older than 2 days:"), olderThan2Days)
	fmt.Println(warn("🔥 Older than 7 days:"), olderThan7Days)
	fmt.Println("")
	if hasOldest {
		fmt.Printf("%s %d days (#%s)\n", accent("Oldest MR:"), oldestMRAgeDays, oldestMRID)
	} else {
		fmt.Println(accent("Oldest MR:"), "n/a")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(mrlistCmd)
}
