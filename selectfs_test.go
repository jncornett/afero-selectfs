package selectfs_test

import (
	"io/ioutil"
	"testing"

	selectfs "github.com/jncornett/afero-selectfs"
	"github.com/spf13/afero"
)

func TestSelectFsOpen(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		err      bool
	}{
		{
			name:     "a://file.txt",
			contents: "a",
		},
		{
			name:     "b://file.txt",
			contents: "b",
		},
		{
			name: "c://file.txt",
			err:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				a afero.Fs = afero.NewMemMapFs()
				b afero.Fs = afero.NewMemMapFs()
			)
			if err := persist(a, "file.txt", "a"); err != nil {
				t.Fatal(err)
			}
			if err := persist(b, "file.txt", "b"); err != nil {
				t.Fatal(err)
			}
			sfs := selectfs.NewFs(map[string]afero.Fs{"a": a, "b": b})
			f, err := sfs.Open(test.name)
			if test.err {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			s := string(buf)
			if test.contents != s {
				t.Errorf("expected contents to be %q, got %q", test.contents, s)
			}
		})
	}
}

func persist(fs afero.Fs, name string, contents string) error {
	f, err := fs.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	return err
}
