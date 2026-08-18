package chat

import (
	"fmt"
	"sort"

	"github.com/Anseltsmm/azkia/internal/config"
	"github.com/Anseltsmm/azkia/internal/message"
	"github.com/Anseltsmm/azkia/internal/session"
	"github.com/Anseltsmm/azkia/internal/tui/styles"
	"github.com/Anseltsmm/azkia/internal/tui/theme"
	"github.com/Anseltsmm/azkia/internal/version"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type SendMsg struct {
	Text        string
	Attachments []message.Attachment
}

type SessionSelectedMsg = session.Session

type SessionClearedMsg struct{}

// NewSessionMsg is sent to clear the current session and start a new one
type NewSessionMsg struct{}

// SlashCommandClosedMsg is sent when the slash command dialog closes so the
// editor can clear the "/" input
type SlashCommandClosedMsg struct{}

type EditorFocusMsg bool

func header(width int) string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		logo(width),
		repo(width),
		"",
		cwd(width),
	)
}

// repoURL is the GitHub repository shown on the welcome screen.
const repoURL = "https://github.com/Anseltsmm/opencode"

// lspsConfigured returns a styled list of configured LSP servers, or an empty
// string when none are configured (so empty sections don't clutter the UI).
func lspsConfigured(width int) string {
	cfg := config.Get()

	// Get LSP names and sort them for consistent ordering
	var lspNames []string
	for name := range cfg.LSP {
		lspNames = append(lspNames, name)
	}
	sort.Strings(lspNames)

	// Nothing to show when no LSPs are configured
	if len(lspNames) == 0 {
		return ""
	}

	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()
	title := ansi.Truncate("LSP Configuration", width, "…")

	lsps := baseStyle.
		Width(width).
		Foreground(t.Primary()).
		Bold(true).
		Render(title)

	var lspViews []string
	for _, name := range lspNames {
		lsp := cfg.LSP[name]
		lspName := baseStyle.
			Foreground(t.Text()).
			Render(fmt.Sprintf("• %s", name))

		cmd := lsp.Command
		cmd = ansi.Truncate(cmd, width-lipgloss.Width(lspName)-3, "…")

		lspPath := baseStyle.
			Foreground(t.TextMuted()).
			Render(fmt.Sprintf(" (%s)", cmd))

		lspViews = append(lspViews,
			baseStyle.
				Width(width).
				Render(
					lipgloss.JoinHorizontal(
						lipgloss.Left,
						lspName,
						lspPath,
					),
				),
		)
	}

	return baseStyle.
		Width(width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lsps,
				lipgloss.JoinVertical(
					lipgloss.Left,
					lspViews...,
				),
			),
		)
}

func logo(width int) string {
	logo := fmt.Sprintf("%s %s", styles.AZKIAIcon, "AZKIA")
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	versionText := baseStyle.
		Foreground(t.TextMuted()).
		Render(version.Version)

	// Brand the logo with the primary color so it stands out on the
	// welcome screen instead of blending in with the regular text.
	logoStyle := baseStyle.
		Bold(true).
		Foreground(t.Primary())

	return logoStyle.
		Width(width).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				logo,
				" ",
				versionText,
			),
		)
}

func repo(width int) string {
	t := theme.CurrentTheme()

	return styles.BaseStyle().
		Foreground(t.TextMuted()).
		Width(width).
		Render(repoURL)
}

func cwd(width int) string {
	cwd := fmt.Sprintf("cwd: %s", config.WorkingDirectory())
	t := theme.CurrentTheme()

	return styles.BaseStyle().
		Foreground(t.TextMuted()).
		Width(width).
		Render(cwd)
}
