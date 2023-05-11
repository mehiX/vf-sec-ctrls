package tree

import "testing"

func TestNewTree(t *testing.T) {
	data := []string{"1", "1.1", "1.2", "1.4", "1.2.1", "1.2.2", "3", "3.1", "3.3", "3.1.1"}

	tdata := New(data)

	if tdata.Value != "" || len(tdata.Children) != 2 {
		t.Fatalf("wrong root. Value: %s, Children: %d", tdata.Value, len(tdata.Children))
	}
}
func TestFindNode(t *testing.T) {

	data := []string{"1", "1.1", "1.2", "1.4", "1.2.1", "1.2.2", "3", "3.1", "3.3", "3.1.1"}

	tdata := New(data)

	f := FindNode(tdata, "1.2")
	if f.Value != "1.2" || len(f.Children) != 2 {
		t.Fatalf("Found wrong node: %v", f)
	}

	f = FindNode(tdata, "1.2.1")
	if f.Value != "1.2.1" || len(f.Children) != 0 {
		t.Fatalf("Found wrong node: %v", f)
	}

	f = FindNode(tdata, "5")
	if f != nil {
		t.Fatalf("should return nil for missing value")
	}
}

func TestAllLeaves(t *testing.T) {

	data := []string{"1", "1.1", "1.2", "1.4", "1.2.1", "1.2.2", "3", "3.1", "3.3", "3.1.1"}

	tdata := New(data)

	all := EdgesFrom(tdata)

	if len(all) != 6 {
		t.Fatalf("Wrong number of leaves: %d", len(all))
	}
}
