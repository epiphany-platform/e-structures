package v0

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	azbi "github.com/mkyc/go-stucts-versioning-tests/azbi/v0"
	"github.com/mkyc/go-stucts-versioning-tests/utils/to"
	"testing"
)

func TestState_Load(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    *State
		wantErr error
	}{
		{
			name: "minimal state",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.1"
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.1"),
				Unused:  []string{},
				AzBI:    nil,
			},
			wantErr: nil,
		},
		{
			name: "unknown field in config",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.1",
	"azbi": {
		"status": "initialized",
		"config": {
			"unknown_key": "unknown_value", 
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"vms_count": 3,
				"use_public_ip": true,
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"address_prefixes": [
					"10.0.1.0/24"
				],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"private_ips": [],
			"public_ips": [
				"123.234.345.456",
				"234.345.456.567",
				"345.456.567.678"
			],
			"rg_name": "epiphany-rg",
			"vm_names": [
				"vm1",
				"vm2",
				"vm3"
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.1"),
				Unused:  []string{"azbi.config.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:             to.StrPtr("epiphany"),
							VmsCount:         to.IntPtr(3),
							UsePublicIP:      to.BooPtr(true),
							Location:         to.StrPtr("northeurope"),
							AddressSpace:     []string{"10.0.0.0/16"},
							AddressPrefixes:  []string{"10.0.1.0/24"},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						PrivateIps: []string{},
						PublicIps:  []string{"123.234.345.456", "234.345.456.567", "345.456.567.678"},
						RgName:     to.StrPtr("epiphany-rg"),
						VmNames:    []string{"vm1", "vm2", "vm3"},
						VnetName:   to.StrPtr("epiphany-vnet"),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in config params",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.1",
	"azbi": {
		"status": "initialized",
		"config": {
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"unknown_key": "unknown_value", 
				"vms_count": 3,
				"use_public_ip": true,
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"address_prefixes": [
					"10.0.1.0/24"
				],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"private_ips": [],
			"public_ips": [
				"123.234.345.456",
				"234.345.456.567",
				"345.456.567.678"
			],
			"rg_name": "epiphany-rg",
			"vm_names": [
				"vm1",
				"vm2",
				"vm3"
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.1"),
				Unused:  []string{"azbi.config.params.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:             to.StrPtr("epiphany"),
							VmsCount:         to.IntPtr(3),
							UsePublicIP:      to.BooPtr(true),
							Location:         to.StrPtr("northeurope"),
							AddressSpace:     []string{"10.0.0.0/16"},
							AddressPrefixes:  []string{"10.0.1.0/24"},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						PrivateIps: []string{},
						PublicIps:  []string{"123.234.345.456", "234.345.456.567", "345.456.567.678"},
						RgName:     to.StrPtr("epiphany-rg"),
						VmNames:    []string{"vm1", "vm2", "vm3"},
						VnetName:   to.StrPtr("epiphany-vnet"),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in output",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.1",
	"azbi": {
		"status": "initialized",
		"config": { 
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"vms_count": 3,
				"use_public_ip": true,
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"address_prefixes": [
					"10.0.1.0/24"
				],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"private_ips": [],
			"unknown_key": "unknown_value",
			"public_ips": [
				"123.234.345.456",
				"234.345.456.567",
				"345.456.567.678"
			],
			"rg_name": "epiphany-rg",
			"vm_names": [
				"vm1",
				"vm2",
				"vm3"
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.1"),
				Unused:  []string{"azbi.output.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:             to.StrPtr("epiphany"),
							VmsCount:         to.IntPtr(3),
							UsePublicIP:      to.BooPtr(true),
							Location:         to.StrPtr("northeurope"),
							AddressSpace:     []string{"10.0.0.0/16"},
							AddressPrefixes:  []string{"10.0.1.0/24"},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						PrivateIps: []string{},
						PublicIps:  []string{"123.234.345.456", "234.345.456.567", "345.456.567.678"},
						RgName:     to.StrPtr("epiphany-rg"),
						VmNames:    []string{"vm1", "vm2", "vm3"},
						VnetName:   to.StrPtr("epiphany-vnet"),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "state major version mismatch",
			args: []byte(`{
	"kind": "state",
	"version": "1.0.0"
}`),
			want:    nil,
			wantErr: MajorVersionMismatchError,
		},
		{
			name: "state minor version mismatch",
			args: []byte(`{
	"kind": "state",
	"version": "0.1.0"
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.1.0"),
				Unused:  []string{},
				AzBI:    nil,
			},
			wantErr: nil,
		},
		{
			name: "state patch version mismatch",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.2"
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.2"),
				Unused:  []string{},
				AzBI:    nil,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &State{}
			err := got.Load(tt.args)

			if tt.wantErr != nil {
				if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("Load() errors mismatch (-want +got):\n%s", diff)
				}
			} else {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("Load() mismatch (-want +got):\n%s", diff)
				}
				if err != nil {
					t.Errorf("Load() unexpected error occured: %v", err)
				}
			}
		})
	}
}
