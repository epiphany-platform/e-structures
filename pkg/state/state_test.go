package state

import (
	"github.com/mkyc/go-stucts-versioning-tests/pkg/to"
	"reflect"
	"testing"
)

func TestState_Load(t *testing.T) {
	tests := []struct {
		name      string
		args      []byte
		wantState State
		wantErr   bool
	}{
		{
			name: "empty state",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.1"
}`),
			wantState: State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.1"),
				Unused:  []string{},
				AzBI:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState := State{}
			if err := gotState.Load(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotState, tt.wantState) {
				t.Errorf("Load() gotState = \n\n%#v\n\n, want = \n\n%#v\n\n", gotState, tt.wantState)
			}
		})
	}
}
