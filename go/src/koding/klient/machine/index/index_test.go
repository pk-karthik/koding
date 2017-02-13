package index_test

import (
	"encoding/json"
	"sort"
	"testing"

	"koding/klient/machine/index"
	"koding/klient/machine/index/indextest"
)

// filetree defines a simple directory structure that will be created for test
// purposes. The values of this map stores file sizes.
var filetree = map[string]int64{
	"a.txt":        128,
	"b.bin":        300 * 1024,
	"c/":           0,
	"c/ca.txt":     2 * 1024,
	"c/cb.bin":     1024 * 1024,
	"d/":           0,
	"d/da.txt":     5 * 1024,
	"d/db.txt":     256,
	"d/dc/":        0,
	"d/dc/dca.txt": 3 * 1024,
	"d/dc/dcb.txt": 1024,
}

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		Op      func(string) error
		Changes index.ChangeSlice
	}{
		"add file": {
			Op: indextest.WriteFile("d/test.bin", 40*1024),
			Changes: index.ChangeSlice{
				index.NewChange("d", index.ChangeMetaUpdate),
				index.NewChange("d/test.bin", index.ChangeMetaAdd),
			},
		},
		"add dir": {
			Op: indextest.AddDir("e"),
			Changes: index.ChangeSlice{
				index.NewChange("e", index.ChangeMetaAdd),
			},
		},
		"remove file": {
			Op: indextest.RmAllFile("c/cb.bin"),
			Changes: index.ChangeSlice{
				index.NewChange("c", index.ChangeMetaUpdate),
				index.NewChange("c/cb.bin", index.ChangeMetaRemove),
			},
		},
		"remove dir": {
			Op: indextest.RmAllFile("c"),
			Changes: index.ChangeSlice{
				index.NewChange("c", index.ChangeMetaRemove),
				index.NewChange("c/ca.txt", index.ChangeMetaRemove),
				index.NewChange("c/cb.bin", index.ChangeMetaRemove),
			},
		},
		"rename file": {
			Op: indextest.MvFile("b.bin", "c/cc.bin"),
			Changes: index.ChangeSlice{
				index.NewChange("b.bin", index.ChangeMetaRemove),
				index.NewChange("c", index.ChangeMetaUpdate),
				index.NewChange("c/cc.bin", index.ChangeMetaAdd),
			},
		},
		"write file": {
			Op: indextest.WriteFile("b.bin", 1024),
			Changes: index.ChangeSlice{
				index.NewChange("b.bin", index.ChangeMetaUpdate),
			},
		},
		"chmod file": {
			Op: indextest.ChmodFile("d/dc/dca.txt", 0600),
			Changes: index.ChangeSlice{
				index.NewChange("d/dc/dca.txt", index.ChangeMetaUpdate),
			},
		},
	}

	for name, test := range tests {
		// capture range variable here
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, clean, err := indextest.GenerateTree(filetree)
			if err != nil {
				t.Fatalf("want err = nil; got %v", err)
			}
			defer clean()

			idx, err := index.NewIndexFiles(root)
			if err != nil {
				t.Fatalf("want err = nil; got %v", err)
			}

			if err := test.Op(root); err != nil {
				t.Fatalf("want err = nil; got %v", err)
			}

			// Synchronize underlying file-system.
			indextest.Sync()

			cs := idx.Compare(root)
			sort.Sort(cs)
			if len(cs) != len(test.Changes) {
				t.Fatalf("want index.Changes count = %d; got %d", len(test.Changes), len(cs))
			}

			// Copy time from result to tests.
			for i, tc := range test.Changes {
				if cs[i].Path() != tc.Path() {
					t.Errorf("want index.Change path = %q; got %q", tc.Path(), cs[i].Path())
				}
				if cs[i].Meta() != tc.Meta() {
					t.Errorf("want index.Change meta = %bb; got %bb", tc.Meta(), cs[i].Meta())
				}
			}

			idx.Apply(root, cs)
			if cs = idx.Compare(root); len(cs) != 0 {
				t.Errorf("want no index.Changes after apply; got %#v", cs)
			}
		})
	}
}

func TestIndexCount(t *testing.T) {
	tests := map[string]struct {
		MaxSize  int64
		Expected int
	}{
		"all items": {
			MaxSize:  -1,
			Expected: 11,
		},
		"less than 100kiB": {
			MaxSize:  100 * 1024,
			Expected: 9,
		},
		"zero": {
			MaxSize:  0,
			Expected: 0,
		},
	}

	for name, test := range tests {
		// capture range variable here
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, clean, err := indextest.GenerateTree(filetree)
			if err != nil {
				t.Fatalf("want err = nil; got %v", err)
			}
			defer clean()

			idx, err := index.NewIndexFiles(root)
			if err != nil {
				t.Fatalf("want err = nil; got %v", err)
			}

			if count := idx.Count(test.MaxSize); count != test.Expected {
				t.Errorf("want count = %d; got %d", test.Expected, count)
			}
		})
	}
}

func TestIndexJSON(t *testing.T) {
	root, clean, err := indextest.GenerateTree(filetree)
	if err != nil {
		t.Fatalf("want err = nil; got %v", err)
	}
	defer clean()

	idx, err := index.NewIndexFiles(root)
	if err != nil {
		t.Fatalf("want err = nil; got %v", err)
	}

	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("want err = nil; got %v", err)
	}

	idx = index.NewIndex()
	if err := json.Unmarshal(data, idx); err != nil {
		t.Fatalf("want err = nil; got %v", err)
	}

	if cs := idx.Compare(root); len(cs) != 0 {
		t.Errorf("want no changes after apply; got %#v", cs)
	}
}
