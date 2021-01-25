package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/epiphany-platform/e-structures/utils/to"
	maps "github.com/mitchellh/mapstructure"
	"reflect"
)

const (
	kind    = "azks"
	version = "v0.0.1"
)

type AzureAd struct {
	Managed             *bool    `json:"managed"`
	TenantId            *string  `json:"tenant_id"`
	AdminGroupObjectIds []string `json:"admin_group_object_ids"`
}

type AutoScalerProfile struct { //TODO consider changing types of string values here to make it more golang'ish
	BalanceSimilarNodeGroups      *bool   `json:"balance_similar_node_groups"`
	MaxGracefulTerminationSec     *string `json:"max_graceful_termination_sec"`
	ScaleDownDelayAfterAdd        *string `json:"scale_down_delay_after_add"`
	ScaleDownDelayAfterDelete     *string `json:"scale_down_delay_after_delete"`
	ScaleDownDelayAfterFailure    *string `json:"scale_down_delay_after_failure"`
	ScanInterval                  *string `json:"scan_interval"`
	ScaleDownUnneeded             *string `json:"scale_down_unneeded"`
	ScaleDownUnready              *string `json:"scale_down_unready"`
	ScaleDownUtilizationThreshold *string `json:"scale_down_utilization_threshold"`
}

type DefaultNodePool struct {
	Size        *int    `json:"size"`
	Min         *int    `json:"min"`
	Max         *int    `json:"max"`
	VmSize      *string `json:"vm_size"`
	DiskSize    *string `json:"disk_size"`
	AutoScaling *bool   `json:"auto_scaling"`
	Type        *string `json:"type"`
}

type Params struct {
	Name             *string `json:"name"`
	Location         *string `json:"location"`
	RsaPublicKeyPath *string `json:"rsa_pub_path"`

	RgName     *string `json:"rg_name"`
	VnetName   *string `json:"vnet_name"`
	SubnetName *string `json:"subnet_name"`

	KubernetesVersion  *string `json:"kubernetes_version"`
	EnableNodePublicIp *bool   `json:"enable_node_public_ip"`
	EnableRbac         *bool   `json:"enable_rbac"`

	DefaultNodePool   *DefaultNodePool   `json:"default_node_pool"`
	AutoScalerProfile *AutoScalerProfile `json:"auto_scaler_profile"`
	AzureAd           *AzureAd           `json:"azure_ad"`

	IdentityType         *string `json:"identity_type"`
	KubeDashboardEnabled *bool   `json:"kube_dashboard_enabled"`
	AdminUsername        *string `json:"admin_username"`
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
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

			RgName:     to.StrPtr("epiphany-rg"),
			VnetName:   to.StrPtr("epiphany-vnet"),
			SubnetName: to.StrPtr("azks"),

			KubernetesVersion:  to.StrPtr("1.18.14"), //TODO ensure that this makes sense
			EnableNodePublicIp: to.BooPtr(false),
			EnableRbac:         to.BooPtr(false),

			DefaultNodePool: &DefaultNodePool{
				Size:        to.IntPtr(2),
				Min:         to.IntPtr(2),
				Max:         to.IntPtr(5),
				VmSize:      to.StrPtr("Standard_DS2_v2"),
				DiskSize:    to.StrPtr("36"),
				AutoScaling: to.BooPtr(true),
				Type:        to.StrPtr("VirtualMachineScaleSets"),
			},
			AutoScalerProfile: &AutoScalerProfile{
				BalanceSimilarNodeGroups:      to.BooPtr(false),
				MaxGracefulTerminationSec:     to.StrPtr("600"),
				ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
				ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
				ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
				ScanInterval:                  to.StrPtr("10s"),
				ScaleDownUnneeded:             to.StrPtr("10m"),
				ScaleDownUnready:              to.StrPtr("10m"),
				ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
			},
			AzureAd: nil,

			IdentityType:         to.StrPtr("SystemAssigned"),
			KubeDashboardEnabled: to.BooPtr(true),
			AdminUsername:        to.StrPtr("operations"),
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
		//TODO fix
	}
	return nil
}

type Output struct {
	KubeConfig *string `json:"kubeconfig"`
}
