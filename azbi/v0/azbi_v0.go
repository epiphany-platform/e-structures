package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	maps "github.com/mitchellh/mapstructure"
	"github.com/mkyc/go-stucts-versioning-tests/utils/to"
)

const (
	kind    = "azbi"
	version = "v0.0.1-alfa"
)

type Params struct {
	Name             *string  `json:"name"`
	VmsCount         *int     `json:"vms_count"`
	UsePublicIP      *bool    `json:"use_public_ip"`
	Location         *string  `json:"location"`
	AddressSpace     []string `json:"address_space"`
	AddressPrefixes  []string `json:"address_prefixes"`
	RsaPublicKeyPath *string  `json:"rsa_pub_path"`
}

type Config struct {
	Kind    *string  `json:"kind"`
	Version *string  `json:"version"`
	Params  *Params  `json:"params"`
	Unused  []string `json:"-"`
}

func NewConfig() *Config {
	return &Config{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Params: &Params{
			Name:             to.StrPtr("epiphany"),
			VmsCount:         to.IntPtr(3),
			UsePublicIP:      to.BooPtr(true),
			Location:         to.StrPtr("northeurope"),
			AddressSpace:     []string{"10.0.0.0/16"},
			AddressPrefixes:  []string{"10.0.1.0/24"},
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
		},
		Unused: []string{},
	}
}

func (c *Config) Save() (b []byte, err error) {
	//TODO validate that all required fields are filled
	return json.MarshalIndent(c, "", "\t")
}

func (c *Config) Load(b []byte) (err error) {
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

var (
	KindMissingValidationError    = errors.New("field 'Kind' cannot be nil")
	VersionMissingValidationError = errors.New("field 'Version' cannot be nil")
	ParamsMissingValidationError  = errors.New("params section missing")
	MinimalParamsValidationError  = errors.New("at least 'vms_count', 'location' and 'name' parameters are required")
	MajorVersionMismatchError     = errors.New("version of loaded structure has MAJOR part different than required")
)

//TODO implement more interesting validation
func (c *Config) isValid() error {
	if c.Kind == nil {
		return KindMissingValidationError
	}
	if c.Version == nil {
		return VersionMissingValidationError
	}
	if c.Params == nil {
		return ParamsMissingValidationError
	}
	if c.Params.Name == nil || c.Params.VmsCount == nil || c.Params.Location == nil {
		return MinimalParamsValidationError
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return err
	}
	constraint, err := semver.NewConstraint(fmt.Sprintf("~%d", v.Major()))
	if err != nil {
		return err
	}
	vv, err := semver.NewVersion(*c.Version)
	if err != nil {
		return err
	}
	if !constraint.Check(vv) {
		return MajorVersionMismatchError
	}
	return nil
}

type Output struct {
	PrivateIps []string `json:"private_ips"`
	PublicIps  []string `json:"public_ips"`
	RgName     *string  `json:"rg_name"`
	VmNames    []string `json:"vm_names"`
	VnetName   *string  `json:"vnet_name"`
}
