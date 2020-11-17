package state

import (
	"encoding/json"
	"errors"
	maps "github.com/mitchellh/mapstructure"
	"github.com/mkyc/go-stucts-versioning-tests/pkg/azbi"
	"github.com/mkyc/go-stucts-versioning-tests/pkg/to"
)

type Status string

const (
	kind    = "state"
	version = "0.0.1"

	Initialized Status = "initialized"
	Applied     Status = "applied"
	Destroyed   Status = "destroyed"
)

type AzBIState struct {
	Status Status       `json:"status"`
	Config *azbi.Config `json:"config"`
	Output *azbi.Output `json:"output"`
}

type State struct {
	Kind    *string    `json:"kind"`
	Version *string    `json:"version"`
	Unused  []string   `json:"-"`
	AzBI    *AzBIState `json:"azbi"`
}

func NewState() *State {
	return &State{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Unused:  []string{},
		AzBI:    &AzBIState{},
	}
}

func (s *State) Save() (b []byte, err error) {
	return json.MarshalIndent(s, "", "\t")
}

func (s *State) Load(b []byte) (err error) {
	var input map[string]interface{}
	if err = json.Unmarshal(b, &input); err != nil {
		return
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &s,
	})
	if err != nil {
		return
	}
	err = d.Decode(input)
	if err != nil {
		return
	}
	s.Unused = md.Unused
	err = s.isValid()
	return
}

//TODO implement more interesting validation
func (s *State) isValid() error {
	if s.Version == nil {
		return errors.New("field 'Version' cannot be nil")
	}
	return nil
}
