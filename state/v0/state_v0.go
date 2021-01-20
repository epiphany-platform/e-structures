package v0

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver"
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	maps "github.com/mitchellh/mapstructure"
)

type Status string

const (
	kind    = "state"
	version = "v0.0.3"

	Initialized Status = "initialized"
	Applied     Status = "applied"
	Destroyed   Status = "destroyed"
)

type AzBIState struct {
	Status Status       `json:"status"`
	Config *azbi.Config `json:"config"`
	Output *azbi.Output `json:"output"`
}

type AzKSState struct {
	Status Status       `json:"status"`
	Config *azks.Config `json:"config"`
	Output *azks.Output `json:"output"`
}

type State struct {
	Kind    *string    `json:"kind"`
	Version *string    `json:"version"`
	Unused  []string   `json:"-"`
	AzBI    *AzBIState `json:"azbi"`
	AzKS    *AzKSState `json:"azks"`
}

//TODO test
func NewState() *State {
	return &State{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Unused:  []string{},
		AzBI:    &AzBIState{},
		AzKS:    &AzKSState{},
	}
}

func (s *State) Marshal() (b []byte, err error) {
	//TODO validate that all required fields are filled
	return json.MarshalIndent(s, "", "\t")
}

func (s *State) Unmarshal(b []byte) (err error) {
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

var (
	KindMissingValidationError    = errors.New("field 'Kind' cannot be nil")
	VersionMissingValidationError = errors.New("field 'Version' cannot be nil")
	MajorVersionMismatchError     = errors.New("version of loaded structure has MAJOR part different than required")
)

//TODO implement more interesting validation
func (s *State) isValid() error {
	if s.Version == nil {
		return VersionMissingValidationError
	}
	if s.Kind == nil {
		return KindMissingValidationError
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return err
	}
	constraint, err := semver.NewConstraint(fmt.Sprintf("~%d", v.Major()))
	if err != nil {
		return err
	}
	vv, err := semver.NewVersion(*s.Version)
	if err != nil {
		return err
	}
	if !constraint.Check(vv) {
		return MajorVersionMismatchError
	}
	return nil
}
