package v0

import (
	"testing"

	"github.com/epiphany-platform/e-structures/utils/to"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					AzureAd:              nil,
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in main structure",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					AzureAd:              nil,
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{"extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in sub structure",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"extra_inner_field" : "extra_inner_value",
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					AzureAd:              nil,
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{"params.extra_inner_field"},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in all possible places",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"name": "epiphany",
		"extra_inner_field" : "extra_inner_value",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					AzureAd:              nil,
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{"params.extra_inner_field", "extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name:    "empty json",
			args:    []byte(`{}`),
			want:    nil,
			wantErr: KindMissingValidationError,
		},
		{
			name: "just kind field",
			args: []byte(`{
	"kind": "azks"
}
`),
			want:    nil,
			wantErr: VersionMissingValidationError,
		},
		{
			name: "just kind and version",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1"
}
`),
			want:    nil,
			wantErr: ParamsMissingValidationError,
		},
		{
			name: "minimal correct json",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "full json",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					AzureAd: &AzureAd{
						Managed:             to.BooPtr(true),
						TenantId:            to.StrPtr("123123123123"),
						AdminGroupObjectIds: []string{"123123123123"},
					},
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'name' parameter missing"},
		},
		{
			name: "empty name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'name' parameter missing"},
		},
		{
			name: "missing location",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'location' parameter missing"},
		},
		{
			name: "empty location",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'location' parameter missing"},
		},
		{
			name: "missing rg_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'rg_name' parameter missing"},
		},
		{
			name: "empty rg_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'rg_name' parameter missing"},
		},
		{
			name: "missing vnet_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'vnet_name' parameter missing"},
		},
		{
			name: "empty vnet_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'vnet_name' parameter missing"},
		},
		{
			name: "missing subnet_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'subnet_name' parameter missing"},
		},
		{
			name: "empty subnet_name",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'subnet_name' parameter missing"},
		},
		{
			name: "missing kubernetes_version",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'kubernetes_version' parameter missing"},
		},
		{
			name: "empty kubernetes_version",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'kubernetes_version' parameter missing"},
		},
		{
			name: "missing enable_node_public_ip",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'enable_node_public_ip' parameter missing"},
		},
		{
			name: "missing enable_rbac",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'enable_rbac' parameter missing"},
		},
		{
			name: "missing identity_type",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'identity_type' parameter missing"},
		},
		{
			name: "empty identity_type",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'identity_type' parameter missing"},
		},
		{
			name: "missing kube_dashboard_enabled",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'kube_dashboard_enabled' parameter missing"},
		},
		{
			name: "missing admin_username",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'admin_username' parameter missing"},
		},
		{
			name: "empty admin_username",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": ""
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'admin_username' parameter missing"},
		},
		{
			name: "missing default_node_pool",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool' parameter missing"},
		},
		{
			name: "empty default_node_pool",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.size' parameter missing"},
		},
		{
			name: "missing default_node_pool.size",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.size' parameter missing"},
		},
		{
			name: "missing default_node_pool.min",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.min' parameter missing"},
		},
		{
			name: "missing default_node_pool.max",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.max' parameter missing"},
		},
		{
			name: "missing default_node_pool.vm_size",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.vm_size' parameter missing"},
		},
		{
			name: "empty default_node_pool.vm_size",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.vm_size' parameter missing"},
		},
		{
			name: "missing default_node_pool.disk_size",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.disk_size' parameter missing"},
		},
		{
			name: "empty default_node_pool.disk_size",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.disk_size' parameter missing"},
		},
		{
			name: "missing default_node_pool.auto_scaling",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.auto_scaling' parameter missing"},
		},
		{
			name: "missing default_node_pool.type",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.type' parameter missing"},
		},
		{
			name: "empty default_node_pool.type",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": ""
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'default_node_pool.type' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.balance_similar_node_groups' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.balance_similar_node_groups",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.balance_similar_node_groups' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.max_graceful_termination_sec",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.max_graceful_termination_sec' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.max_graceful_termination_sec",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.max_graceful_termination_sec' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_delay_after_add",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_add' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_delay_after_add",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_add' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_delay_after_delete",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_delete' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_delay_after_delete",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_delete' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_delay_after_failure",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_failure' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_delay_after_failure",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_delay_after_failure' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scan_interval",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scan_interval' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scan_interval",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scan_interval' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_unneeded",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unneeded' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_unneeded",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unneeded' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_unready",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unready' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_unready",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_unready' parameter missing"},
		},
		{
			name: "missing auto_scaler_profile.scale_down_utilization_threshold",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_utilization_threshold' parameter missing"},
		},
		{
			name: "empty auto_scaler_profile.scale_down_utilization_threshold",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": ""
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'auto_scaler_profile.scale_down_utilization_threshold' parameter missing"},
		},
		{
			name: "missing azure_ad",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "null azure_ad",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
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
					IdentityType:         to.StrPtr("SystemAssigned"),
					KubeDashboardEnabled: to.BooPtr(true),
					AdminUsername:        to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing azure_ad.managed",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'azure_ad.managed' parameter missing"},
		},
		{
			name: "missing azure_ad.tenant_id",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'azure_ad.tenant_id' parameter missing"},
		},
		{
			name: "empty azure_ad.tenant_id",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'azure_ad.tenant_id' parameter missing"},
		},
		{
			name: "missing azure_ad.admin_group_object_ids",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123"
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'azure_ad.admin_group_object_ids' parameter list is missing or its length is 0"},
		},
		{
			name: "empty azure_ad.admin_group_object_ids",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": []
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'azure_ad.admin_group_object_ids' parameter list is missing or its length is 0"},
		},
		{
			name: "empty azure_ad.admin_group_object_ids element",
			args: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_size": "36",
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				""
			]
		}, 
		"identity_type": "SystemAssigned",
		"kube_dashboard_enabled": true,
		"admin_username": "operations"
	}
}`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"one of Azure AD Admin Group IDs lists value is empty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Config{}
			err := got.Unmarshal(tt.args)

			if tt.wantErr != nil {
				errMsg := ""
				if err != nil {
					errMsg = err.Error()
				}
				if diff := cmp.Diff(tt.wantErr.Error(), errMsg, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("Unmarshal() errors mismatch (-want +got):\n%s", diff)
				}
			} else {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
				}
				if err != nil {
					t.Errorf("Unmarshal() unexpected error occured: %v", err)
				}
			}
		})
	}
}
