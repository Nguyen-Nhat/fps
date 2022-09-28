package jiratest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetReportedIssueLinks(t *testing.T) {
	result := getReportedIssueLinks("  ")
	assert.Equal(t, []string(nil), result)

	result = getReportedIssueLinks("  pes-123")
	assert.Equal(t, []string{"PES-123"}, result)

	result = getReportedIssueLinks("  pes-123, , pes-556")
	assert.Equal(t, []string{"PES-123", "PES-556"}, result)

	result = getReportedIssueLinks("  pes-123, pes-454, PES-666  ")
	assert.Equal(t, []string{"PES-123", "PES-454", "PES-666"}, result)
}

func TestIsIssueLinksFiltered(t *testing.T) {
	var result bool
	result = isIssueLinksFiltered([]string{"PES-123", "PES-567"}, []string{})
	assert.False(t, result)

	result = isIssueLinksFiltered([]string{"PES-123", "PES-567"}, []string{"PES-555"})
	assert.True(t, result)

	result = isIssueLinksFiltered([]string{"PES-123", "PES-555"}, []string{"PES-555"})
	assert.False(t, result)

	result = isIssueLinksFiltered([]string{"PES-123", "PES-555", "PES-567"}, []string{"PES-555", "PES-123"})
	assert.False(t, result)

	result = isIssueLinksFiltered([]string{"PES-123", "PES-555", "PES-567"}, []string{"PES-666", "PES-222"})
	assert.True(t, result)
}

func TestNormalizeIssueLinks(t *testing.T) {
	result := normalizeIssueLinks([]string{" pes-235", " pes-666 "})
	assert.Equal(t, []string{"PES-235", "PES-666"}, result)
}
