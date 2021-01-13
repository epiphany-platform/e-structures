package v0

import (
	"testing"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		{
			name: "unknown field in config",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.2",
	"azbi": {
		"status": "initialized",
		"config": {
			"unknown_key": "unknown_value", 
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "main", 
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					}
				}],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"rg_name": "epiphany-rg",
			"vm_groups": [
				{
					"vm_group_name": "vm-group0",
					"vms": [
						{
							"private_ips": [
								"10.0.1.4"
							],
							"public_ip": "123.234.345.456",
							"vm_name": "epiphany-vm-group0-0"
						},
						{
							"private_ips": [
								"10.0.1.5"
							],
							"public_ip": "123.234.345.457",
							"vm_name": "epiphany-vm-group0-1"
						},
						{
							"private_ips": [
								"10.0.1.6"
							],
							"public_ip": "123.234.345.458",
							"vm_name": "epiphany-vm-group0-2"
						}
					]
				}
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.2"),
				Unused:  []string{"azbi.config.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:         to.StrPtr("epiphany"),
							Location:     to.StrPtr("northeurope"),
							AddressSpace: []string{"10.0.0.0/16"},
							Subnets: []azbi.Subnet{
								{
									Name:            to.StrPtr("main"),
									AddressPrefixes: []string{"10.0.1.0/24"},
								},
							},
							VmGroups: []azbi.VmGroup{
								{
									Name:        to.StrPtr("vm-group0"),
									VmCount:     to.IntPtr(3),
									VmSize:      to.StrPtr("Standard_DS2_v2"),
									UsePublicIP: to.BooPtr(true),
									SubnetNames: []string{"main"},
									Image: &azbi.Image{
										Publisher: to.StrPtr("Canonical"),
										Offer:     to.StrPtr("UbuntuServer"),
										Sku:       to.StrPtr("18.04-LTS"),
										Version:   to.StrPtr("18.04.202006101"),
									},
								},
							},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						RgName:   to.StrPtr("epiphany-rg"),
						VnetName: to.StrPtr("epiphany-vnet"),
						VmGroups: []azbi.OutputVmGroup{
							{
								Name: to.StrPtr("vm-group0"),
								Vms: []azbi.OutputVm{
									{
										Name:       to.StrPtr("epiphany-vm-group0-0"),
										PrivateIps: []string{"10.0.1.4"},
										PublicIp:   to.StrPtr("123.234.345.456"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-1"),
										PrivateIps: []string{"10.0.1.5"},
										PublicIp:   to.StrPtr("123.234.345.457"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-2"),
										PrivateIps: []string{"10.0.1.6"},
										PublicIp:   to.StrPtr("123.234.345.458"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in config params",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.2",
	"azbi": {
		"status": "initialized",
		"config": {
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"unknown_key": "unknown_value",
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "main", 
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					}
				}],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"rg_name": "epiphany-rg",
			"vm_groups": [
				{
					"vm_group_name": "vm-group0",
					"vms": [
						{
							"private_ips": [
								"10.0.1.4"
							],
							"public_ip": "123.234.345.456",
							"vm_name": "epiphany-vm-group0-0"
						},
						{
							"private_ips": [
								"10.0.1.5"
							],
							"public_ip": "123.234.345.457",
							"vm_name": "epiphany-vm-group0-1"
						},
						{
							"private_ips": [
								"10.0.1.6"
							],
							"public_ip": "123.234.345.458",
							"vm_name": "epiphany-vm-group0-2"
						}
					]
				}
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.2"),
				Unused:  []string{"azbi.config.params.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:         to.StrPtr("epiphany"),
							Location:     to.StrPtr("northeurope"),
							AddressSpace: []string{"10.0.0.0/16"},
							Subnets: []azbi.Subnet{
								{
									Name:            to.StrPtr("main"),
									AddressPrefixes: []string{"10.0.1.0/24"},
								},
							},
							VmGroups: []azbi.VmGroup{
								{
									Name:        to.StrPtr("vm-group0"),
									VmCount:     to.IntPtr(3),
									VmSize:      to.StrPtr("Standard_DS2_v2"),
									UsePublicIP: to.BooPtr(true),
									SubnetNames: []string{"main"},
									Image: &azbi.Image{
										Publisher: to.StrPtr("Canonical"),
										Offer:     to.StrPtr("UbuntuServer"),
										Sku:       to.StrPtr("18.04-LTS"),
										Version:   to.StrPtr("18.04.202006101"),
									},
								},
							},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						RgName:   to.StrPtr("epiphany-rg"),
						VnetName: to.StrPtr("epiphany-vnet"),
						VmGroups: []azbi.OutputVmGroup{
							{
								Name: to.StrPtr("vm-group0"),
								Vms: []azbi.OutputVm{
									{
										Name:       to.StrPtr("epiphany-vm-group0-0"),
										PrivateIps: []string{"10.0.1.4"},
										PublicIp:   to.StrPtr("123.234.345.456"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-1"),
										PrivateIps: []string{"10.0.1.5"},
										PublicIp:   to.StrPtr("123.234.345.457"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-2"),
										PrivateIps: []string{"10.0.1.6"},
										PublicIp:   to.StrPtr("123.234.345.458"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in output",
			args: []byte(`{
	"kind": "state",
	"version": "0.0.2",
	"azbi": {
		"status": "initialized",
		"config": { 
			"kind": "azbi",
			"version": "0.0.1",
			"params": {
				"name": "epiphany",
				"location": "northeurope",
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "main", 
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					}
				}],
				"rsa_pub_path": "/shared/vms_rsa.pub"
			}
		},
		"output": {
			"rg_name": "epiphany-rg",
			"unknown_key": "unknown_value",
			"vm_groups": [
				{
					"vm_group_name": "vm-group0",
					"vms": [
						{
							"private_ips": [
								"10.0.1.4"
							],
							"public_ip": "123.234.345.456",
							"vm_name": "epiphany-vm-group0-0"
						},
						{
							"private_ips": [
								"10.0.1.5"
							],
							"public_ip": "123.234.345.457",
							"vm_name": "epiphany-vm-group0-1"
						},
						{
							"private_ips": [
								"10.0.1.6"
							],
							"public_ip": "123.234.345.458",
							"vm_name": "epiphany-vm-group0-2"
						}
					]
				}
			],
			"vnet_name": "epiphany-vnet"
		}
	}
}`),
			want: &State{
				Kind:    to.StrPtr("state"),
				Version: to.StrPtr("0.0.2"),
				Unused:  []string{"azbi.output.unknown_key"},
				AzBI: &AzBIState{
					Status: "initialized",
					Config: &azbi.Config{
						Kind:    to.StrPtr("azbi"),
						Version: to.StrPtr("0.0.1"),
						Params: &azbi.Params{
							Name:         to.StrPtr("epiphany"),
							Location:     to.StrPtr("northeurope"),
							AddressSpace: []string{"10.0.0.0/16"},
							Subnets: []azbi.Subnet{
								{
									Name:            to.StrPtr("main"),
									AddressPrefixes: []string{"10.0.1.0/24"},
								},
							},
							VmGroups: []azbi.VmGroup{
								{
									Name:        to.StrPtr("vm-group0"),
									VmCount:     to.IntPtr(3),
									VmSize:      to.StrPtr("Standard_DS2_v2"),
									UsePublicIP: to.BooPtr(true),
									SubnetNames: []string{"main"},
									Image: &azbi.Image{
										Publisher: to.StrPtr("Canonical"),
										Offer:     to.StrPtr("UbuntuServer"),
										Sku:       to.StrPtr("18.04-LTS"),
										Version:   to.StrPtr("18.04.202006101"),
									},
								},
							},
							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
						},
						Unused: nil,
					},
					Output: &azbi.Output{
						RgName:   to.StrPtr("epiphany-rg"),
						VnetName: to.StrPtr("epiphany-vnet"),
						VmGroups: []azbi.OutputVmGroup{
							{
								Name: to.StrPtr("vm-group0"),
								Vms: []azbi.OutputVm{
									{
										Name:       to.StrPtr("epiphany-vm-group0-0"),
										PrivateIps: []string{"10.0.1.4"},
										PublicIp:   to.StrPtr("123.234.345.456"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-1"),
										PrivateIps: []string{"10.0.1.5"},
										PublicIp:   to.StrPtr("123.234.345.457"),
									},
									{
										Name:       to.StrPtr("epiphany-vm-group0-2"),
										PrivateIps: []string{"10.0.1.6"},
										PublicIp:   to.StrPtr("123.234.345.458"),
									},
								},
							},
						},
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
			err := got.Unmarshal(tt.args)

			if tt.wantErr != nil {
				if diff := cmp.Diff(tt.wantErr, err, cmpopts.EquateErrors()); diff != "" {
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
