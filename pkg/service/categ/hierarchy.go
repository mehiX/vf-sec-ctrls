package categ

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mehix/vf-sec-ctrls/pkg/domain/categ"
	"github.com/mehix/vf-sec-ctrls/tree"
	"github.com/xuri/excelize/v2"
)

type Hierarchy struct {
	data map[string]categ.Entry
	Root *tree.Node
}

func NewHierarchy() *Hierarchy {
	return &Hierarchy{data: make(map[string]categ.Entry, 0)}
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

		origCols, err := rows.Columns()
		if err != nil {
			log.Printf("Row %d: %v\n", counter+1, err)
			continue
		}

		// some rows may have less columns filled with data
		cols := make([]string, 21)
		copy(cols, origCols)

		entry := categ.Entry{
			Type: strings.TrimSpace(cols[0]),
			ID:   categ.EntryID(strings.TrimSpace(cols[1])),
			Name: strings.TrimSpace(cols[2]),
		}

		entry.Description = strings.TrimSpace(cols[3])
		entry.C = strings.TrimSpace(cols[4])
		entry.I = strings.TrimSpace(cols[5])
		entry.A = strings.TrimSpace(cols[6])
		entry.T = strings.TrimSpace(cols[7])
		entry.PD = strings.TrimSpace(cols[8])
		entry.NSI = strings.TrimSpace(cols[9])
		entry.SESE = strings.TrimSpace(cols[10])
		entry.OTCI = strings.TrimSpace(cols[11])
		entry.CSRDirection = strings.TrimSpace(cols[12])
		entry.SPSA = strings.TrimSpace(cols[13])
		entry.GDPR = strings.TrimSpace(cols[14]) == "GDPR"
		entry.ExternalSupplier = strings.TrimSpace(cols[15]) != ""
		entry.AssetType = strings.TrimSpace(cols[16])
		entry.OperationalCapability = strings.TrimSpace(cols[17])
		entry.PartOfGISR = strings.ToLower(strings.TrimSpace(cols[18])) == "yes"
		entry.LastUpdated = strings.TrimSpace(cols[19])
		entry.OldID = strings.TrimSpace(cols[20])

		h.data[string(entry.ID)] = entry
	}

	h.Root = BuildTree(h.Entries())

	return h, nil
}

func NewFromList(entries []categ.Entry) *Hierarchy {
	h := &Hierarchy{
		data: make(map[string]categ.Entry),
	}
	for _, e := range entries {
		h.data[string(e.ID)] = e
	}
	h.Root = BuildTree(entries)

	return h
}

func (h *Hierarchy) Entries() []categ.Entry {
	entries := make([]categ.Entry, 0, len(h.data))
	for _, v := range h.data {
		entries = append(entries, v)
	}
	sortEntries(entries)
	return entries
}

func (h *Hierarchy) Entry(id string) *categ.Entry {
	e, ok := h.data[id]
	if !ok {
		return nil
	}

	return &e
}

func BuildTree(entries []categ.Entry) *tree.Node {

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

// ControlsByCategory returns the direct children of this category.
// It does not return the searched category.
// So if the category turns out to be a control, then an empty list will be returned
func (h *Hierarchy) ControlsByCategory(categoryID string) []categ.Entry {
	n := tree.FindNode(h.Root, categoryID)
	if n == nil {
		return []categ.Entry{}
	}

	controlIDs := tree.EdgesFrom(n)
	var controls []categ.Entry
	for _, c := range controlIDs {
		if c != categoryID {
			controls = append(controls, h.data[c])
		}
	}

	sortEntries(controls)

	return controls

}

func (h *Hierarchy) Parents(id string) []categ.Entry {

	parts := strings.Split(id, ".")
	if len(parts) == 1 {
		return []categ.Entry{}
	}

	parents := make([]categ.Entry, 0)
	for i := 0; i < len(parts)-1; i++ {
		parents = append(parents, h.data[strings.Join(parts[:i+1], ".")])
	}

	return parents
}
