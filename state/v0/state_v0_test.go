package v0

import (
	"testing"

	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
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
			},
			wantErr: nil,
		},
		//		{
		//			name: "unknown fields in multiple places",
		//			args: []byte(`{
		//	"kind": "state",
		//	"version": "0.0.2",
		//	"azbi": {
		//		"status": "initialized",
		//		"config": {
		//			"unknown_key_1": "unknown_value_1",
		//			"meta": {
		//				"kind": "azbiConfig",
		//				"version": "v0.1.0",
		//				"module_version": "v0.0.1"
		//			},
		//			"params": {
		//				"unknown_key_2": "unknown_value_2",
		//				"name": "epiphany",
		//				"location": "northeurope",
		//				"address_space": [
		//					"10.0.0.0/16"
		//				],
		//				"subnets": [
		//					{
		//						"name": "main",
		//						"address_prefixes": [
		//							"10.0.1.0/24"
		//						]
		//					}
		//				],
		//				"vm_groups": [{
		//					"unknown_key_3": "unknown_value_3",
		//					"name": "vm-group0",
		//					"vm_count": 3,
		//					"vm_size": "Standard_DS2_v2",
		//					"use_public_ip": true,
		//					"subnet_names": ["main"],
		//					"vm_image": {
		//						"publisher": "Canonical",
		//						"offer": "UbuntuServer",
		//						"sku": "18.04-LTS",
		//						"version": "18.04.202006101"
		//					},
		//					"data_disks": []
		//				}],
		//				"rsa_pub_path": "/shared/vms_rsa.pub"
		//			}
		//		},
		//		"output": {
		//			"unknown_key_4": "unknown_value_4",
		//			"rg_name": "epiphany-rg",
		//			"vm_groups": [
		//				{
		//					"vm_group_name": "vm-group0",
		//					"unknown_key_5": "unknown_value_5",
		//					"vms": [
		//						{
		//							"private_ips": [
		//								"10.0.1.4"
		//							],
		//							"public_ip": "123.234.345.456",
		//							"vm_name": "epiphany-vm-group0-0"
		//						},
		//						{
		//							"private_ips": [
		//								"10.0.1.5"
		//							],
		//							"public_ip": "123.234.345.457",
		//							"vm_name": "epiphany-vm-group0-1"
		//						},
		//						{
		//							"private_ips": [
		//								"10.0.1.6"
		//							],
		//							"public_ip": "123.234.345.458",
		//							"vm_name": "epiphany-vm-group0-2"
		//						}
		//					]
		//				}
		//			],
		//			"vnet_name": "epiphany-vnet"
		//		}
		//	}
		//}`),
		//			want: &State{
		//				Kind:    to.StrPtr("state"),
		//				Version: to.StrPtr("0.0.2"),
		//				Unused: []string{
		//					"azbi.config.params.vm_groups[0].unknown_key_3",
		//					"azbi.config.params.unknown_key_2",
		//					"azbi.config.unknown_key_1",
		//					"azbi.output.vm_groups[0].unknown_key_5",
		//					"azbi.output.unknown_key_4",
		//				},
		//				AzBI: &AzBIState{
		//					Status: "initialized",
		//					Config: &azbi.Config{
		//						Meta: &azbi.Meta{
		//							Kind:          to.StrPtr("azbiConfig"),
		//							Version:       to.StrPtr("v0.1.0"),
		//							ModuleVersion: to.StrPtr("v0.0.1"),
		//						},
		//						Params: &azbi.Params{
		//							Name:         to.StrPtr("epiphany"),
		//							Location:     to.StrPtr("northeurope"),
		//							AddressSpace: []string{"10.0.0.0/16"},
		//							Subnets: []azbi.Subnet{
		//								{
		//									Name:            to.StrPtr("main"),
		//									AddressPrefixes: []string{"10.0.1.0/24"},
		//								},
		//							},
		//							VmGroups: []azbi.VmGroup{
		//								{
		//									Name:        to.StrPtr("vm-group0"),
		//									VmCount:     to.IntPtr(3),
		//									VmSize:      to.StrPtr("Standard_DS2_v2"),
		//									UsePublicIP: to.BooPtr(true),
		//									SubnetNames: []string{"main"},
		//									VmImage: &azbi.VmImage{
		//										Publisher: to.StrPtr("Canonical"),
		//										Offer:     to.StrPtr("UbuntuServer"),
		//										Sku:       to.StrPtr("18.04-LTS"),
		//										Version:   to.StrPtr("18.04.202006101"),
		//									},
		//									DataDisks: []azbi.DataDisk{},
		//								},
		//							},
		//							RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
		//						},
		//						Unused: nil,
		//					},
		//					Output: &azbi.Output{
		//						RgName:   to.StrPtr("epiphany-rg"),
		//						VnetName: to.StrPtr("epiphany-vnet"),
		//						VmGroups: []azbi.OutputVmGroup{
		//							{
		//								Name: to.StrPtr("vm-group0"),
		//								Vms: []azbi.OutputVm{
		//									{
		//										Name:       to.StrPtr("epiphany-vm-group0-0"),
		//										PrivateIps: []string{"10.0.1.4"},
		//										PublicIp:   to.StrPtr("123.234.345.456"),
		//									},
		//									{
		//										Name:       to.StrPtr("epiphany-vm-group0-1"),
		//										PrivateIps: []string{"10.0.1.5"},
		//										PublicIp:   to.StrPtr("123.234.345.457"),
		//									},
		//									{
		//										Name:       to.StrPtr("epiphany-vm-group0-2"),
		//										PrivateIps: []string{"10.0.1.6"},
		//										PublicIp:   to.StrPtr("123.234.345.458"),
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//			wantErr: nil,
		//		},
		{
			name: "state major version mismatch",
			args: []byte(`{
	"kind": "state",
	"version": "1.0.0"
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		//		{
		//			name: "state minor version mismatch",
		//			args: []byte(`{
		//	"kind": "state",
		//	"version": "0.1.0"
		//}`),
		//			want: &State{
		//				Kind:    to.StrPtr("state"),
		//				Version: to.StrPtr("0.1.0"),
		//				Unused:  []string{},
		//				AzBI:    nil,
		//			},
		//			wantErr: nil,
		//		},
		//		{
		//			name: "state patch version mismatch",
		//			args: []byte(`{
		//	"kind": "state",
		//	"version": "0.0.2"
		//}`),
		//			want: &State{
		//				Kind:    to.StrPtr("state"),
		//				Version: to.StrPtr("0.0.2"),
		//				Unused:  []string{},
		//				AzBI:    nil,
		//			},
		//			wantErr: nil,
		//		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &State{}
			err := got.Unmarshal(tt.args)

			if tt.wantErr != nil {
				if err != nil {
					if _, ok := err.(*validator.InvalidValidationError); ok {
						t.Fatal(err)
					}
					errs := err.(validator.ValidationErrors)
					if len(errs) != len(tt.wantErr.(test.TestValidationErrors)) {
						t.Fatalf("incorrect length of found errors. Got: \n%s\nExpected: \n%s", errs.Error(), tt.wantErr.Error())
					}
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
				} else {
					t.Errorf("No errors got. All expected errors: \n%s", tt.wantErr.Error())
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
