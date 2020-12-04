package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/Masterminds/semver"
	"github.com/epiphany-platform/e-structures/utils/to"
	maps "github.com/mitchellh/mapstructure"
)

const (
	kind    = "azbi"
	version = "v0.0.2"
)

type Subnet struct {
	Name            *string  `json:"name"`
	AddressPrefixes []string `json:"address_prefixes"`
}

type Params struct {
	Name             *string  `json:"name"`
	VmsCount         *int     `json:"vms_count"`
	UsePublicIP      *bool    `json:"use_public_ip"`
	Location         *string  `json:"location"`
	AddressSpace     []string `json:"address_space"`
	Subnets          []Subnet `json:"subnets"`
	RsaPublicKeyPath *string  `json:"rsa_pub_path"`
}

type Config struct {
	Kind    *string  `json:"kind"`
	Version *string  `json:"version"`
	Params  *Params  `json:"params"`
	Unused  []string `json:"-"`
}

//TODO test
func NewConfig() *Config {
	return &Config{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Params: &Params{
			Name:         to.StrPtr("epiphany"),
			VmsCount:     to.IntPtr(3),
			UsePublicIP:  to.BooPtr(true),
			Location:     to.StrPtr("northeurope"),
			AddressSpace: []string{"10.0.0.0/16"},
			Subnets: []Subnet{
				{
					Name:            to.StrPtr("main"),
					AddressPrefixes: []string{"10.0.1.0/24"},
				},
			},
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
		},
		Unused: []string{},
	}
}

func (c *Config) Marshall() (b []byte, err error) {
	//TODO validate that all required fields are filled
	return json.MarshalIndent(c, "", "\t")
}

func (c *Config) Unmarshall(b []byte) (err error) {
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
	MajorVersionMismatchError     = errors.New("version of loaded structure has MAJOR part different than required")
)

type MinimalParamsValidationError struct {
	msg string
}

func (e MinimalParamsValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.msg)
}

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
	if c.Params.Name == nil {
		return &MinimalParamsValidationError{"'name' parameter missing"}
	}
	if c.Params.VmsCount == nil {
		return &MinimalParamsValidationError{"'vms_count' parameter missing"}
	}
	if c.Params.Location == nil {
		return &MinimalParamsValidationError{"'location' parameter missing"}
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
	if c.Params != nil && !reflect.DeepEqual(c.Params, &Params{}) {
		if c.Params.Subnets == nil || len(c.Params.Subnets) < 1 {
			return &MinimalParamsValidationError{"'subnets' list parameter missing or is 0 length"}
		}
		for _, s := range c.Params.Subnets {
			if s.Name == nil || len(*s.Name) < 1 {
				return &MinimalParamsValidationError{"one of subnets is missing 'name' field or name is empty"}
			}
			if s.AddressPrefixes == nil || len(s.AddressPrefixes) < 1 {
				return &MinimalParamsValidationError{"'address_prefixes' list parameter in one of subnets missing or is 0 length"}
			}
		}
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
