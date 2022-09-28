package jiratest

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestRunScript ...
type TestRunScript struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Detail ...
type Detail struct {
	Name            string
	IssueLinks      []string
	Objective       string
	Precondition    string
	WebLinks        []string
	Folder          string
	ConfluenceLinks []string
	Steps           []string
}

type testResult struct {
	Name               string        `json:"name"`
	IssueLinks         []string      `json:"issueLinks"`
	Objective          string        `json:"objective"`
	Precondition       string        `json:"precondition"`
	WebLinks           []string      `json:"webLinks"`
	Folder             string        `json:"folder"`
	ConfluenceLinks    []string      `json:"confluenceLinks"`
	TestScript         TestRunScript `json:"testScript"`
	TestRunStatus      string        `json:"testrun_status"`
	TestRunEnvironment string        `json:"testrun_environment"`
	TestRunComment     string        `json:"testrun_comment"`
	TestRunDuration    float64       `json:"testrun_duration"`
	TestRunDate        string        `json:"testrun_date"`
}

var outputFile *os.File
var outputFileMut sync.Mutex

const directoryEnv = "JIRA_PWD"
const reportedIssueLinksEnv = "REPORTED_ISSUE_LINKS"

func getReportedIssueLinks(env string) []string {
	list := strings.Split(env, ",")
	var result []string
	for _, e := range list {
		s := strings.TrimSpace(e)
		if s == "" {
			continue
		}
		s = strings.ToUpper(s)
		result = append(result, s)
	}
	return result
}

func isIssueLinksFiltered(issueLinks []string, reported []string) bool {
	if len(reported) == 0 {
		return false
	}
	for _, s := range reported {
		for _, issueLink := range issueLinks {
			if s == issueLink {
				return false
			}
		}
	}
	return true
}

func normalizeIssueLinks(links []string) []string {
	var result []string
	for _, link := range links {
		s := strings.TrimSpace(link)
		s = strings.ToUpper(s)
		result = append(result, s)
	}
	return result
}

func writeResult(result testResult) {
	directory := os.Getenv(directoryEnv)
	if directory == "" {
		return
	}

	outputFileMut.Lock()
	defer outputFileMut.Unlock()

	if outputFile == nil {
		filename := path.Join(directory, "testrun.tmp.json")
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		outputFile = file
	}

	err := json.NewEncoder(outputFile).Encode(result)
	if err != nil {
		panic(err)
	}
}

// Step adds a step to the Test Script
func (detail *Detail) Step(step string) {
	detail.Steps = append(detail.Steps, step)
}

func stepsToTestScript(steps []string) TestRunScript {
	return TestRunScript{
		Type: "PLAIN_TEXT",
		Text: strings.Join(steps, "</br>"),
	}
}

func getTestEnvironment() string {
	osEnv := os.Getenv("ENV")
	if osEnv != "" {
		return osEnv
	}
	return "local"
}

// Setup set up and tear down a Functional Test Case
// The usage MUST looks like *defer detail.Setup(t)()*
func (detail *Detail) Setup(t *testing.T) func() {
	start := time.Now()

	return func() {
		normalizedLinks := normalizeIssueLinks(detail.IssueLinks)
		if isIssueLinksFiltered(normalizedLinks, getReportedIssueLinks(os.Getenv(reportedIssueLinksEnv))) {
			return
		}

		name := detail.Name
		if name == "" {
			name = t.Name()
		}

		d := time.Since(start)

		status := "Pass"
		if t.Failed() {
			status = "Fail"
		}

		result := testResult{
			Name:               name,
			IssueLinks:         normalizedLinks,
			Objective:          detail.Objective,
			Precondition:       detail.Precondition,
			WebLinks:           detail.WebLinks,
			Folder:             detail.Folder,
			ConfluenceLinks:    detail.ConfluenceLinks,
			TestScript:         stepsToTestScript(detail.Steps),
			TestRunStatus:      status,
			TestRunEnvironment: getTestEnvironment(),
			TestRunDuration:    float64(d.Milliseconds()) / 1000.0,
			TestRunDate:        start.Format(time.RFC3339),
		}

		writeResult(result)
	}
}
