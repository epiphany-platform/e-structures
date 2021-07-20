package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/globals"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestState_Init(t *testing.T) {
	tests := []struct {
		name          string
		moduleVersion string
		want          *State
	}{
		{
			name:          "happy path",
			moduleVersion: "v1.1.1",
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: globals.Initialized,
				Unused: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := &State{}
			got.Init(tt.moduleVersion)
			a.Equal(tt.want, got)
		})
	}
}

func TestState_Backup(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr error
	}{
		{
			name:    "happy path",
			state:   &State{},
			wantErr: nil,
		},
		{
			name:    "file already exists",
			state:   &State{},
			wantErr: os.ErrExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p, err := createTempDirectory("azbi-state-backup")
			if errors.Is(tt.wantErr, os.ErrExist) {
				err = ioutil.WriteFile(filepath.Join(p, "backup-file.json"), []byte("content"), 0644)
				t.Logf("path: %s", filepath.Join(p, "backup-file.json"))
				a.NoError(err)
			}
			err = tt.state.Backup(filepath.Join(p, "backup-file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				a.Equal(tt.wantErr, err)
			} else {
				a.NoError(err)
			}
		})
	}
}
