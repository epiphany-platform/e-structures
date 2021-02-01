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
	RsaPublicKeyPath *string `json:"rsa_pub_path"` // TODO check why this field is not validated

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
	KubeDashboardEnabled *bool   `json:"kube_dashboard_enabled"` // TODO remove https://docs.microsoft.com/en-us/azure/aks/kubernetes-dashboard
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
			Location:         to.StrPtr("northeurope"), //TODO possibly delete this value in future
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),

			RgName:     to.StrPtr("epiphany-rg"),
			VnetName:   to.StrPtr("epiphany-vnet"),
			SubnetName: to.StrPtr("azks"),

			KubernetesVersion:  to.StrPtr("1.18.14"), //TODO ensure that this default version is correct
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
		if c.Params.Name == nil || *c.Params.Name == "" {
			return &MinimalParamsValidationError{"'name' parameter missing"}
		}
		if c.Params.Location == nil || *c.Params.Location == "" {
			return &MinimalParamsValidationError{"'location' parameter missing"}
		}
		if c.Params.RgName == nil || *c.Params.RgName == "" {
			return &MinimalParamsValidationError{"'rg_name' parameter missing"}
		}
		if c.Params.VnetName == nil || *c.Params.VnetName == "" {
			return &MinimalParamsValidationError{"'vnet_name' parameter missing"}
		}
		if c.Params.SubnetName == nil || *c.Params.SubnetName == "" {
			return &MinimalParamsValidationError{"'subnet_name' parameter missing"}
		}
		if c.Params.KubernetesVersion == nil || *c.Params.KubernetesVersion == "" { // TODO semver format could also be validated
			return &MinimalParamsValidationError{"'kubernetes_version' parameter missing"}
		}
		if c.Params.EnableNodePublicIp == nil {
			return &MinimalParamsValidationError{"'enable_node_public_ip' parameter missing"}
		}
		if c.Params.EnableRbac == nil {
			return &MinimalParamsValidationError{"'enable_rbac' parameter missing"}
		}

		if c.Params.IdentityType == nil || *c.Params.IdentityType == "" {
			return &MinimalParamsValidationError{"'identity_type' parameter missing"}
		}
		if c.Params.KubeDashboardEnabled == nil {
			return &MinimalParamsValidationError{"'kube_dashboard_enabled' parameter missing"}
		}
		fmt.Println("[DEPRECATION] 'kube_dashboard_enabled' parameter will soon be deprecated due to Azure removing Dashboard support in AKS.")
		if c.Params.AdminUsername == nil || *c.Params.AdminUsername == "" {
			return &MinimalParamsValidationError{"'admin_username' parameter missing"}
		}

		if c.Params.DefaultNodePool == nil {
			return &MinimalParamsValidationError{"'default_node_pool' parameter missing"}
		} else {
			if c.Params.DefaultNodePool.Size == nil {
				return &MinimalParamsValidationError{"'default_node_pool.size' parameter missing"}
			}
			if c.Params.DefaultNodePool.Min == nil {
				return &MinimalParamsValidationError{"'default_node_pool.min' parameter missing"}
			}
			if c.Params.DefaultNodePool.Max == nil {
				return &MinimalParamsValidationError{"'default_node_pool.max' parameter missing"}
			}
			if c.Params.DefaultNodePool.VmSize == nil || *c.Params.DefaultNodePool.VmSize == "" {
				return &MinimalParamsValidationError{"'default_node_pool.vm_size' parameter missing"}
			}
			if c.Params.DefaultNodePool.DiskSize == nil || *c.Params.DefaultNodePool.DiskSize == "" {
				return &MinimalParamsValidationError{"'default_node_pool.disk_size' parameter missing"}
			}
			if c.Params.DefaultNodePool.AutoScaling == nil {
				return &MinimalParamsValidationError{"'default_node_pool.auto_scaling' parameter missing"}
			}
			if c.Params.DefaultNodePool.Type == nil || *c.Params.DefaultNodePool.Type == "" {
				return &MinimalParamsValidationError{"'default_node_pool.type' parameter missing"}
			}
		}

		if c.Params.AutoScalerProfile == nil {
			return &MinimalParamsValidationError{"'auto_scaler_profile' parameter missing"}
		} else {
			if c.Params.AutoScalerProfile.BalanceSimilarNodeGroups == nil {
				return &MinimalParamsValidationError{"'auto_scaler_profile.balance_similar_node_groups' parameter missing"}
			}
			if c.Params.AutoScalerProfile.MaxGracefulTerminationSec == nil || *c.Params.AutoScalerProfile.MaxGracefulTerminationSec == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.max_graceful_termination_sec' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownDelayAfterAdd == nil || *c.Params.AutoScalerProfile.ScaleDownDelayAfterAdd == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_add' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownDelayAfterDelete == nil || *c.Params.AutoScalerProfile.ScaleDownDelayAfterDelete == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_delete' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownDelayAfterFailure == nil || *c.Params.AutoScalerProfile.ScaleDownDelayAfterFailure == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_failure' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScanInterval == nil || *c.Params.AutoScalerProfile.ScanInterval == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scan_interval' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownUnneeded == nil || *c.Params.AutoScalerProfile.ScaleDownUnneeded == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unneeded' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownUnready == nil || *c.Params.AutoScalerProfile.ScaleDownUnready == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unready' parameter missing"}
			}
			if c.Params.AutoScalerProfile.ScaleDownUtilizationThreshold == nil || *c.Params.AutoScalerProfile.ScaleDownUtilizationThreshold == "" { // TODO format could also be validated
				return &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_utilization_threshold' parameter missing"}
			}
		}

		if c.Params.AzureAd == nil {
			// Azure AD can be null
		} else {
			if c.Params.AzureAd.Managed == nil {
				return &MinimalParamsValidationError{"'azure_ad.managed' parameter missing"}
			}
			if c.Params.AzureAd.TenantId == nil || *c.Params.AzureAd.TenantId == "" {
				return &MinimalParamsValidationError{"'azure_ad.tenant_id' parameter missing"}
			}
			if c.Params.AzureAd.AdminGroupObjectIds == nil || len(c.Params.AzureAd.AdminGroupObjectIds) < 1 {
				return &MinimalParamsValidationError{"'azure_ad.admin_group_object_ids' parameter list is missing or its length is 0"}
			}
			for _, ago := range c.Params.AzureAd.AdminGroupObjectIds {
				if ago == "" {
					return &MinimalParamsValidationError{"one of Azure AD Admin Group IDs lists value is empty"}
				}
			}
		}
	}
	return nil
}

type Output struct {
	KubeConfig *string `json:"kubeconfig"`
}
