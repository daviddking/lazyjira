package tui

import (
	"testing"

	"github.com/textfuel/lazyjira/v2/pkg/config"
	"github.com/textfuel/lazyjira/v2/pkg/internal/testkit"
	"github.com/textfuel/lazyjira/v2/pkg/jira/jiratest"
)

func TestApp_SideWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		width      int
		configured int
		want       int
	}{
		{"vertical layout returns full width", 60, 40, 60},
		{"default when unset", 200, 0, 40},
		{"explicit width honored when wide", 200, 50, 50},
		{"clamped to half the terminal", 200, 150, 100},
		{"shrinks to 35 percent on narrow terminal", 100, 40, 35},
		{"floors at minimum width", 100, 10, 25},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			app := &App{width: testCase.width, cfg: &config.Config{}}
			app.cfg.GUI.SidePanelWidth = testCase.configured

			testkit.AssertEqual(t, "sideWidth", app.sideWidth(), testCase.want)
		})
	}
}

func TestApp_StackedLayoutKeepsListAndDetailVisible(t *testing.T) {
	t.Parallel()
	for _, height := range []int{20, 40, 60} {
		app := newAppWithFake(t, &jiratest.FakeClient{T: t})
		app.cfg.GUI.StackedLayout = true
		app.width = 120
		app.height = height
		app.side = sideLeft
		app.leftFocus = focusIssues
		app.layoutPanels()

		testkit.AssertEqual(t, "stacked panel width", app.panelSideW, app.width)
		issuesH, detailH := app.panelIssuesH, app.panelDetailH
		if issuesH < 6 || issuesH > 20 || detailH < 3 {
			t.Fatalf("height %d: list/detail heights = %d/%d", height, issuesH, detailH)
		}
		total := app.panelStatusH + issuesH + app.panelInfoH + app.panelProjectsH + detailH + app.panelLogH
		testkit.AssertEqual(t, "panels fit terminal", total, height-1)

		app.side = sideRight
		app.layoutPanels()
		testkit.AssertEqual(t, "list height after focus change", app.panelIssuesH, issuesH)
		testkit.AssertEqual(t, "detail height after focus change", app.panelDetailH, detailH)
	}
}
