package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
)

type modelStats struct {
	ModelName           string                  `json:"modelName"`
	TotalItems          int                     `json:"totalItems"`
	ApprovedItems       int                     `json:"approvedItems"`
	ApprovalRate        float64                 `json:"approvalRate"`
	AvgScore            float64                 `json:"avgScore"`
	AvgWords            float64                 `json:"avgWords"`
	LanguageStats       map[string]languageStat `json:"languageStats"`
	LevelStats          map[string]levelStat    `json:"levelStats"`
	TopRejectionReasons []rejectionReason       `json:"topRejectionReasons"`
	TopTopics           []topicStat             `json:"topTopics"`
}

func (ms *modelStats) asMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*Model Statistics: %s*\n\n", ms.ModelName))
	sb.WriteString("📊 *Key Metrics:*\n")
	sb.WriteString(fmt.Sprintf("• Total dictionaries: %d\n", ms.TotalItems))
	sb.WriteString(fmt.Sprintf("• Approved: %d (%.1f%%)\n", ms.ApprovedItems, ms.ApprovalRate))
	sb.WriteString(fmt.Sprintf("• Average score: %.1f\n", ms.AvgScore))
	sb.WriteString(fmt.Sprintf("• Average word count: %.1f\n\n", ms.AvgWords))

	sb.WriteString("📚 *By Difficulty Level:*\n")

	for _, level := range types.AllLanguageLevels() {
		levelStr := level.String()
		if stat, exists := ms.LevelStats[levelStr]; exists && stat.TotalItems > 0 {
			sb.WriteString(fmt.Sprintf("• %s: %.1f%% approved (%d of %d), avg. score: %.1f\n",
				levelStr, stat.ApprovalRate, stat.ApprovedItems, stat.TotalItems, stat.AvgScore))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("🌐 *Top Language Pairs:*\n")
	type langStatPair struct {
		lang string
		stat languageStat
	}
	var langPairs []langStatPair
	for lang, stat := range ms.LanguageStats {
		langPairs = append(langPairs, langStatPair{lang, stat})
	}
	sort.Slice(langPairs, func(i, j int) bool {
		return langPairs[i].stat.TotalItems > langPairs[j].stat.TotalItems
	})

	for i := 0; i < min(len(langPairs), 5); i++ {
		pair := langPairs[i]
		sb.WriteString(fmt.Sprintf("• %s: %.1f%% approved (%d of %d), avg. score: %.1f\n",
			pair.lang, pair.stat.ApprovalRate, pair.stat.ApprovedItems, pair.stat.TotalItems, pair.stat.AvgScore))
	}
	sb.WriteString("\n")

	sb.WriteString("📌 *Top Topics:*\n")
	for _, topic := range ms.TopTopics {
		sb.WriteString(fmt.Sprintf("• %s: %.1f%% approved (%d of %d), avg. score: %.1f\n",
			topic.Topic, topic.ApprovalRate, topic.ApprovedItems, topic.TotalItems, topic.AvgScore))
	}
	sb.WriteString("\n")

	if len(ms.TopRejectionReasons) > 0 {
		sb.WriteString("❌ *Top Rejection Reasons:*\n")
		for _, reason := range ms.TopRejectionReasons {
			percentage := 0.0
			if (ms.TotalItems - ms.ApprovedItems) > 0 {
				percentage = float64(reason.Count) / float64(ms.TotalItems-ms.ApprovedItems) * 100
			}
			sb.WriteString(fmt.Sprintf("• %s: %d (%.1f%%)\n", reason.Reason, reason.Count, percentage))
		}
	}
	sb.WriteString("\n----------------------------\n")
	return sb.String()
}

type languageStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovalRate  float64 `json:"approvalRate"`
	AvgScore      float64 `json:"avgScore"`
}

type levelStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovalRate  float64 `json:"approvalRate"`
	AvgScore      float64 `json:"avgScore"`
}

type rejectionReason struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

type topicStat struct {
	Topic         string  `json:"topic"`
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovalRate  float64 `json:"approvalRate"`
	AvgScore      float64 `json:"avgScore"`
}

func generateModelStats(items []applingoprocessing.SchemaItem) []modelStats {
	modelStatsMap := make(map[string]*modelStats)

	for _, item := range items {
		modelParts := strings.Split(item.PromptCraft, "::")
		if len(modelParts) != 2 {
			continue
		}

		promptName := modelParts[0]
		stats, exists := modelStatsMap[promptName]
		if !exists {
			stats = &modelStats{
				ModelName:           promptName,
				LanguageStats:       make(map[string]languageStat),
				LevelStats:          make(map[string]levelStat),
				TopTopics:           []topicStat{},
				TopRejectionReasons: []rejectionReason{},
			}
			modelStatsMap[promptName] = stats
		}

		stats.TotalItems++
		stats.AvgScore += float64(item.Score)
		stats.AvgWords += float64(item.Words)

		if item.Upload == 1 {
			stats.ApprovedItems++
		}

		langStat, exists := stats.LanguageStats[item.Languages]
		if !exists {
			langStat = languageStat{}
		}
		langStat.TotalItems++
		langStat.AvgScore += float64(item.Score)
		if item.Upload == 1 {
			langStat.ApprovedItems++
		}
		stats.LanguageStats[item.Languages] = langStat

		lvlStat, exists := stats.LevelStats[item.Level]
		if !exists {
			lvlStat = levelStat{}
		}
		lvlStat.TotalItems++
		lvlStat.AvgScore += float64(item.Score)
		if item.Upload == 1 {
			lvlStat.ApprovedItems++
		}
		stats.LevelStats[item.Level] = lvlStat

		if item.Upload == 0 && item.Reason != "" {
			found := false
			for i, reason := range stats.TopRejectionReasons {
				if reason.Reason == item.Reason {
					stats.TopRejectionReasons[i].Count++
					found = true
					break
				}
			}
			if !found {
				stats.TopRejectionReasons = append(stats.TopRejectionReasons, rejectionReason{
					Reason: item.Reason,
					Count:  1,
				})
			}
		}

		foundTopic := false
		for i, topic := range stats.TopTopics {
			if topic.Topic == item.Topic {
				stats.TopTopics[i].TotalItems++
				stats.TopTopics[i].AvgScore += float64(item.Score)
				if item.Upload == 1 {
					stats.TopTopics[i].ApprovedItems++
				}
				foundTopic = true
				break
			}
		}
		if !foundTopic {
			stats.TopTopics = append(stats.TopTopics, topicStat{
				Topic:      item.Topic,
				TotalItems: 1,
				AvgScore:   float64(item.Score),
				ApprovedItems: func() int {
					if item.Upload == 1 {
						return 1
					}
					return 0
				}(),
			})
		}
	}

	for _, stats := range modelStatsMap {
		if stats.TotalItems > 0 {
			stats.ApprovalRate = float64(stats.ApprovedItems) / float64(stats.TotalItems) * 100
			stats.AvgScore /= float64(stats.TotalItems)
			stats.AvgWords /= float64(stats.TotalItems)
		}

		for lang, langStat := range stats.LanguageStats {
			if langStat.TotalItems > 0 {
				langStat.ApprovalRate = float64(langStat.ApprovedItems) / float64(langStat.TotalItems) * 100
				langStat.AvgScore /= float64(langStat.TotalItems)
				stats.LanguageStats[lang] = langStat
			}
		}

		for level, levelStat := range stats.LevelStats {
			if levelStat.TotalItems > 0 {
				levelStat.ApprovalRate = float64(levelStat.ApprovedItems) / float64(levelStat.TotalItems) * 100
				levelStat.AvgScore /= float64(levelStat.TotalItems)
				stats.LevelStats[level] = levelStat
			}
		}

		for i, topic := range stats.TopTopics {
			if topic.TotalItems > 0 {
				stats.TopTopics[i].ApprovalRate = float64(topic.ApprovedItems) / float64(topic.TotalItems) * 100
				stats.TopTopics[i].AvgScore /= float64(topic.TotalItems)
			}
		}

		sort.Slice(stats.TopTopics, func(i, j int) bool {
			return stats.TopTopics[i].TotalItems > stats.TopTopics[j].TotalItems
		})
		if len(stats.TopTopics) > 5 {
			stats.TopTopics = stats.TopTopics[:5]
		}

		sort.Slice(stats.TopRejectionReasons, func(i, j int) bool {
			return stats.TopRejectionReasons[i].Count > stats.TopRejectionReasons[j].Count
		})
		if len(stats.TopRejectionReasons) > 3 {
			stats.TopRejectionReasons = stats.TopRejectionReasons[:3]
		}
	}

	result := make([]modelStats, 0, len(modelStatsMap))
	for _, stats := range modelStatsMap {
		result = append(result, *stats)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ModelName < result[j].ModelName
	})
	return result
}
