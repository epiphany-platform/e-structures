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
	version = "v0.1.0"
)

type Subnet struct {
	Name            *string  `json:"name"`
	AddressPrefixes []string `json:"address_prefixes"`
}

type Image struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}

type VmGroup struct {
	Name        *string  `json:"name"`
	VmCount     *int     `json:"vm_count"`
	VmSize      *string  `json:"vm_size"`
	UsePublicIP *bool    `json:"use_public_ip"`
	SubnetNames []string `json:"subnet_names"`
	Image       *Image   `json:"image"`
}

type Params struct {
	Name             *string   `json:"name"`
	Location         *string   `json:"location"`
	AddressSpace     []string  `json:"address_space"`
	Subnets          []Subnet  `json:"subnets"`
	VmGroups         []VmGroup `json:"vm_groups"`
	RsaPublicKeyPath *string   `json:"rsa_pub_path"`
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
			Name:             to.StrPtr("epiphany"),
			Location:         to.StrPtr("northeurope"),
			AddressSpace:     []string{"10.0.0.0/16"},
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
			Subnets: []Subnet{
				{
					Name:            to.StrPtr("main"),
					AddressPrefixes: []string{"10.0.1.0/24"},
				},
			},
			VmGroups: []VmGroup{
				{
					Name:        to.StrPtr("vm-group0"),
					VmCount:     to.IntPtr(1),
					VmSize:      to.StrPtr("Standard_DS2_v2"),
					UsePublicIP: to.BooPtr(true),
					SubnetNames: []string{"main"},
					Image: &Image{
						Publisher: to.StrPtr("Canonical"),
						Offer:     to.StrPtr("UbuntuServer"),
						Sku:       to.StrPtr("18.04-LTS"),
						Version:   to.StrPtr("18.04.202006101"),
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
		if c.Params.VmGroups == nil {
			return &MinimalParamsValidationError{"'vm_groups' list parameter missing"}
		}
		for _, s := range c.Params.Subnets {
			if s.Name == nil || len(*s.Name) < 1 {
				return &MinimalParamsValidationError{"one of subnets is missing 'name' field or name is empty"}
			}
			if s.AddressPrefixes == nil || len(s.AddressPrefixes) < 1 {
				return &MinimalParamsValidationError{"'address_prefixes' list parameter in one of subnets missing or is 0 length"}
			}
		}
		if len(c.Params.VmGroups) > 0 {
			for _, vmGroup := range c.Params.VmGroups {
				if vmGroup.Name == nil || len(*vmGroup.Name) < 1 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'name' field or name is empty"}
				}
				if vmGroup.VmCount == nil || *vmGroup.VmCount < 0 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'vm_count' field or there is a negative number"}
				}
				if vmGroup.VmSize == nil || len(*vmGroup.VmSize) < 1 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'vm_size' field or vm_size is empty"}
				}
				if vmGroup.UsePublicIP == nil {
					return &MinimalParamsValidationError{"one of vm groups is missing 'use_public_ip' field"}
				}
				if vmGroup.VmSize == nil || len(*vmGroup.VmSize) < 1 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'vm_size' field or vm_size is empty"}
				}
				if vmGroup.SubnetNames == nil || len(vmGroup.SubnetNames) < 1 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'subnet_names' list field or its length is 0"}
				}
				if vmGroup.Image == nil {
					return &MinimalParamsValidationError{"one of vm groups is missing 'image' field"}
				} else {
					if vmGroup.Image.Publisher == nil || len(*vmGroup.Image.Publisher) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'image.publisher' field or this field is empty"}
					}
					if vmGroup.Image.Offer == nil || len(*vmGroup.Image.Offer) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'image.offer' field or this field is empty"}
					}
					if vmGroup.Image.Sku == nil || len(*vmGroup.Image.Sku) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'image.sku' field or this field is empty"}
					}
					if vmGroup.Image.Version == nil || len(*vmGroup.Image.Version) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'image.version' field or this field is empty"}
					}
				}
			}
		}
	}
	return nil
}

type OutputVm struct {
	Name       string   `json:"vm_name"`
	PrivateIps []string `json:"private_ips"`
	PublicIp   string   `json:"public_ip"`
}

type OutputVmGroup struct {
	Name string     `json:"vm_group_name"`
	Vms  []OutputVm `json:"vms"`
}

type Output struct {
	RgName   *string         `json:"rg_name"`
	VnetName *string         `json:"vnet_name"`
	VmGroups []OutputVmGroup `json:"vm_groups"`
}
