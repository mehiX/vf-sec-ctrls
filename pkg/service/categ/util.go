package categ

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mehix/vf-sec-ctrls/pkg/domain/categ"
)

func sortEntries(entries []categ.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		partsI := strings.Split(string(entries[i].ID), ".")
		sI := ""
		for _, p := range partsI {
			sI += fmt.Sprintf("%02s", p)
		}
		partsJ := strings.Split(string(entries[j].ID), ".")
		sJ := ""
		for _, p := range partsJ {
			sJ += fmt.Sprintf("%02s", p)
		}
		return sI < sJ
	})
}

func sortIDs(ids []string) {
	sort.Slice(ids, func(i, j int) bool {
		partsI := strings.Split(ids[i], ".")
		sI := ""
		for _, p := range partsI {
			sI += fmt.Sprintf("%02s", p)
		}
		partsJ := strings.Split(ids[j], ".")
		sJ := ""
		for _, p := range partsJ {
			sJ += fmt.Sprintf("%02s", p)
		}
		return sI < sJ
	})

}
