package dialog

import (
	"strings"

	utilComponents "github.com/Anseltsmm/azkia/internal/tui/components/util"
	"github.com/Anseltsmm/azkia/internal/tui/layout"
	"github.com/Anseltsmm/azkia/internal/tui/styles"
	"github.com/Anseltsmm/azkia/internal/tui/theme"
	"github.com/Anseltsmm/azkia/internal/tui/util"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Command represents a command that can be executed
type Command struct {
	ID          string
	Title       string
	Description string
	Handler     func(cmd Command) tea.Cmd
}

func (ci Command) Render(selected bool, width int) string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	descStyle := baseStyle.Width(width).Foreground(t.TextMuted())
	itemStyle := baseStyle.Width(width).
		Foreground(t.Text()).
		Background(t.Background())

	if selected {
		itemStyle = itemStyle.
			Background(t.Primary()).
			Foreground(t.Background()).
			Bold(true)
		descStyle = descStyle.
			Background(t.Primary()).
			Foreground(t.Background())
	}

	title := itemStyle.Padding(0, 1).Render(ci.Title)
	if ci.Description != "" {
		description := descStyle.Padding(0, 1).Render(ci.Description)
		return lipgloss.JoinVertical(lipgloss.Left, title, description)
	}
	return title
}

// CommandSelectedMsg is sent when a command is selected
type CommandSelectedMsg struct {
	Command Command
}

// CloseCommandDialogMsg is sent when the command dialog is closed
type CloseCommandDialogMsg struct{}

// OpenCommandDialogMsg is sent to open the command dialog with an initial filter query
// (used by slash commands, e.g. typing "/new" in the editor)
type OpenCommandDialogMsg struct {
	Query string
}

// CommandDialog interface for the command selection dialog
type CommandDialog interface {
	tea.Model
	layout.Bindings
	SetCommands(commands []Command)
	SetQuery(query string)
}

type commandDialogCmp struct {
	listView    utilComponents.SimpleList[Command]
	allCommands []Command
	query       string
	width       int
	height      int
}

type commandKeyMap struct {
	Enter  key.Binding
	Escape key.Binding
}

var commandKeys = commandKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select command"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}

func (c *commandDialogCmp) Init() tea.Cmd {
	return c.listView.Init()
}

func (c *commandDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, commandKeys.Enter):
			selectedItem, idx := c.listView.GetSelectedItem()
			if idx != -1 {
				return c, util.CmdHandler(CommandSelectedMsg{
					Command: selectedItem,
				})
			}
		case key.Matches(msg, commandKeys.Escape):
			return c, util.CmdHandler(CloseCommandDialogMsg{})
		case msg.Type == tea.KeyBackspace:
			if len(c.query) > 0 {
				c.query = c.query[:len(c.query)-1]
				c.applyFilter()
			}
			return c, nil
		default:
			// Typing filters the command list (slash-command style)
			if s := msg.String(); len(s) == 1 {
				c.query += s
				c.applyFilter()
				return c, nil
			}
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
	}

	u, cmd := c.listView.Update(msg)
	c.listView = u.(utilComponents.SimpleList[Command])
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

// applyFilter filters the commands list by the current query (e.g. "/new")
func (c *commandDialogCmp) applyFilter() {
	slug := strings.ToLower(strings.TrimPrefix(c.query, "/"))
	if slug == "" {
		c.listView.SetItems(c.allCommands)
		return
	}

	filtered := make([]Command, 0)
	for _, cmd := range c.allCommands {
		if strings.HasPrefix(strings.ToLower(cmd.ID), slug) {
			filtered = append(filtered, cmd)
		}
	}
	c.listView.SetItems(filtered)
}

func (c *commandDialogCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	maxWidth := 40

	commands := c.listView.GetItems()

	for _, cmd := range commands {
		if len(cmd.Title) > maxWidth-4 {
			maxWidth = len(cmd.Title) + 4
		}
		if cmd.Description != "" {
			if len(cmd.Description) > maxWidth-4 {
				maxWidth = len(cmd.Description) + 4
			}
		}
	}

	c.listView.SetMaxWidth(maxWidth)

	title := baseStyle.
		Foreground(t.Primary()).
		Bold(true).
		Width(maxWidth).
		Padding(0, 1).
		Render("Commands")

	// Show the current filter query (slash command style)
	queryText := c.query
	if !strings.HasPrefix(queryText, "/") {
		queryText = "/" + queryText
	}
	query := baseStyle.
		Foreground(t.TextMuted()).
		Width(maxWidth).
		Padding(0, 1).
		Render(queryText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		baseStyle.Width(maxWidth).Render(""),
		query,
		baseStyle.Width(maxWidth).Render(""),
		baseStyle.Width(maxWidth).Render(c.listView.View()),
		baseStyle.Width(maxWidth).Render(""),
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (c *commandDialogCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(commandKeys)
}

func (c *commandDialogCmp) SetCommands(commands []Command) {
	c.allCommands = commands
	c.query = ""
	c.applyFilter()
}

func (c *commandDialogCmp) SetQuery(query string) {
	c.query = query
	c.applyFilter()
}

// NewCommandDialogCmp creates a new command selection dialog
func NewCommandDialogCmp() CommandDialog {
	listView := utilComponents.NewSimpleList[Command](
		[]Command{},
		10,
		"No commands available",
		false,
	)
	return &commandDialogCmp{
		listView: listView,
	}
}
