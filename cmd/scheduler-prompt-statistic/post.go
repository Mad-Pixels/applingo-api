package main

import (
	"fmt"
	"strings"
	"time"
)

const MaxPostLength = 4000

type Post struct {
	Content string
}

func GeneratePosts(stats []modelStats) []Post {
	var posts []Post
	var currentPostContent strings.Builder

	currentDate := time.Now().Format("01/02/2006")
	header := fmt.Sprintf("*Model Quality Report - %s*\n\n", currentDate)

	totalItems := 0
	totalApproved := 0
	modelCount := len(stats)

	for _, stat := range stats {
		totalItems += stat.TotalItems
		totalApproved += stat.ApprovedItems
	}

	overallApprovalRate := 0.0
	if totalItems > 0 {
		overallApprovalRate = float64(totalApproved) / float64(totalItems) * 100
	}

	summary := "📈 *Overall Statistics:*\n"
	summary += fmt.Sprintf("• Number of models: %d\n", modelCount)
	summary += fmt.Sprintf("• Total dictionaries: %d\n", totalItems)
	summary += fmt.Sprintf("• Approved: %d (%.1f%%)\n\n", totalApproved, overallApprovalRate)
	summary += "👇 *Detailed Statistics by Model*\n\n"

	currentPostContent.WriteString(header)
	currentPostContent.WriteString(summary)

	for _, stat := range stats {
		modelMarkdown := stat.asMarkdown()

		if currentPostContent.Len()+len(modelMarkdown) > MaxPostLength {
			posts = append(posts, Post{Content: currentPostContent.String()})

			currentPostContent.Reset()
			currentPostContent.WriteString(fmt.Sprintf("*Model Quality Report - %s (continued)*\n\n", currentDate))
		}

		currentPostContent.WriteString(modelMarkdown)
	}

	if currentPostContent.Len() > 0 {
		posts = append(posts, Post{Content: currentPostContent.String()})
	}

	return posts
}
