package v0

import (
	"encoding/json"
	"errors"
	"github.com/epiphany-platform/e-structures/utils/validators"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	maps "github.com/mitchellh/mapstructure"

	"github.com/go-playground/validator/v10"
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
	Status Status       `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *azbi.Config `json:"config" validate:"omitempty"`
	Output *azbi.Output `json:"output" validate:"omitempty"`
}

func (s *AzBIState) GetConfig() *azbi.Config {
	if s == nil {
		return nil
	}
	return s.Config
}

func (s *AzBIState) GetOutput() *azbi.Output {
	if s == nil {
		return nil
	}
	return s.Output
}

type AzKSState struct {
	Status Status       `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *azks.Config `json:"config" validate:"omitempty"`
	Output *azks.Output `json:"output" validate:"omitempty"`
}

func (s *AzKSState) GetConfig() *azks.Config {
	if s == nil {
		return nil
	}
	return s.Config
}

func (s *AzKSState) GetOutput() *azks.Output {
	if s == nil {
		return nil
	}
	return s.Output
}

type State struct {
	Kind    *string    `json:"kind" validate:"required,eq=state"`
	Version *string    `json:"version" validate:"required,major=~0"`
	Unused  []string   `json:"-"`
	AzBI    *AzBIState `json:"azbi" validate:"omitempty"`
	AzKS    *AzKSState `json:"azks" validate:"omitempty"`
}

func (s *State) GetAzBIState() *AzBIState {
	if s == nil {
		return nil
	}
	return s.AzBI
}

func (s *State) GetAzKSState() *AzKSState {
	if s == nil {
		return nil
	}
	return s.AzKS
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

//TODO implement more interesting validation
func (s *State) isValid() error {
	if s == nil {
		return errors.New("state is nil")
	}
	validate := validator.New()
	err := validate.RegisterValidation("major", validators.HasMajorVersionLike)
	if err != nil {
		return err
	}
	err = validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}
