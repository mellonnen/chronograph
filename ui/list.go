package ui

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mellonnen/chronograph/models"
)

// item wraps a listable resource to satisfy the list.Item interface.
type item struct {
	models.Listable
}

func (i item) Title() string       { return i.GetName() }
func (i item) FilterValue() string { return i.GetName() }
func (i item) Description() string { return i.GetDescription() }

// listModel represents a list that contain some listable resource.
type listModel struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap

	itemType Resource
}

// newList specifies a new list model for the provided listables and resource type.
func newList[L models.Listable](listables []L, resourceType Resource) listModel {
	delegateKeys := newDelegateKeyMap(resourceType)
	m := listModel{
		list:         list.New(itemsFromListable(listables), newDelegate(delegateKeys), 0, 0),
		keys:         newListKeyMap(resourceType),
		delegateKeys: delegateKeys,
		itemType:     resourceType,
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			m.keys.create,
			m.keys.toggleHelp,
		}
	}
	m.list.Title = strings.Title(fmt.Sprintf("%ss", resourceType))
	return m
}

// newDelegate creates a new delegate for the list.
// A delegate is the ACTUAL list element. So all logic that directly deals
// with a current member of the list should be handled by the delegate.
func newDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// Removal of resources works by having the delegate detect
			// key stroke. It then messages to the model to delete the list
			// item from the database. If the model succeeds, the same message will
			// propagate back to the delegate that then removes the list item from the UI.

			// Detect removal of items.
			case key.Matches(msg, keys.remove):
				return removeResourceCmd(m.Index())
			}

			// The message has propagated back -> we can delete the item.
		case removeResourceMsg:
			m.RemoveItem(msg.index)
			if len(m.Items()) == 0 {
				keys.remove.SetEnabled(false)
			}
		}
		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d
}

// listKeyMap specifies which keys the list should detect.
type listKeyMap struct {
	create     key.Binding
	toggleHelp key.Binding
}

// newListKeyMap returns a key map for the list.
func newListKeyMap(resourceType Resource) *listKeyMap {
	return &listKeyMap{
		create:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", fmt.Sprintf("add %s", resourceType))),
		toggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

// delegateKeyMap specifies the keys that a delegate should detect,
// that is events that are linked to a SINGLE list item.
type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// newDelegateKeyMap returns a new key map for the delegate.
func newDelegateKeyMap(resourceType Resource) *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", fmt.Sprintf("choose %s", resourceType))),
		remove: key.NewBinding(key.WithKeys("x", "backspace"), key.WithHelp("x", fmt.Sprintf("remove %s", resourceType))),
	}
}

// update updates the list.
func (m listModel) update(msg tea.Msg) (listModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.create):
			m.delegateKeys.remove.SetEnabled(true)
			workspace := models.Workspace{
				Name:        "Misc",
				Description: sql.NullString{String: "General workspace", Valid: true},
			}

			cmds = append(cmds, addWorkspaceCmd(workspace))

		case key.Matches(msg, m.keys.toggleHelp):
			m.list.SetShowHelp(!m.list.ShowHelp())
		}

	case addWorkspaceMsg:
		m.delegateKeys.remove.SetEnabled(true)
		cmds = append(cmds, m.list.InsertItem(-1, item{msg.Workspace}))
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// view returns the view for the list.
func (l listModel) view() string {
	return l.list.View()
}

// Helpers

// itemsFromListable loops over a slice of listables and converts them to a slice of []list.Item
func itemsFromListable[L models.Listable](listable []L) []list.Item {
	l := make([]list.Item, len(listable))
	for i, x := range listable {
		l[i] = item{x}
	}
	return l
}
