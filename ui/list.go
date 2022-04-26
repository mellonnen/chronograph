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

type item struct {
	models.Listable
}

func (i item) Title() string       { return i.GetName() }
func (i item) FilterValue() string { return i.GetName() }
func (i item) Description() string { return i.GetDescription() }

func itemsFromListable[L models.Listable](listable []L) []list.Item {
	l := make([]list.Item, len(listable))
	for i, x := range listable {
		l[i] = item{x}
	}
	return l
}

type listModel struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap

	itemType Resource
}

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

func newDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// signal that the resource should be removed.
			case key.Matches(msg, keys.remove):
				return removeResourceCmd(m.Index())
			}

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

type listKeyMap struct {
	create     key.Binding
	toggleHelp key.Binding
}

func newListKeyMap(resourceType Resource) *listKeyMap {
	return &listKeyMap{
		create:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", fmt.Sprintf("add %s", resourceType))),
		toggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap(resourceType Resource) *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", fmt.Sprintf("choose %s", resourceType))),
		remove: key.NewBinding(key.WithKeys("x", "backspace"), key.WithHelp("x", fmt.Sprintf("remove %s", resourceType))),
	}
}

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

func (l listModel) view() string {
	return l.list.View()
}
