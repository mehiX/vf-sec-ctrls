package categ

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mehix/vf-sec-ctrls/tree"
	"github.com/xuri/excelize/v2"
)

type Hierarchy struct {
	data map[string]Entry
	Root *tree.Node
}

func NewHierarchy() *Hierarchy {
	return &Hierarchy{data: make(map[string]Entry, 0)}
}

type EntryID string

type Entry struct {
	Type string
	ID   EntryID
	Name string
}

func (e Entry) Indent(s string) string {
	n := len(strings.Split(string(e.ID), ".")) - 1
	return strings.Repeat(s, n)
}

func (e Entry) String() string {
	return fmt.Sprintf("%s - %s [%s]", e.ID, e.Name, e.Type)
}

func NewFromFile(fn string, sheetName string) (*Hierarchy, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h, err := NewFromReader(f, sheetName)

	return h, err
}

func NewFromReader(in io.ReadCloser, sheetName string) (*Hierarchy, error) {
	f, err := excelize.OpenReader(in)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if sheetName == "" {
		lst := f.GetSheetList()
		if len(lst) == 0 {
			return nil, fmt.Errorf("no worksheets available in reader")
		}
		sheetName = lst[0]
	}
	rows, err := f.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	h := NewHierarchy()
	for counter := 0; rows.Next(); counter++ {
		// first row is header
		if counter == 0 {
			continue
		}

		cols, err := rows.Columns()
		if err != nil {
			log.Printf("Row %d: %v\n", counter+1, err)
		}

		entry := Entry{
			Type: strings.TrimSpace(cols[0]),
			ID:   EntryID(strings.TrimSpace(cols[1])),
			Name: strings.TrimSpace(cols[2]),
		}

		h.data[string(entry.ID)] = entry
	}

	h.Root = BuildTree(h.Entries())

	return h, nil
}

func (h *Hierarchy) Entries() []Entry {
	entries := make([]Entry, 0, len(h.data))
	for _, v := range h.data {
		entries = append(entries, v)
	}
	sortEntries(entries)
	return entries
}

func BuildTree(entries []Entry) *tree.Node {

	ids := make([]string, len(entries))
	for idx := range entries {
		ids[idx] = string(entries[idx].ID)
	}

	return tree.New(ids)
}

func (h *Hierarchy) Find(id string) *tree.Node {
	return tree.FindNode(h.Root, id)
}

func (h *Hierarchy) Print(pretty bool) {
	data := h.Entries()
	sort.Slice(data, func(i, j int) bool {
		return data[i].ID < data[j].ID
	})

	for _, d := range data {
		if pretty {
			fmt.Printf("%s%s\n", d.Indent("  "), d)
		} else {
			fmt.Println(d)
		}
	}
}

func (h *Hierarchy) ControlIDs() []string {
	all := tree.EdgesFrom(h.Root)

	sortIDs(all)

	return all
}

func (h *Hierarchy) ControlsByCategory(categoryID string) []Entry {
	n := tree.FindNode(h.Root, categoryID)
	if n == nil {
		return []Entry{}
	}

	controlIDs := tree.EdgesFrom(n)
	var controls []Entry
	for _, c := range controlIDs {
		controls = append(controls, h.data[c])
	}

	sortEntries(controls)

	return controls

}
