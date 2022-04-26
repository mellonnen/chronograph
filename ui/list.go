package ui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

type listKeyMap struct {
	create     key.Binding
	choose     key.Binding
	remove     key.Binding
	toggleHelp key.Binding
}

func newListKeyMap(itemType string) listKeyMap {
	return listKeyMap{
		create:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", fmt.Sprintf("add %s", itemType))),
		choose:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", fmt.Sprintf("choose %s", itemType))),
		remove:     key.NewBinding(key.WithKeys("r"), key.WithHelp("a", fmt.Sprintf("remove %s", itemType))),
		toggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func newList(in interface{}) (list.Model, listKeyMap, error) {
	var listKeys listKeyMap
	var l list.Model
	switch v := in.(type) {
	case []models.Workspace:
		listKeys = newListKeyMap("workspace")
		l = list.New(itemsFromListable(v), list.NewDefaultDelegate(), 0, 0)
		l.Title = "Workspaces"
	case []models.Repo:
		listKeys = newListKeyMap("repo")
		l = list.New(itemsFromListable(v), list.NewDefaultDelegate(), 0, 0)
		l.Title = "Repos"
	case []models.Task:
		listKeys = newListKeyMap("task")
		l = list.New(itemsFromListable(v), list.NewDefaultDelegate(), 0, 0)
		l.Title = "Tasks"
	default:
		return list.Model{}, listKeyMap{}, errors.New("non-listable data")
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.create,
			listKeys.choose,
			listKeys.remove,
			listKeys.toggleHelp,
		}
	}
	return l, listKeys, nil
}
