package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gCyrille/xoctx/internal/lab"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("111"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("105"))

	listHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	activeIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("108"))

	inactiveIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	activeNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("108"))

	inactiveNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	badgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("108")).
			Italic(true)

	// Shared color palette
	cPurple = lipgloss.Color("170")
	cGreen  = lipgloss.Color("108")
	cDim    = lipgloss.Color("245")

	padStyle = lipgloss.NewStyle().Padding(1)
)

// LabTheme returns a custom huh theme that matches the list/show styling.
var LabTheme = func() *huh.Theme {
	theme := huh.ThemeBase()

	// Title / header
	theme.Group.Title = lipgloss.NewStyle().
		Bold(true).Foreground(cPurple)
	theme.Blurred.Title = theme.Focused.Title

	// Remove the left border ("┃") from the select field
	theme.Focused.Base = lipgloss.NewStyle().PaddingLeft(1)
	theme.Blurred.Base = lipgloss.NewStyle().PaddingLeft(1)

	// Use "●" instead of ">" as the selection cursor
	theme.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(cGreen).
		SetString("● ")
	theme.Blurred.SelectSelector = theme.Focused.SelectSelector

	// Selected option
	theme.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(cGreen).Bold(true)
	theme.Blurred.SelectedOption = theme.Focused.SelectedOption

	// Unselected option
	theme.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(cDim)
	theme.Blurred.UnselectedOption = theme.Focused.UnselectedOption

	theme.Form.Base = lipgloss.NewStyle().Padding(1)

	return theme
}()

func RenderSummary(profile *lab.Profile) string {
	var lines []string

	lines = append(lines, titleStyle.Render(fmt.Sprintf("Lab: %s", profile.Name)))

	keys := make([]string, 0, len(profile.Vars))
	for k := range profile.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := profile.Vars[k]
		keyStr := keyStyle.Render(k)
		var valStr string
		if lab.IsSensitive(k) {
			valStr = valueStyle.Render(maskValue(v))
		} else {
			valStr = valueStyle.Render(v)
		}
		lines = append(lines, fmt.Sprintf("  %s %s", keyStr, valStr))
	}

	return strings.Join(lines, "\n")
}

func maskValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}

func RenderList(labs []string, current string) string {
	if len(labs) == 0 {
		return listHeaderStyle.Render("No contexts configured yet.\n\nCreate a .env file in ~/.xoctx/")
	}

	var lines []string
	lines = append(lines, listHeaderStyle.Render("Available XO Contexts"))

	for _, l := range labs {
		if l == current {
			lines = append(lines, fmt.Sprintf("  %s %s %s",
				activeIconStyle.Render("●"),
				activeNameStyle.Render(l),
				badgeStyle.Render("(active)"),
			))
		} else {
			lines = append(lines, fmt.Sprintf("  %s %s",
				inactiveIconStyle.Render("○"),
				inactiveNameStyle.Render(l),
			))
		}
	}

	return strings.Join(lines, "\n")
}
