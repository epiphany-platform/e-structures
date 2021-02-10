package v0

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	maps "github.com/mitchellh/mapstructure"
)

const (
	kind    = "azbi"
	version = "v0.1.2"
)

type DataDisk struct {
	GbSize *int `json:"disk_size_gb" validate:"required,min=1"`
}

type Subnet struct {
	Name            *string  `json:"name" validate:"required,min=1"`
	AddressPrefixes []string `json:"address_prefixes" validate:"required,min=1,dive,required"`
}

type VmImage struct {
	Publisher *string `json:"publisher" validate:"required,min=1"`
	Offer     *string `json:"offer" validate:"required,min=1"`
	Sku       *string `json:"sku" validate:"required,min=1"`
	Version   *string `json:"version" validate:"required,min=1"`
}

type VmGroup struct {
	Name        *string    `json:"name" validate:"required"`
	VmCount     *int       `json:"vm_count" validate:"required,min=1"`
	VmSize      *string    `json:"vm_size" validate:"required"`
	UsePublicIP *bool      `json:"use_public_ip" validate:"required"`
	SubnetNames []string   `json:"subnet_names" validate:"required,min=1,dive,required"` // TODO custom validator to ensure subnet existence
	VmImage     *VmImage   `json:"vm_image" validate:"required,dive"`
	DataDisks   []DataDisk `json:"data_disks" validate:"required,dive"`
}

type Params struct {
	Name             *string   `json:"name" validate:"required"`
	Location         *string   `json:"location"`      // TODO add validation
	AddressSpace     []string  `json:"address_space"` // TODO add validation
	Subnets          []Subnet  `json:"subnets" validate:"required,min=1,dive"`
	VmGroups         []VmGroup `json:"vm_groups" validate:"required,dive"`
	RsaPublicKeyPath *string   `json:"rsa_pub_path"` // TODO add validation
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
	Kind    *string  `json:"kind" validate:"required,eq=azbi"`
	Version *string  `json:"version" validate:"required"` // TODO custom validator to ensure major version match
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

func (c *Config) isValid() error {
	if c == nil {
		return errors.New("azbi config is nil")
	}
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		fmt.Println(err.Error())
		return err
	}

	//v, err := semver.NewVersion(version)
	//if err != nil {
	//	return err
	//}
	//constraint, err := semver.NewConstraint(fmt.Sprintf("~%d", v.Major()))
	//if err != nil {
	//	return err
	//}
	//vv, err := semver.NewVersion(*c.Version)
	//if err != nil {
	//	return err
	//}
	//if !constraint.Check(vv) {
	//	return MajorVersionMismatchError
	//}
	//if c.Params != nil && !reflect.DeepEqual(c.Params, &Params{}) {

	//	if c.Params.Location == nil {
	//		return &MinimalParamsValidationError{"'location' parameter missing"}
	//	}
	//	if c.Params.Subnets == nil || len(c.Params.Subnets) < 1 {
	//		return &MinimalParamsValidationError{"'subnets' list parameter missing or is 0 length"}
	//	}
	//	if len(c.Params.VmGroups) > 0 {
	//		for _, vmGroup := range c.Params.VmGroups {
	//			if vmGroup.SubnetNames == nil || len(vmGroup.SubnetNames) < 1 {
	//				return &MinimalParamsValidationError{"one of vm groups is missing 'subnet_names' list field or its length is 0"}
	//			}
	//			for _, sn := range vmGroup.SubnetNames {
	//				if sn == "" {
	//					return &MinimalParamsValidationError{"one of vm groups subnet names lists value is empty"}
	//				}
	//				found := false
	//				for _, s := range c.Params.Subnets {
	//					if sn == *s.Name {
	//						found = true
	//					}
	//				}
	//				if !found {
	//					return &MinimalParamsValidationError{"one of vm groups subnet names wasn't found among subnets"}
	//				}
	//			}
	//		}
	//	}
	//}
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
