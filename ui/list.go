package ui

import (
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

type listKeyMap struct {
	create     key.Binding
	choose     key.Binding
	remove     key.Binding
	toggleHelp key.Binding
}

type listModel struct {
	list list.Model
	Keys *listKeyMap

	itemType Resource
}

func newListKeyMap(resourceType Resource) *listKeyMap {
	return &listKeyMap{
		create:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", fmt.Sprintf("add %s", resourceType))),
		choose:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", fmt.Sprintf("choose %s", resourceType))),
		remove:     key.NewBinding(key.WithKeys("r"), key.WithHelp("a", fmt.Sprintf("remove %s", resourceType))),
		toggleHelp: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	}
}

func newList[L models.Listable](listables []L, resourceType Resource) listModel {
	m := listModel{
		list:     list.New(itemsFromListable(listables), list.NewDefaultDelegate(), 0, 0),
		Keys:     newListKeyMap(resourceType),
		itemType: resourceType,
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			m.Keys.create,
			m.Keys.choose,
			m.Keys.remove,
			m.Keys.toggleHelp,
		}
	}
	m.list.Title = strings.Title(fmt.Sprintf("%ss", resourceType))
	return m
}

func (l listModel) update(msg tea.Msg) tea.Cmd {
	newList, cmd := l.list.Update(msg)
	l.list = newList
	return cmd
}

func (l listModel) index() int {
	return l.list.Index()
}

func (l listModel) removeItem(index int) {
	l.list.RemoveItem(index)
}

func (l listModel) addItem(listableItem models.Listable) tea.Cmd {
	return l.list.InsertItem(-1, item{listableItem})
}

func (l listModel) filterState() list.FilterState {
	return l.list.FilterState()
}

func (l listModel) view() string {
	return l.list.View()
}
