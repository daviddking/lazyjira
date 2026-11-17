package views

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/textfuel/lazyjira/v2/pkg/jira"
)

func TestIssueHyperlinkHost(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ host, want string }{
		{"jira.example.com/jira/", "https://jira.example.com/jira/browse/TEST-1"},
		{" https://jira.example.com/?query=1#fragment ", "https://jira.example.com/browse/TEST-1"},
		{"http://localhost:8080/jira", "http://localhost:8080/jira/browse/TEST-1"},
		{"", ""}, {"https://", ""}, {"file:///tmp/jira", ""},
	} {
		t.Run(tc.host, func(t *testing.T) {
			list := NewIssuesList()
			list.SetJiraHost("https://previous.example.com")
			list.SetJiraHost(tc.host)
			want := "TEST-1"
			if tc.want != "" {
				want = ansi.SetHyperlink(tc.want) + want + ansi.ResetHyperlink()
			}
			if got := list.issueKeyLink("TEST-1"); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
			if got := list.issueKeyLink(""); got != "" {
				t.Fatalf("empty key rendered as %q", got)
			}
		})
	}
}

func TestIssueHyperlinksPreserveLayout(t *testing.T) {
	t.Parallel()
	list := makeFocusedIssuesList([]jira.Issue{{Key: "TEST-1", Summary: "First issue"}, {Key: "TEST-22", Summary: "Second issue"}})
	for _, width := range []int{100, 25, 10} {
		list.SetSize(width, 10)
		for _, selected := range []string{"TEST-1", "TEST-22"} {
			list.SelectByKey(selected)
			list.SetJiraHost("")
			before := ansi.Strip(list.View())
			list.SetJiraHost("jira.example.com")
			output := list.View()
			if ansi.Strip(output) != before {
				t.Fatal("hyperlinks changed visible text or alignment")
			}
			if lipgloss.Width(output) != lipgloss.Width(before) {
				t.Fatal("hyperlinks changed panel width")
			}
			if strings.Count(output, "\x1b]8;;https://") != strings.Count(output, ansi.ResetHyperlink()) {
				t.Fatal("clipping left an unclosed hyperlink")
			}
			if width >= 25 {
				for _, key := range []string{"TEST-1", "TEST-22"} {
					link := ansi.SetHyperlink("https://jira.example.com/browse/"+key) + key + ansi.ResetHyperlink()
					if !strings.Contains(output, link) {
						t.Fatalf("missing hyperlink for %s", key)
					}
				}
			}
		}
	}
}
