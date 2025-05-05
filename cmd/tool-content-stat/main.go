package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	defaultDaysWatchdog = 7
)

// ReportData stores all necessary data for the report
type ReportData struct {
	Date          string                          `json:"date"`
	Period        string                          `json:"period"`
	TotalItems    int                             `json:"totalItems"`
	ModelStats    map[string]ModelStat            `json:"modelStats"`
	LevelStats    map[string]map[string]LevelStat `json:"levelStats"`
	TopicStats    map[string]map[string]TopicStat `json:"topicStats"`
	LanguageStats map[string]map[string]LangStat  `json:"languageStats"`
	TopReasons    map[string]map[string]int       `json:"topReasons"`
	TimelineData  map[string][]TimePoint          `json:"timelineData"`
}

// ModelStat statistics by model
type ModelStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovedRate  float64 `json:"approvedRate"`
	AvgScore      float64 `json:"avgScore"`
	AvgWords      float64 `json:"avgWords"`
}

// LevelStat statistics by difficulty level
type LevelStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovedRate  float64 `json:"approvedRate"`
	AvgScore      float64 `json:"avgScore"`
}

// TopicStat statistics by topic
type TopicStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovedRate  float64 `json:"approvedRate"`
	AvgScore      float64 `json:"avgScore"`
}

// LangStat statistics by language pair
type LangStat struct {
	TotalItems    int     `json:"totalItems"`
	ApprovedItems int     `json:"approvedItems"`
	ApprovedRate  float64 `json:"approvedRate"`
	AvgScore      float64 `json:"avgScore"`
}

// TimePoint for timeline data
type TimePoint struct {
	Date         string  `json:"date"`
	Count        int     `json:"count"`
	ApprovedRate float64 `json:"approvedRate"`
	AvgScore     float64 `json:"avgScore"`
}

// ReasonCount for sorting rejection reasons
type ReasonCount struct {
	Reason string
	Count  int
}

// TopicCount for sorting topics
type TopicCount struct {
	Topic string
	Count int
}

// LanguageCount for sorting language pairs
type LanguageCount struct {
	Language string
	Count    int
}

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dynamo := cloud.NewDynamo(cfg)
	table := applingoprocessing.TableName

	// Get all records
	items, err := fetchAllItems(ctx, dynamo, table)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Filter items for specified period
	now := time.Now()
	startDate := now.AddDate(0, 0, -defaultDaysWatchdog)

	var periodItems []applingoprocessing.SchemaItem
	for _, item := range items {
		itemTime := time.Unix(int64(item.Created), 0)
		if itemTime.After(startDate) {
			periodItems = append(periodItems, item)
		}
	}

	// Generate report data
	period := fmt.Sprintf("Last %d days", defaultDaysWatchdog)
	reportData := generateReportData(periodItems, now, period)

	// Generate HTML report
	htmlReport := generateHTMLReport(reportData)

	// Save report to file
	reportFileName := fmt.Sprintf("applingo_report_%s.html", now.Format("2006-01-02"))
	err = os.WriteFile(reportFileName, []byte(htmlReport), 0644)
	if err != nil {
		log.Fatalf("Error saving report: %v", err)
	}

	fmt.Printf("Report successfully saved to file: %s\n", reportFileName)
}

// fetchAllItems gets all items from DynamoDB table
func fetchAllItems(ctx context.Context, dynamo *cloud.Dynamo, table string) ([]applingoprocessing.SchemaItem, error) {
	var allItems []applingoprocessing.SchemaItem
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		scanInput := dynamo.BuildScanInput(table, 100, lastEvaluatedKey)
		result, err := dynamo.Scan(ctx, table, scanInput)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		var items []applingoprocessing.SchemaItem
		err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %v", err)
		}

		allItems = append(allItems, items...)

		if result.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = result.LastEvaluatedKey
	}

	return allItems, nil
}

// generateReportData generates data for the report
func generateReportData(items []applingoprocessing.SchemaItem, now time.Time, period string) ReportData {
	report := ReportData{
		Date:          now.Format("2006-01-02"),
		Period:        period,
		TotalItems:    len(items),
		ModelStats:    make(map[string]ModelStat),
		LevelStats:    make(map[string]map[string]LevelStat),
		TopicStats:    make(map[string]map[string]TopicStat),
		LanguageStats: make(map[string]map[string]LangStat),
		TopReasons:    make(map[string]map[string]int),
		TimelineData:  make(map[string][]TimePoint),
	}

	// Maps for counting topics and languages by model
	topicCounts := make(map[string]map[string]int)
	langCounts := make(map[string]map[string]int)

	// Maps for timeline data
	dailyData := make(map[string]map[string][]applingoprocessing.SchemaItem)

	// Process items
	for _, item := range items {
		// Extract model name from PromptCraft (format: craftName::model)
		modelParts := strings.Split(item.PromptCraft, "::")
		if len(modelParts) != 2 {
			continue
		}

		promptName := modelParts[0]

		// Update model statistics
		stats, exists := report.ModelStats[promptName]
		if !exists {
			stats = ModelStat{}
		}

		stats.TotalItems++
		stats.AvgScore += float64(item.Score)
		stats.AvgWords += float64(item.Words)

		if item.Upload == 1 {
			stats.ApprovedItems++
		}

		report.ModelStats[promptName] = stats

		// Collect rejection reasons
		if item.Upload == 0 {
			if _, exists := report.TopReasons[promptName]; !exists {
				report.TopReasons[promptName] = make(map[string]int)
			}
			report.TopReasons[promptName][item.Reason]++
		}

		// Statistics by difficulty level for each model
		if _, exists := report.LevelStats[promptName]; !exists {
			report.LevelStats[promptName] = make(map[string]LevelStat)
		}

		levelStat, exists := report.LevelStats[promptName][item.Level]
		if !exists {
			levelStat = LevelStat{}
		}

		levelStat.TotalItems++
		levelStat.AvgScore += float64(item.Score)

		if item.Upload == 1 {
			levelStat.ApprovedItems++
		}

		report.LevelStats[promptName][item.Level] = levelStat

		// Count topics for future top-5 selection
		if _, exists := topicCounts[promptName]; !exists {
			topicCounts[promptName] = make(map[string]int)
		}
		topicCounts[promptName][item.Topic]++

		// Prepare topic statistics
		if _, exists := report.TopicStats[promptName]; !exists {
			report.TopicStats[promptName] = make(map[string]TopicStat)
		}

		topicStat, exists := report.TopicStats[promptName][item.Topic]
		if !exists {
			topicStat = TopicStat{}
		}

		topicStat.TotalItems++
		topicStat.AvgScore += float64(item.Score)

		if item.Upload == 1 {
			topicStat.ApprovedItems++
		}

		report.TopicStats[promptName][item.Topic] = topicStat

		// Process language pairs
		if _, exists := langCounts[promptName]; !exists {
			langCounts[promptName] = make(map[string]int)
		}
		langCounts[promptName][item.Languages]++

		if _, exists := report.LanguageStats[promptName]; !exists {
			report.LanguageStats[promptName] = make(map[string]LangStat)
		}

		langStat, exists := report.LanguageStats[promptName][item.Languages]
		if !exists {
			langStat = LangStat{}
		}

		langStat.TotalItems++
		langStat.AvgScore += float64(item.Score)

		if item.Upload == 1 {
			langStat.ApprovedItems++
		}

		report.LanguageStats[promptName][item.Languages] = langStat

		// Process timeline data
		itemTime := time.Unix(int64(item.Created), 0)
		dateKey := itemTime.Format("2006-01-02")

		if _, exists := dailyData[promptName]; !exists {
			dailyData[promptName] = make(map[string][]applingoprocessing.SchemaItem)
		}

		dailyData[promptName][dateKey] = append(dailyData[promptName][dateKey], item)
	}

	// Calculate averages for models
	for name, stats := range report.ModelStats {
		if stats.TotalItems > 0 {
			stats.AvgScore /= float64(stats.TotalItems)
			stats.AvgWords /= float64(stats.TotalItems)
			stats.ApprovedRate = float64(stats.ApprovedItems) / float64(stats.TotalItems) * 100
			report.ModelStats[name] = stats
		}
	}

	// Calculate averages for difficulty levels
	for modelName, levelStats := range report.LevelStats {
		for level, stats := range levelStats {
			if stats.TotalItems > 0 {
				stats.AvgScore /= float64(stats.TotalItems)
				stats.ApprovedRate = float64(stats.ApprovedItems) / float64(stats.TotalItems) * 100
				report.LevelStats[modelName][level] = stats
			}
		}
	}

	// Calculate averages for topics
	for modelName, topicStats := range report.TopicStats {
		for topic, stats := range topicStats {
			if stats.TotalItems > 0 {
				stats.AvgScore /= float64(stats.TotalItems)
				stats.ApprovedRate = float64(stats.ApprovedItems) / float64(stats.TotalItems) * 100
				report.TopicStats[modelName][topic] = stats
			}
		}
	}

	// Calculate averages for language pairs
	for modelName, langStats := range report.LanguageStats {
		for lang, stats := range langStats {
			if stats.TotalItems > 0 {
				stats.AvgScore /= float64(stats.TotalItems)
				stats.ApprovedRate = float64(stats.ApprovedItems) / float64(stats.TotalItems) * 100
				report.LanguageStats[modelName][lang] = stats
			}
		}
	}

	// Generate timeline data
	for modelName, days := range dailyData {
		var timePoints []TimePoint

		// Get sorted dates
		var dates []string
		for date := range days {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		for _, date := range dates {
			items := days[date]
			timePoint := TimePoint{
				Date:  date,
				Count: len(items),
			}

			var totalScore float64
			var approvedCount int

			for _, item := range items {
				totalScore += float64(item.Score)
				if item.Upload == 1 {
					approvedCount++
				}
			}

			if len(items) > 0 {
				timePoint.AvgScore = totalScore / float64(len(items))
				timePoint.ApprovedRate = float64(approvedCount) / float64(len(items)) * 100
			}

			timePoints = append(timePoints, timePoint)
		}

		report.TimelineData[modelName] = timePoints
	}

	return report
}

// generateHTMLReport generates the HTML report
func generateHTMLReport(data ReportData) string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Applingo Model Quality Report - ` + data.Date + `</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@3.7.1/dist/chart.min.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: 20px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .approved {
            color: green;
        }
        .rejected {
            color: red;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .summary {
            display: flex;
            justify-content: space-between;
            flex-wrap: wrap;
            margin-bottom: 20px;
        }
        .summary-card {
            background-color: #f8f9fa;
            border-radius: 5px;
            padding: 15px;
            width: 30%;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            margin-bottom: 10px;
        }
        .summary-card h3 {
            margin-top: 0;
        }
        .summary-card p {
            font-size: 24px;
            font-weight: bold;
            margin: 5px 0;
        }
        .level-good {
            background-color: #d4edda;
        }
        .level-medium {
            background-color: #fff3cd;
        }
        .level-bad {
            background-color: #f8d7da;
        }
        .chart-container {
            position: relative;
            height: 300px;
            margin-bottom: 30px;
        }
        .flex-container {
            display: flex;
            justify-content: space-between;
        }
        .chart-half {
            width: 48%;
        }
        .recommendation {
            background-color: #e2f0fd;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 15px;
        }
        .recommendation h4 {
            margin-top: 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Applingo Model Quality Report</h1>
        <p>Date: ` + data.Date + `</p>
        <p>Period: ` + data.Period + `</p>
        
        <div class="summary">
            <div class="summary-card">
                <h3>Total Dictionaries</h3>
                <p>` + fmt.Sprintf("%d", data.TotalItems) + `</p>
            </div>
        </div>
        
        <h2>1. Model Summary Statistics</h2>
        <table>
            <tr>
                <th>Model</th>
                <th>Total Dictionaries</th>
                <th>Approved</th>
                <th>Approval Rate</th>
                <th>Average Score</th>
                <th>Average Word Count</th>
            </tr>`

	// Sort models by name for stable display
	var modelNames []string
	for name := range data.ModelStats {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	for _, name := range modelNames {
		stats := data.ModelStats[name]
		html += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%.1f%%</td>
                <td>%.1f</td>
                <td>%.1f</td>
            </tr>`, name, stats.TotalItems, stats.ApprovedItems, stats.ApprovedRate, stats.AvgScore, stats.AvgWords)
	}

	html += `
        </table>
        
        <div class="chart-container">
            <canvas id="modelComparisonChart"></canvas>
        </div>
        
        <h2>2. Analysis by Difficulty Level</h2>`

	for _, modelName := range modelNames {
		html += fmt.Sprintf(`
        <h3>%s</h3>
        <table>
            <tr>
                <th>Level</th>
                <th>Total</th>
                <th>Approved</th>
                <th>Approval Rate</th>
                <th>Average Score</th>
            </tr>`, modelName)

		// Sort difficulty levels in order A1, A2, B1, B2, C1, C2
		levelOrder := []string{"A1", "A2", "B1", "B2", "C1", "C2"}

		for _, level := range levelOrder {
			if stats, exists := data.LevelStats[modelName][level]; exists {
				// Determine CSS class for row highlighting
				cssClass := ""
				if stats.ApprovedRate >= 75 {
					cssClass = "level-good"
				} else if stats.ApprovedRate >= 50 {
					cssClass = "level-medium"
				} else {
					cssClass = "level-bad"
				}

				html += fmt.Sprintf(`
            <tr class="%s">
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%.1f%%</td>
                <td>%.1f</td>
            </tr>`, cssClass, level, stats.TotalItems, stats.ApprovedItems, stats.ApprovedRate, stats.AvgScore)
			}
		}

		html += `
        </table>`
	}

	html += `
        <div class="flex-container">
            <div class="chart-half">
                <div class="chart-container">
                    <canvas id="levelApprovalChart"></canvas>
                </div>
            </div>
            <div class="chart-half">
                <div class="chart-container">
                    <canvas id="levelScoreChart"></canvas>
                </div>
            </div>
        </div>
        
        <h2>3. Analysis by Language Pair</h2>`

	for _, modelName := range modelNames {
		html += fmt.Sprintf(`
        <h3>%s</h3>
        <table>
            <tr>
                <th>Language Pair</th>
                <th>Total</th>
                <th>Approved</th>
                <th>Approval Rate</th>
                <th>Average Score</th>
            </tr>`, modelName)

		// Sort languages by count
		var langCounts []LanguageCount
		for lang, stats := range data.LanguageStats[modelName] {
			langCounts = append(langCounts, LanguageCount{lang, stats.TotalItems})
		}

		sort.Slice(langCounts, func(i, j int) bool {
			return langCounts[i].Count > langCounts[j].Count
		})

		for _, langCount := range langCounts {
			lang := langCount.Language
			stats := data.LanguageStats[modelName][lang]

			// Determine CSS class for row highlighting
			cssClass := ""
			if stats.ApprovedRate >= 75 {
				cssClass = "level-good"
			} else if stats.ApprovedRate >= 50 {
				cssClass = "level-medium"
			} else {
				cssClass = "level-bad"
			}

			html += fmt.Sprintf(`
            <tr class="%s">
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%.1f%%</td>
                <td>%.1f</td>
            </tr>`, cssClass, lang, stats.TotalItems, stats.ApprovedItems, stats.ApprovedRate, stats.AvgScore)
		}

		html += `
        </table>`
	}

	html += `
        <div class="chart-container">
            <canvas id="languageChart"></canvas>
        </div>
        
        <h2>4. Analysis by Top-5 Topics</h2>`

	for _, modelName := range modelNames {
		html += fmt.Sprintf(`
        <h3>%s</h3>
        <table>
            <tr>
                <th>Topic</th>
                <th>Total</th>
                <th>Approved</th>
                <th>Approval Rate</th>
                <th>Average Score</th>
            </tr>`, modelName)

		// Sort topics by count
		var topicCounts []TopicCount
		for topic, stats := range data.TopicStats[modelName] {
			topicCounts = append(topicCounts, TopicCount{topic, stats.TotalItems})
		}

		sort.Slice(topicCounts, func(i, j int) bool {
			return topicCounts[i].Count > topicCounts[j].Count
		})

		// Show top-5 topics or fewer if there are less than 5
		limit := 5
		if len(topicCounts) < limit {
			limit = len(topicCounts)
		}

		for i := 0; i < limit; i++ {
			topic := topicCounts[i].Topic
			stats := data.TopicStats[modelName][topic]

			// Determine CSS class for row highlighting
			cssClass := ""
			if stats.ApprovedRate >= 75 {
				cssClass = "level-good"
			} else if stats.ApprovedRate >= 50 {
				cssClass = "level-medium"
			} else {
				cssClass = "level-bad"
			}

			html += fmt.Sprintf(`
            <tr class="%s">
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%.1f%%</td>
                <td>%.1f</td>
            </tr>`, cssClass, topic, stats.TotalItems, stats.ApprovedItems, stats.ApprovedRate, stats.AvgScore)
		}

		html += `
        </table>`
	}

	html += `
        <div class="chart-container">
            <canvas id="topicChart"></canvas>
        </div>
        
        <h2>5. Top Rejection Reasons by Model</h2>`

	// For each model show top-5 rejection reasons
	for _, modelName := range modelNames {
		if reasonCounts, exists := data.TopReasons[modelName]; exists && len(reasonCounts) > 0 {
			html += fmt.Sprintf(`
        <h3>%s</h3>
        <table>
            <tr>
                <th>Rejection Reason</th>
                <th>Count</th>
                <th>%% of Rejected</th>
            </tr>`, modelName)

			// Sort reasons by frequency
			var reasons []ReasonCount
			totalRejected := data.ModelStats[modelName].TotalItems - data.ModelStats[modelName].ApprovedItems

			for reason, count := range reasonCounts {
				reasons = append(reasons, ReasonCount{reason, count})
			}

			sort.Slice(reasons, func(i, j int) bool {
				return reasons[i].Count > reasons[j].Count
			})

			// Show top-5 reasons or fewer if there are less than 5
			limit := 5
			if len(reasons) < limit {
				limit = len(reasons)
			}

			for i := 0; i < limit; i++ {
				percentage := 0.0
				if totalRejected > 0 {
					percentage = float64(reasons[i].Count) / float64(totalRejected) * 100
				}

				html += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%d</td>
                <td>%.1f%%</td>
            </tr>`, reasons[i].Reason, reasons[i].Count, percentage)
			}

			html += `
        </table>`
		} else {
			html += fmt.Sprintf(`
        <h3>%s</h3>
        <p>No rejected dictionaries or all dictionaries approved.</p>`, modelName)
		}
	}

	html += `
        <div class="chart-container">
            <canvas id="rejectionChart"></canvas>
        </div>
        
        <h2>6. Performance Trend Over Time</h2>
        <div class="chart-container">
            <canvas id="timelineChart"></canvas>
        </div>
        
        <h2>7. Comparative Analysis</h2>
        <table>
            <tr>
                <th>Metric</th>`

	// Column headers with model names
	for _, name := range modelNames {
		html += fmt.Sprintf(`
                <th>%s</th>`, name)
	}

	html += `
            </tr>
            <tr>
                <td>Approval Rate</td>`

	for _, name := range modelNames {
		html += fmt.Sprintf(`
                <td>%.1f%%</td>`, data.ModelStats[name].ApprovedRate)
	}

	html += `
            </tr>
            <tr>
                <td>Average Score</td>`

	for _, name := range modelNames {
		html += fmt.Sprintf(`
                <td>%.1f</td>`, data.ModelStats[name].AvgScore)
	}

	html += `
            </tr>
            <tr>
                <td>Average Word Count</td>`

	for _, name := range modelNames {
		html += fmt.Sprintf(`
                <td>%.1f</td>`, data.ModelStats[name].AvgWords)
	}

	html += `
            </tr>
        </table>
        
        <h2>8. Key Findings and Recommendations</h2>`

	// Generate findings and recommendations for each model
	for _, modelName := range modelNames {
		stats := data.ModelStats[modelName]

		html += fmt.Sprintf(`
        <h3>%s</h3>`, modelName)

		// Overall performance assessment
		html += `
        <div class="recommendation">
            <h4>Overall Performance</h4>`

		if stats.ApprovedRate >= 75 {
			html += fmt.Sprintf(`
            <p>Model performance is <strong>strong</strong> with an approval rate of %.1f%%. The average quality score is %.1f.</p>`, stats.ApprovedRate, stats.AvgScore)
		} else if stats.ApprovedRate >= 50 {
			html += fmt.Sprintf(`
            <p>Model performance is <strong>moderate</strong> with an approval rate of %.1f%%. The average quality score is %.1f.</p>`, stats.ApprovedRate, stats.AvgScore)
		} else {
			html += fmt.Sprintf(`
            <p>Model performance is <strong>below target</strong> with an approval rate of only %.1f%%. The average quality score is %.1f.</p>`, stats.ApprovedRate, stats.AvgScore)
		}

		html += `
        </div>`

		// Strength areas (best performing levels and topics)
		html += `
        <div class="recommendation">
            <h4>Strength Areas</h4>
            <p>This model performs best with:</p>
            <ul>`

		// Find best performing levels
		type PerformanceItem struct {
			Name         string
			ApprovedRate float64
		}

		var levelPerformance []PerformanceItem
		for level, levelStats := range data.LevelStats[modelName] {
			if levelStats.TotalItems >= 2 {
				levelPerformance = append(levelPerformance, PerformanceItem{level, levelStats.ApprovedRate})
			}
		}

		sort.Slice(levelPerformance, func(i, j int) bool {
			return levelPerformance[i].ApprovedRate > levelPerformance[j].ApprovedRate
		})

		// Find best performing topics
		var topicPerformance []PerformanceItem
		for topic, topicStats := range data.TopicStats[modelName] {
			if topicStats.TotalItems >= 2 {
				topicPerformance = append(topicPerformance, PerformanceItem{topic, topicStats.ApprovedRate})
			}
		}

		sort.Slice(topicPerformance, func(i, j int) bool {
			return topicPerformance[i].ApprovedRate > topicPerformance[j].ApprovedRate
		})

		// Output best performing levels
		if len(levelPerformance) > 0 {
			limit := 2
			if len(levelPerformance) < limit {
				limit = len(levelPerformance)
			}

			for i := 0; i < limit; i++ {
				level := levelPerformance[i]
				if level.ApprovedRate > 0 {
					html += fmt.Sprintf(`
                <li>Difficulty level <strong>%s</strong> (%.1f%% approval rate)</li>`, level.Name, level.ApprovedRate)
				}
			}
		}

		// Output best performing topics
		if len(topicPerformance) > 0 {
			limit := 2
			if len(topicPerformance) < limit {
				limit = len(topicPerformance)
			}

			for i := 0; i < limit; i++ {
				topic := topicPerformance[i]
				if topic.ApprovedRate > 0 {
					html += fmt.Sprintf(`
                <li>Topic <strong>%s</strong> (%.1f%% approval rate)</li>`, topic.Name, topic.ApprovedRate)
				}
			}
		}

		html += `
            </ul>
        </div>`

		// Improvement areas
		html += `
        <div class="recommendation">
            <h4>Areas for Improvement</h4>`

		// Get top rejection reasons
		if reasonCounts, exists := data.TopReasons[modelName]; exists && len(reasonCounts) > 0 {
			var reasons []ReasonCount
			for reason, count := range reasonCounts {
				reasons = append(reasons, ReasonCount{reason, count})
			}

			sort.Slice(reasons, func(i, j int) bool {
				return reasons[i].Count > reasons[j].Count
			})

			topReason := reasons[0].Reason

			html += fmt.Sprintf(`
            <p>The most common reason for dictionary rejection is: <strong>%s</strong></p>
            <p>Recommendations:</p>
            <ul>`, topReason)

			// Custom recommendations based on typical rejection reasons
			if strings.Contains(strings.ToLower(topReason), "language") {
				html += `
                <li>Review and improve language pair handling in prompt</li>
                <li>Add more explicit language validation checks</li>
                <li>Consider focusing the model on fewer language pairs for better quality</li>`
			} else if strings.Contains(strings.ToLower(topReason), "level") {
				html += `
                <li>Better align content with the specified difficulty level</li>
                <li>Add clearer examples of appropriate content for each level in the prompt</li>
                <li>Consider specialized prompts for different difficulty levels</li>`
			} else if strings.Contains(strings.ToLower(topReason), "topic") || strings.Contains(strings.ToLower(topReason), "off-topic") {
				html += `
                <li>Improve topic adherence in the prompt instructions</li>
                <li>Add more examples of on-topic content</li>
                <li>Implement stronger topic validation criteria</li>`
			} else if strings.Contains(strings.ToLower(topReason), "repetitive") {
				html += `
                <li>Enhance diversity of generated content</li>
                <li>Add instructions to avoid repetition in prompt</li>
                <li>Implement checks for entry uniqueness</li>`
			} else {
				html += `
                <li>Review the prompt structure and examples</li>
                <li>Add more specific quality guidelines</li>
                <li>Consider model fine-tuning to address persistent issues</li>`
			}

			html += `
            </ul>`
		} else {
			html += `
            <p>No specific improvement areas identified based on rejections.</p>`
		}

		html += `
        </div>`

		// Specific difficulty level recommendations
		html += `
        <div class="recommendation">
            <h4>Difficulty Level Strategy</h4>
            <ul>`

		// Check for problematic levels
		problemLevelFound := false
		for _, level := range []string{"A1", "A2", "B1", "B2", "C1", "C2"} {
			if stats, exists := data.LevelStats[modelName][level]; exists && stats.TotalItems >= 2 && stats.ApprovedRate < 50 {
				problemLevelFound = true
				html += fmt.Sprintf(`
                <li>Consider revising prompts for <strong>%s</strong> level (%.1f%% approval rate)</li>`, level, stats.ApprovedRate)
			}
		}

		if !problemLevelFound {
			html += `
                <li>All difficulty levels show reasonable performance</li>`
		}

		html += `
            </ul>
        </div>
        
        <div class="recommendation">
            <h4>Summary of Actions</h4>
            <ol>`

		// Overall recommendations based on model performance
		if stats.ApprovedRate < 50 {
			html += `
                <li>Comprehensive prompt review and revision is recommended</li>
                <li>Consider additional review steps in the generation process</li>
                <li>Focus on improving the most common rejection reasons</li>`
		} else if stats.ApprovedRate < 75 {
			html += `
                <li>Targeted improvements for specific difficulty levels and topics</li>
                <li>Address the top rejection reasons</li>
                <li>Monitor performance trends over time</li>`
		} else {
			html += `
                <li>Maintain current quality standards</li>
                <li>Consider expanding to additional topics or language pairs</li>
                <li>Regularly review for consistency</li>`
		}

		html += `
            </ol>
        </div>`
	}

	// Add JavaScript for charts
	html += `
    </div>

    <script>
        // Helper function to generate colors
        function generateColors(count) {
            const colors = [
                'rgba(54, 162, 235, 0.7)',
                'rgba(255, 99, 132, 0.7)',
                'rgba(75, 192, 192, 0.7)',
                'rgba(255, 159, 64, 0.7)',
                'rgba(153, 102, 255, 0.7)',
                'rgba(255, 205, 86, 0.7)',
                'rgba(201, 203, 207, 0.7)',
                'rgba(255, 99, 71, 0.7)',
                'rgba(46, 139, 87, 0.7)',
                'rgba(106, 90, 205, 0.7)'
            ];
            
            let result = [];
            for (let i = 0; i < count; i++) {
                result.push(colors[i % colors.length]);
            }
            return result;
        }

        // Model Comparison Chart
        const modelCtx = document.getElementById('modelComparisonChart').getContext('2d');
        new Chart(modelCtx, {
            type: 'bar',
            data: {
                labels: [`

	// Add model names as labels for the chart
	for i, name := range modelNames {
		if i > 0 {
			html += ", "
		}
		html += fmt.Sprintf("'%s'", name)
	}

	html += `],
                datasets: [
                    {
                        label: 'Approval Rate (%)',
                        data: [`

	// Add approval rates as data
	for i, name := range modelNames {
		if i > 0 {
			html += ", "
		}
		html += fmt.Sprintf("%.1f", data.ModelStats[name].ApprovedRate)
	}

	html += `],
                        backgroundColor: 'rgba(54, 162, 235, 0.7)',
                        borderColor: 'rgba(54, 162, 235, 1)',
                        borderWidth: 1
                    },
                    {
                        label: 'Average Score',
                        data: [`

	// Add average scores as data
	for i, name := range modelNames {
		if i > 0 {
			html += ", "
		}
		html += fmt.Sprintf("%.1f", data.ModelStats[name].AvgScore)
	}

	html += `],
                        backgroundColor: 'rgba(255, 99, 132, 0.7)',
                        borderColor: 'rgba(255, 99, 132, 1)',
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Model Performance Comparison'
                    },
                    legend: {
                        position: 'top',
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100
                    }
                }
            }
        });

        // Level Approval Chart
        const levelApprovalCtx = document.getElementById('levelApprovalChart').getContext('2d');
        new Chart(levelApprovalCtx, {
            type: 'radar',
            data: {
                labels: ['A1', 'A2', 'B1', 'B2', 'C1', 'C2'],
                datasets: [`

	// Generate datasets for each model for the level approval chart
	for i, name := range modelNames {
		if i > 0 {
			html += ","
		}

		html += fmt.Sprintf(`
                    {
                        label: '%s',
                        data: [`, name)

		// Add approval rates for each level
		for j, level := range []string{"A1", "A2", "B1", "B2", "C1", "C2"} {
			if j > 0 {
				html += ", "
			}

			if stats, exists := data.LevelStats[name][level]; exists {
				html += fmt.Sprintf("%.1f", stats.ApprovedRate)
			} else {
				html += "0"
			}
		}

		html += fmt.Sprintf(`],
                        backgroundColor: 'rgba(%d, %d, %d, 0.2)',
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        borderWidth: 1
                    }`, 54+i*50, 162-i*30, 235-i*50, 54+i*50, 162-i*30, 235-i*50)
	}

	html += `
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Approval Rate by Difficulty Level (%)'
                    }
                },
                scales: {
                    r: {
                        min: 0,
                        max: 100,
                        ticks: {
                            stepSize: 20
                        }
                    }
                }
            }
        });

        // Level Score Chart
        const levelScoreCtx = document.getElementById('levelScoreChart').getContext('2d');
        new Chart(levelScoreCtx, {
            type: 'radar',
            data: {
                labels: ['A1', 'A2', 'B1', 'B2', 'C1', 'C2'],
                datasets: [`

	// Generate datasets for each model for the level score chart
	for i, name := range modelNames {
		if i > 0 {
			html += ","
		}

		html += fmt.Sprintf(`
                    {
                        label: '%s',
                        data: [`, name)

		// Add average scores for each level
		for j, level := range []string{"A1", "A2", "B1", "B2", "C1", "C2"} {
			if j > 0 {
				html += ", "
			}

			if stats, exists := data.LevelStats[name][level]; exists {
				html += fmt.Sprintf("%.1f", stats.AvgScore)
			} else {
				html += "0"
			}
		}

		html += fmt.Sprintf(`],
                        backgroundColor: 'rgba(%d, %d, %d, 0.2)',
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        borderWidth: 1
                    }`, 255-i*50, 99+i*30, 132+i*50, 255-i*50, 99+i*30, 132+i*50)
	}

	html += `
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Average Score by Difficulty Level'
                    }
                },
                scales: {
                    r: {
                        min: 0,
                        max: 100,
                        ticks: {
                            stepSize: 20
                        }
                    }
                }
            }
        });

        // Language Chart
        const langCtx = document.getElementById('languageChart').getContext('2d');
        new Chart(langCtx, {
            type: 'bar',
            data: {
                labels: [`

	// Create a set of all language pairs
	languagePairs := make(map[string]bool)
	for _, modelName := range modelNames {
		for lang := range data.LanguageStats[modelName] {
			languagePairs[lang] = true
		}
	}

	// Convert to a sorted slice
	var uniqueLangs []string
	for lang := range languagePairs {
		uniqueLangs = append(uniqueLangs, lang)
	}
	sort.Strings(uniqueLangs)

	// Add language pairs as labels
	for i, lang := range uniqueLangs {
		if i > 0 {
			html += ", "
		}
		html += fmt.Sprintf("'%s'", lang)
	}

	html += `],
                datasets: [`

	// Generate datasets for each model for the language chart
	for i, name := range modelNames {
		if i > 0 {
			html += ","
		}

		html += fmt.Sprintf(`
                    {
                        label: '%s',
                        data: [`, name)

		// Add approval rates for each language pair
		for j, lang := range uniqueLangs {
			if j > 0 {
				html += ", "
			}

			if stats, exists := data.LanguageStats[name][lang]; exists {
				html += fmt.Sprintf("%.1f", stats.ApprovedRate)
			} else {
				html += "0"
			}
		}

		html += fmt.Sprintf(`],
                        backgroundColor: 'rgba(%d, %d, %d, 0.7)',
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        borderWidth: 1
                    }`, 75+i*50, 192-i*30, 192+i*30, 75+i*50, 192-i*30, 192+i*30)
	}

	html += `
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Approval Rate by Language Pair (%)'
                    },
                    legend: {
                        position: 'top',
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100
                    }
                }
            }
        });

        // Topic Chart
        const topicCtx = document.getElementById('topicChart').getContext('2d');
        
        // Get top 5 topics across all models
        const topTopics = [`

	// Create a map to count occurrences of each topic across all models
	allTopicCounts := make(map[string]int)
	for _, modelName := range modelNames {
		for topic, stats := range data.TopicStats[modelName] {
			allTopicCounts[topic] += stats.TotalItems
		}
	}

	// Convert to a slice and sort
	var sortedTopics []TopicCount
	for topic, count := range allTopicCounts {
		sortedTopics = append(sortedTopics, TopicCount{topic, count})
	}

	sort.Slice(sortedTopics, func(i, j int) bool {
		return sortedTopics[i].Count > sortedTopics[j].Count
	})

	// Get top 5 topics
	topTopicsLimit := 5
	if len(sortedTopics) < topTopicsLimit {
		topTopicsLimit = len(sortedTopics)
	}

	// Add top topics as labels
	for i := 0; i < topTopicsLimit; i++ {
		if i > 0 {
			html += ", "
		}
		html += fmt.Sprintf("'%s'", sortedTopics[i].Topic)
	}

	html += `];
        
        new Chart(topicCtx, {
            type: 'bar',
            data: {
                labels: topTopics,
                datasets: [`

	// Generate datasets for each model for the topic chart
	for i, name := range modelNames {
		if i > 0 {
			html += ","
		}

		html += fmt.Sprintf(`
                    {
                        label: '%s',
                        data: [`, name)

		// Add approval rates for top topics
		for j := 0; j < topTopicsLimit; j++ {
			topic := sortedTopics[j].Topic
			if j > 0 {
				html += ", "
			}

			if stats, exists := data.TopicStats[name][topic]; exists {
				html += fmt.Sprintf("%.1f", stats.ApprovedRate)
			} else {
				html += "0"
			}
		}

		html += fmt.Sprintf(`],
                        backgroundColor: 'rgba(%d, %d, %d, 0.7)',
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        borderWidth: 1
                    }`, 153+i*30, 102+i*40, 255-i*30, 153+i*30, 102+i*40, 255-i*30)
	}

	html += `
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Approval Rate by Top Topics (%)'
                    },
                    legend: {
                        position: 'top',
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100
                    }
                }
            }
        });

        // Rejection Reasons Chart
        const rejectionCtx = document.getElementById('rejectionChart').getContext('2d');
        
        // Get common rejection reasons
        const commonReasons = {};
        
        // Count total occurrences of each reason across all models
        `

	// Loop through models to count rejection reasons
	for _, name := range modelNames {
		html += fmt.Sprintf(`
        if (%t) {
            const reasons_%s = `, len(data.TopReasons[name]) > 0, name)

		// Add rejection reasons as a JavaScript object
		html += "{"
		for reason, count := range data.TopReasons[name] {
			html += fmt.Sprintf("\"%s\": %d, ", strings.ReplaceAll(reason, "\"", "\\\""), count)
		}
		html += "};"

		html += fmt.Sprintf(`
            for (const reason in reasons_%s) {
                if (!commonReasons[reason]) {
                    commonReasons[reason] = 0;
                }
                commonReasons[reason] += reasons_%s[reason];
            }
        }`, name, name)
	}

	html += `
        
        // Convert to arrays and sort
        const reasonEntries = Object.entries(commonReasons).sort((a, b) => b[1] - a[1]);
        const topReasons = reasonEntries.slice(0, 5).map(entry => entry[0]);
        
        new Chart(rejectionCtx, {
            type: 'pie',
            data: {
                labels: topReasons,
                datasets: [
                    {
                        data: reasonEntries.slice(0, 5).map(entry => entry[1]),
                        backgroundColor: generateColors(5),
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Common Rejection Reasons'
                    },
                    legend: {
                        position: 'right',
                    }
                }
            }
        });

        // Timeline Chart
        const timelineCtx = document.getElementById('timelineChart').getContext('2d');
        new Chart(timelineCtx, {
            type: 'line',
            data: {
                datasets: [`

	// Generate datasets for each model for the timeline chart
	for i, name := range modelNames {
		if i > 0 {
			html += ","
		}

		if timePoints, exists := data.TimelineData[name]; exists && len(timePoints) > 0 {
			// Create arrays for dates and values
			html += fmt.Sprintf(`
                    {
                        label: '%s - Approval Rate',
                        data: [`, name)

			// Add date points
			for j, point := range timePoints {
				if j > 0 {
					html += ", "
				}
				html += fmt.Sprintf("{x: '%s', y: %.1f}", point.Date, point.ApprovedRate)
			}

			html += fmt.Sprintf(`],
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        backgroundColor: 'rgba(%d, %d, %d, 0.1)',
                        borderWidth: 2,
                        tension: 0.4,
                        fill: false
                    },
                    {
                        label: '%s - Avg Score',
                        data: [`, 54+i*50, 162-i*30, 235-i*50, 54+i*50, 162-i*30, 235-i*50, name)

			// Add score points
			for j, point := range timePoints {
				if j > 0 {
					html += ", "
				}
				html += fmt.Sprintf("{x: '%s', y: %.1f}", point.Date, point.AvgScore)
			}

			html += fmt.Sprintf(`],
                        borderColor: 'rgba(%d, %d, %d, 1)',
                        backgroundColor: 'rgba(%d, %d, %d, 0.1)',
                        borderWidth: 2,
                        tension: 0.4,
                        fill: false,
                        yAxisID: 'y1'
                    }`, 255-i*50, 99+i*30, 132+i*50, 255-i*50, 99+i*30, 132+i*50)
		}
	}

	html += `
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Performance Trend Over Time'
                    },
                    legend: {
                        position: 'top',
                    }
                },
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'day'
                        }
                    },
                    y: {
                        beginAtZero: true,
                        max: 100,
                        title: {
                            display: true,
                            text: 'Approval Rate (%)'
                        }
                    },
                    y1: {
                        beginAtZero: true,
                        max: 100,
                        position: 'right',
                        title: {
                            display: true,
                            text: 'Average Score'
                        },
                        grid: {
                            drawOnChartArea: false
                        }
                    }
                }
            }
        });
    </script>
</body>
</html>`

	return html
}
