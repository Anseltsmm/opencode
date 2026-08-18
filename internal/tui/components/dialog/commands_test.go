package dialog

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCommandDialogFilter(t *testing.T) {
	d := NewCommandDialogCmp()
	d.SetCommands([]Command{
		{ID: "models", Title: "Switch Model"},
		{ID: "new", Title: "New Session"},
		{ID: "help", Title: "Toggle Help"},
	})
	d.SetQuery("/")

	// Simulate typing "m"
	updated, _ := d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	d = updated.(CommandDialog)

	items := d.(*commandDialogCmp).listView.GetItems()
	if len(items) != 1 || items[0].ID != "models" {
		t.Fatalf("expected 1 item (models), got %d: %+v", len(items), items)
	}

	// Backspace should restore all commands
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	d = updated.(CommandDialog)
	items = d.(*commandDialogCmp).listView.GetItems()
	if len(items) != 3 {
		t.Fatalf("expected 3 items after backspace, got %d", len(items))
	}
}
