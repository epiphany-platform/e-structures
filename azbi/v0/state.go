package v0

import (
	"encoding/json"
	"errors"
	"github.com/epiphany-platform/e-structures/globals"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
	maps "github.com/mitchellh/mapstructure"
	"io/ioutil"
	"os"
)

type State struct {
	Meta   *Meta          `json:"meta" validate:"required"`
	Status globals.Status `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *Config        `json:"config" validate:"omitempty"`
	Output *Output        `json:"output" validate:"omitempty"`
	Unused []string       `json:"-"`
}

func (s *State) Init(moduleVersion string) {
	*s = State{
		Meta: &Meta{
			Kind:          to.StrPtr(stateKind),
			Version:       to.StrPtr(stateVersion),
			ModuleVersion: to.StrPtr(moduleVersion),
		},
		Status: globals.Initialized,
		Config: nil, // TODO should it be nil?
		Output: nil, // TODO should it be nil?
		Unused: []string{},
	}
}

func (s *State) Backup(path string) error {
	return globals.Backup(s, path)
}

func (s *State) Load(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	// TODO backup raw here

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var input map[string]interface{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return err
	}

	// TODO check if current version here

	state := &State{}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{Metadata: &md, TagName: "json", Result: &state})
	if err != nil {
		return err
	}

	err = d.Decode(input)
	if err != nil {
		return err
	}

	state.Unused = md.Unused

	// TODO rethink if validation should be done here

	err = state.Valid()
	if err != nil {
		return err
	}

	*s = *state
	return nil
}

func (s *State) Save(path string) error {
	return globals.Save(s, path)
}

func (s *State) Print() ([]byte, error) {
	return globals.Print(s)
}

func (s *State) Valid() error {
	if s == nil {
		return errors.New("expected state is nil")
	}
	validate := validator.New()
	err := validate.RegisterValidation("version", validators.HasVersion)
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

func (s *State) Upgrade(_ string) error {
	return nil
}

type Output struct {
	RgName   *string         `json:"rg_name"`
	VnetName *string         `json:"vnet_name"`
	VmGroups []OutputVmGroup `json:"vm_groups"`
}

type OutputVmGroup struct {
	Name *string    `json:"vm_group_name"`
	Vms  []OutputVm `json:"vms"`
}

type OutputVm struct {
	Name       *string          `json:"vm_name"`
	PrivateIps []string         `json:"private_ips"`
	PublicIp   *string          `json:"public_ip"`
	DataDisks  []OutputDataDisk `json:"data_disks"`
}

type OutputDataDisk struct {
	Size *int `json:"size"`
	Lun  *int `json:"lun"`
}
