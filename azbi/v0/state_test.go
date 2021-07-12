package v0

import (
	"github.com/epiphany-platform/e-structures/globals"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/stretchr/testify/assert"
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
