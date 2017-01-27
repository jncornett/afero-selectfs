package selectfs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
)

const DefaultSplit = "://"

type Fs struct {
	M     map[string]afero.Fs
	Split string
}

func NewFs(m map[string]afero.Fs) *Fs {
	if m == nil {
		m = make(map[string]afero.Fs)
	}
	return &Fs{
		M:     m,
		Split: DefaultSplit,
	}
}

func (s Fs) Create(name string) (afero.File, error) {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return nil, errInvalidSelector(sel)
	} else {
		return fs.Create(p)
	}
}

func (s Fs) Mkdir(name string, perm os.FileMode) error {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Mkdir(p, perm)
	}
}

func (s Fs) MkdirAll(path string, perm os.FileMode) error {
	sel, p := parseName(path)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.MkdirAll(p, perm)
	}
}

func (s Fs) Open(name string) (afero.File, error) {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return nil, errInvalidSelector(sel)
	} else {
		return fs.Open(p)
	}
}

func (s Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return nil, errInvalidSelector(sel)
	} else {
		return fs.OpenFile(p, flag, perm)
	}
}

func (s Fs) Remove(name string) error {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Remove(p)
	}
}

func (s Fs) RemoveAll(path string) error {
	sel, p := parseName(path)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Remove(p)
	}
}

// NOTE the protocol for newname is ignored
func (s Fs) Rename(oldname string, newname string) error {
	sel, oldP := parseName(oldname)
	newSel, newP := parseName(newname)
	if sel != newSel {
		return errCrossFs
	}
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Rename(oldP, newP)
	}
}

func (s Fs) Stat(name string) (os.FileInfo, error) {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return nil, errInvalidSelector(sel)
	} else {
		return fs.Stat(p)
	}
}

func (s Fs) Name() string {
	return "SelectFs"
}

func (s Fs) Chmod(name string, mode os.FileMode) error {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Chmod(p, mode)
	}
}

func (s Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	sel, p := parseName(name)
	if fs, ok := s.M[sel]; !ok {
		return errInvalidSelector(sel)
	} else {
		return fs.Chtimes(p, atime, mtime)
	}
}

func parseName(name string) (selector string, path string) {
	split := strings.SplitN(name, "://", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", split[0]
}

func errInvalidSelector(sel string) error {
	return fmt.Errorf("Invalid selector %q", sel)
}

var errCrossFs = errors.New("Cross-filesystem operation not supported")
