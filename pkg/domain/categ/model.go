package categ

import (
	"fmt"
	"strings"
)

type EntryID string

type Entry struct {
	Type                  string
	ID                    EntryID
	Name                  string
	Description           string
	C                     string
	I                     string
	A                     string
	T                     string
	PD                    string
	NSI                   string
	SESE                  string
	OTCI                  string
	CSRDirection          string // CS&R direction for control type
	SPSA                  string
	GDPR                  bool
	ExternalSupplier      bool
	AssetType             string
	OperationalCapability string
	PartOfGISR            bool
	LastUpdated           string
	OldID                 string
}

func (e Entry) Indent(s string) string {
	n := len(strings.Split(string(e.ID), ".")) - 1
	return strings.Repeat(s, n)
}

func (e Entry) String() string {
	return fmt.Sprintf("%s - %s [%s]", e.ID, e.Name, e.Type)
}
