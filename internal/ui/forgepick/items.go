// Package forgepick renders the registry forges as single-choice rows for the
// screen that chooses a repo's forge.
package forgepick

import (
	"charm.land/bubbles/v2/list"
	"github.com/bkenks/gitkeeper/internal/config"
	"github.com/bkenks/gitkeeper/internal/styles"
)

// item is a registry forge as shown in a selection list.
type item struct {
	name, host string
	selected   bool
}

func (i item) Title() string {
	box := styles.GlyphCheckOff
	if i.selected {
		box = styles.GlyphCheckOn
	}
	return box + " " + i.name
}

func (i item) Description() string { return i.host }

func (i item) FilterValue() string { return i.name }

// Items builds one row per forge, marking the one named selected.
func Items(forges []config.Forge, selected string) []list.Item {
	items := make([]list.Item, len(forges))
	for i, f := range forges {
		items[i] = item{name: f.Name, host: f.Host, selected: f.Name == selected}
	}
	return items
}

// IndexOf returns the position of the named forge, or -1.
func IndexOf(forges []config.Forge, name string) int {
	for i, f := range forges {
		if f.Name == name {
			return i
		}
	}
	return -1
}
