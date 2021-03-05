package v0

import (
	"encoding/json"
	"errors"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
	maps "github.com/mitchellh/mapstructure"
)

const (
	kind    = "hi"
	version = "v0.0.1"
)

type MountPoint struct {
	Lun  *int    `json:"lun" validate:"required,min=0"`
	Path *string `json:"path" validate:"required,min=1"`
}

type Host struct {
	Name *string `json:"name" validate:"required,min=1"`
	Ip   *string `json:"ip" validate:"required,min=1"`
}

type VmGroup struct {
	Name        *string      `json:"name" validate:"required,min=1"`
	AdminUser   *string      `json:"admin_user" validate:"required,min=1"`
	Hosts       []Host       `json:"hosts" validate:"required,min=1,dive"`
	MountPoints []MountPoint `json:"mount_point" validate:"omitempty,dive"`
}

type Params struct {
	VmGroups          []VmGroup `json:"vm_groups" validate:"required,dive"`
	RsaPrivateKeyPath *string   `json:"rsa_private_path" validate:"required,min=1"`
}

type Config struct {
	Kind    *string  `json:"kind" validate:"required,eq=hi"`
	Version *string  `json:"version" validate:"required,version=~0"`
	Params  *Params  `json:"params" validate:"required"`
	Unused  []string `json:"-"`
}

func (c *Config) GetParams() *Params {
	if c == nil {
		return nil
	}
	return c.Params
}

//TODO test
func NewConfig() *Config {
	return &Config{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Params: &Params{
			VmGroups: []VmGroup{
				{
					Name:      to.StrPtr("vm-group0"),
					AdminUser: to.StrPtr("operations"),
					Hosts: []Host{
						{
							Name: to.StrPtr("epiphany-vm-group0-1"),
							Ip:   to.StrPtr("10.0.1.4"),
						},
					},
					MountPoints: []MountPoint{
						{
							Lun:  to.IntPtr(10),
							Path: to.StrPtr("/data/test"),
						},
					},
				},
			},
		},
		Unused: []string{},
	}
}

func (c *Config) Marshal() (b []byte, err error) {
	//TODO validate that all required fields are filled
	return json.MarshalIndent(c, "", "\t")
}

func (c *Config) Unmarshal(b []byte) (err error) {
	var input map[string]interface{}
	if err = json.Unmarshal(b, &input); err != nil {
		return
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &c,
	})
	if err != nil {
		return
	}
	err = d.Decode(input)
	if err != nil {
		return
	}
	c.Unused = md.Unused
	err = c.isValid()
	return
}

func (c *Config) isValid() error {
	if c == nil {
		return errors.New("hi config is nil")
	}
	validate := validator.New()

	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}
	err = validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}
