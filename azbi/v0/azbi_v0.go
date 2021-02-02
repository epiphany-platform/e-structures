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
	version = "v0.1.2"
)

type DataDisk struct {
	GbSize *int `json:"disk_size_gb"`
}

type Subnet struct {
	Name            *string  `json:"name"`
	AddressPrefixes []string `json:"address_prefixes"`
}

type VmImage struct {
	Publisher *string `json:"publisher"`
	Offer     *string `json:"offer"`
	Sku       *string `json:"sku"`
	Version   *string `json:"version"`
}

type VmGroup struct {
	Name        *string    `json:"name"`
	VmCount     *int       `json:"vm_count"`
	VmSize      *string    `json:"vm_size"`
	UsePublicIP *bool      `json:"use_public_ip"`
	SubnetNames []string   `json:"subnet_names"`
	VmImage     *VmImage   `json:"vm_image"`
	DataDisks   []DataDisk `json:"data_disks"`
}

type Params struct {
	Name             *string   `json:"name"`
	Location         *string   `json:"location"`
	AddressSpace     []string  `json:"address_space"`
	Subnets          []Subnet  `json:"subnets"`
	VmGroups         []VmGroup `json:"vm_groups"`
	RsaPublicKeyPath *string   `json:"rsa_pub_path"` // TODO check why this field is not validated
}

func (p *Params) GetRsaPublicKeyV() string {
	if p == nil {
		return ""
	}
	return *p.RsaPublicKeyPath
}

func (p *Params) GetNameV() string {
	if p == nil {
		return ""
	}
	return *p.Name
}

func (p *Params) GetLocationV() string {
	if p == nil {
		return ""
	}
	return *p.Location
}

// ExtractEmptySubnets gets params and extracts from it list of Subnet unassigned to any of VmGroup.
func (p *Params) ExtractEmptySubnets() []Subnet {
	if p == nil {
		return nil
	}
	if p.Subnets == nil || len(p.Subnets) == 0 {
		return nil
	}
	if p.VmGroups == nil || len(p.VmGroups) == 0 {
		return p.Subnets
	}
	m := make(map[string]Subnet)
	for _, subnet := range p.Subnets {
		m[*subnet.Name] = subnet
	}
	for _, vmGroup := range p.VmGroups {
		for _, subnet := range p.Subnets {
			for _, subnetName := range vmGroup.SubnetNames {
				if *subnet.Name == subnetName {
					_, ok := m[subnetName]
					if ok {
						delete(m, subnetName)
						break
					}
				}
			}
		}
	}
	result := make([]Subnet, 0)
	for _, s := range p.Subnets {
		if v, ok := m[*s.Name]; ok {
			result = append(result, v)
		}
	}
	return result
}

type Config struct {
	Kind    *string  `json:"kind"`
	Version *string  `json:"version"`
	Params  *Params  `json:"params"`
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
					VmImage: &VmImage{
						Publisher: to.StrPtr("Canonical"),
						Offer:     to.StrPtr("UbuntuServer"),
						Sku:       to.StrPtr("18.04-LTS"),
						Version:   to.StrPtr("18.04.202006101"),
					},
					DataDisks: []DataDisk{
						{
							GbSize: to.IntPtr(10),
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
		if c.Params.Name == nil {
			return &MinimalParamsValidationError{"'name' parameter missing"}
		}
		if c.Params.Location == nil {
			return &MinimalParamsValidationError{"'location' parameter missing"}
		}
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
			for _, ap := range s.AddressPrefixes {
				if ap == "" {
					return &MinimalParamsValidationError{"'address_prefixes' list value in one of subnets missing is empty"}
				}
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
				if vmGroup.SubnetNames == nil || len(vmGroup.SubnetNames) < 1 {
					return &MinimalParamsValidationError{"one of vm groups is missing 'subnet_names' list field or its length is 0"}
				}
				for _, sn := range vmGroup.SubnetNames {
					if sn == "" {
						return &MinimalParamsValidationError{"one of vm groups subnet names lists value is empty"}
					}
					found := false
					for _, s := range c.Params.Subnets {
						if sn == *s.Name {
							found = true
						}
					}
					if !found {
						return &MinimalParamsValidationError{"one of vm groups subnet names wasn't found among subnets"}
					}
				}
				if vmGroup.VmImage == nil {
					return &MinimalParamsValidationError{"one of vm groups is missing 'vm_image' field"}
				} else {
					if vmGroup.VmImage.Publisher == nil || len(*vmGroup.VmImage.Publisher) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'vm_image.publisher' field or this field is empty"}
					}
					if vmGroup.VmImage.Offer == nil || len(*vmGroup.VmImage.Offer) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'vm_image.offer' field or this field is empty"}
					}
					if vmGroup.VmImage.Sku == nil || len(*vmGroup.VmImage.Sku) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'vm_image.sku' field or this field is empty"}
					}
					if vmGroup.VmImage.Version == nil || len(*vmGroup.VmImage.Version) < 1 {
						return &MinimalParamsValidationError{"one of vm groups is missing 'vm_image.version' field or this field is empty"}
					}
				}
				if vmGroup.DataDisks == nil {
					return &MinimalParamsValidationError{"one of vm groups is missing 'data_disks' list"}
				}
				for _, dd := range vmGroup.DataDisks {
					if dd.GbSize == nil || *dd.GbSize < 1 {
						return &MinimalParamsValidationError{"one of vm groups data disks sizes is empty or size is less than 1"}
					}
				}
			}
		}
	}
	return nil
}

type OutputDataDisk struct {
	Size *int `json:"size"`
	Lun  *int `json:"lun"`
}

type OutputVm struct {
	Name       *string          `json:"vm_name"`
	PrivateIps []string         `json:"private_ips"`
	PublicIp   *string          `json:"public_ip"`
	DataDisks  []OutputDataDisk `json:"data_disks"`
}

type OutputVmGroup struct {
	Name *string    `json:"vm_group_name"`
	Vms  []OutputVm `json:"vms"`
}

type Output struct {
	RgName   *string         `json:"rg_name"`
	VnetName *string         `json:"vnet_name"`
	VmGroups []OutputVmGroup `json:"vm_groups"`
}

func (o *Output) GetRgNameV() string {
	if o == nil {
		return ""
	}
	return *o.RgName
}

func (o *Output) GetVnetNameV() string {
	if o == nil {
		return ""
	}
	return *o.VnetName
}
