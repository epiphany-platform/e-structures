package v0

import (
	"reflect"
	"testing"

	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

func TestConfig_Load_general(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
	"kind": "azbi",
	"version": "v0.1.0",
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [{
			"name": "main",
			"address_prefixes": [
				"10.0.1.0/24"
			]
		}],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 3,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}, 
			"data_disks": [
				{
					"disk_size_gb": 10
				}
			]
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							VmCount:     to.IntPtr(3),
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
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"kind": "azbi",
	"version": "v0.1.0",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"extra_inner_field" : "extra_inner_value",
		"location": "northeurope",
		"name": "epiphany",
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
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}, 
			"data_disks": [
				{
					"disk_size_gb": 10
				}
			]
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:     to.StrPtr("northeurope"),
					Name:         to.StrPtr("epiphany"),
					AddressSpace: []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
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
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
				},
				Unused: []string{"params.extra_inner_field", "extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "empty json",
			json: []byte(`{}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Kind",
					Field: "Kind",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "minimal correct json",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 1,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": false,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "major version mismatch",
			json: []byte(`{
			"kind": "azbi",
			"version": "100.0.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			json: []byte(`{
			"kind": "azbi",
			"version": "0.100.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.100.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			json: []byte(`{
			"kind": "azbi",
			"version": "0.0.100",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.100"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

func TestConfig_Load_Params(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "just vm_groups in params",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					}, 
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Location",
					Field: "Location",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing requested subnets list",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "VmGroups[0].SubnetNames[0]",
					Tag:   "insubnets",
				},
			},
		},
		{
			name: "empty subnets list",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing subnet params",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty subnet params",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "",
						"address_prefixes": []
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty subnet address prefixes element and not cidr",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "main",
						"address_prefixes": [
							""
						]
					},
					{
						"name": "main",
						"address_prefixes": [
							"10.0.1.0"
						]
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes[0]",
					Field: "AddressPrefixes[0]",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[1].AddressPrefixes[0]",
					Field: "AddressPrefixes[0]",
					Tag:   "cidr",
				},
			},
		},
		{
			name: "multiple subnets configuration",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "main",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					},
					{
						"name": "second",
						"address_prefixes": [
							"10.0.2.0/24"
						]
					}
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"main"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing address_space when present subnets",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "excluded_without",
				},
			},
		},
		{
			name: "missing subnets when present address_space",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
				"address_space": [
					"10.0.0.0/16"
				],
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "required_with",
				},
			},
		},
		{
			name: "emtpy address_space",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [],
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace",
					Field: "AddressSpace",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty address_space element or not cidr",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"", "10.0.1.0"
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace[0]",
					Field: "AddressSpace[0]",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace[1]",
					Field: "AddressSpace[1]",
					Tag:   "cidr",
				},
			},
		},
		{
			name: "empty params.rsa_pub_path",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "",  
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "min",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

func TestConfig_Load_VmGroups(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "vm_group without networking",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "/shared/vms_rsa.pub",
				"vm_groups": [{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing vm_groups parameter",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name"
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups",
					Field: "VmGroups",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups parameter",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
				"vm_groups": []
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},

		{
			name: "missing vm_groups parameters",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"vm_groups": [{}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmSize",
					Field: "VmSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].UsePublicIP",
					Field: "UsePublicIP",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage",
					Field: "VmImage",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks",
					Field: "DataDisks",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups parameters",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",
				"vm_groups": [{
					"name": "",
					"vm_count": 0,
					"vm_size": "",
					"use_public_ip": true,
					"vm_image": {},
					"data_disks": [{}]
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmSize",
					Field: "VmSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "required",
				},
			},
		},

		{
			name: "negative vm_groups vm_count parameter",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_count": -1,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": ["main"],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty vm_groups subnet_names list",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"subnet_names": [],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames",
					Field: "SubnetNames",
					Tag:   "min",
				},
			},
		},
		{
			name: "vm_groups subnet_names list empty value",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"subnet_names": [""],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "SubnetNames[0]",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "VmGroups[0].SubnetNames[0]",
					Tag:   "required",
				},
			},
		},
		{
			name: "vm_groups subnet_names list value not existing in subnets",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"subnet_names": ["incorrect"],
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "VmGroups[0].SubnetNames[0]",
					Tag:   "insubnets",
				},
			},
		},
		{
			name: "empty vm_groups.data_disks list value",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": [
						{}
					]
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "required",
				},
			},
		},
		{
			name: "zero vm_groups.data_disks list value",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": [
						{
							"disk_size_gb": 0
						}
					]
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
			},
		},
		{
			name: "negative vm_groups.data_disks list value",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_image": {
						"publisher": "Canonical",
						"offer": "UbuntuServer",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": [
						{
							"disk_size_gb": -1
						}
					]
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
			},
		},
		{
			name: "multiple vm_groups configuration",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "first",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				],
				"vm_groups": [
					{
						"name": "first",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["first"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					},
					{
						"name": "second",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["first"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					}
				]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple vm_groups and subnets configuration",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name", 
				"address_space": [
					"10.0.0.0/16"
				], 
				"subnets": [
					{
						"name": "first",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					},
					{
						"name": "second",
						"address_prefixes": [
							"10.0.2.0/24"
						]
					}
				],
				"vm_groups": [
					{
						"name": "first",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["first"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					},
					{
						"name": "second",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["second"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					}
				]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"second"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "2 vm_groups and 3 subnets configuration",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "first",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					},
					{
						"name": "second",
						"address_prefixes": [
							"10.0.2.0/24"
						]
					},
					{
						"name": "third",
						"address_prefixes": [
							"10.0.3.0/24"
						]
					}
				],
				"vm_groups": [
					{
						"name": "first",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["first", "third"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					},
					{
						"name": "second",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["second", "third"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": []
					}
				]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
						{
							Name:            to.StrPtr("third"),
							AddressPrefixes: []string{"10.0.3.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"first", "third"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"second", "third"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple vm_groups and subnets and data disks configuration",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
				"address_space": [
					"10.0.0.0/16"
				],
				"subnets": [
					{
						"name": "first",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					},
					{
						"name": "second",
						"address_prefixes": [
							"10.0.2.0/24"
						]
					}
				],
				"vm_groups": [
					{
						"name": "first",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["first"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": [
							{
								"disk_size_gb": 10
							},
							{
								"disk_size_gb": 20
							}
						]
					},
					{
						"name": "second",
						"vm_count": 3,
						"vm_size": "Standard_DS2_v2",
						"use_public_ip": true,
						"subnet_names": ["second"],
						"vm_image": {
							"publisher": "Canonical",
							"offer": "UbuntuServer",
							"sku": "18.04-LTS",
							"version": "18.04.202006101"
						},
						"data_disks": [
							{
								"disk_size_gb": 30
							},
							{
								"disk_size_gb": 40
							}
						]
					}
				]
			}
		}
		`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"first"},
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
								{
									GbSize: to.IntPtr(20),
								},
							},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BooPtr(true),
							SubnetNames: []string{"second"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{
								{
									GbSize: to.IntPtr(30),
								},
								{
									GbSize: to.IntPtr(40),
								},
							},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

func TestConfig_Load_VmGroup_VmImage(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "missing vm_groups.vm_image parameters",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_image": {},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image parameters",
			json: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"rsa_pub_path": "some-file-name",  
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
					"vm_image": {
						"publisher": "",
						"offer": "",
						"sku": "",
						"version": ""
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

func configLoadTestingBody(t *testing.T, json []byte, want *Config, wantErr error) {
	got := &Config{}
	err := got.Unmarshal(json)

	if wantErr != nil {
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				t.Fatal(err)
			}
			errs := err.(validator.ValidationErrors)
			if len(errs) != len(wantErr.(test.TestValidationErrors)) {
				t.Fatalf("incorrect length of found errors. Got: \n%s\nExpected: \n%s", errs.Error(), wantErr.Error())
			}
			for _, e := range errs {
				found := false
				for _, we := range wantErr.(test.TestValidationErrors) {
					if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), wantErr.Error())
				}
			}
		} else {
			t.Errorf("No errors got. All expected errors: \n%s", wantErr.Error())
		}
	} else {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
		}
		if err != nil {
			t.Errorf("Unmarshal() unexpected error occured: %v", err)
		}
	}
}

func TestParams_ExtractEmptySubnets(t *testing.T) {
	tests := []struct {
		name   string
		params *Params
		want   []Subnet
	}{
		{
			name: "happy path",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet2"),
					AddressPrefixes: []string{"2.2.2.2/24"},
				},
			},
		},
		{
			name:   "nil params",
			params: nil,
			want:   nil,
		},
		{
			name: "nil subnets",
			params: &Params{
				Subnets: nil,
			},
			want: nil,
		},
		{
			name: "empty subnets",
			params: &Params{
				Subnets: []Subnet{},
			},
			want: nil,
		},
		{
			name: "nil vm_groups",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: nil,
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
			},
		},
		{
			name: "empty vm_groups",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
			},
		},
		{
			name: "no empty subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups no empty subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
					{
						SubnetNames: []string{"subnet2"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups reuse one subnet",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups some free subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
					{
						Name:            to.StrPtr("subnet3"),
						AddressPrefixes: []string{"3.3.3.3/24"},
					},
					{
						Name:            to.StrPtr("subnet4"),
						AddressPrefixes: []string{"4.4.4.4/24"},
					},
					{
						Name:            to.StrPtr("subnet5"),
						AddressPrefixes: []string{"5.5.5.5/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet2", "subnet5"},
					},
					{
						SubnetNames: []string{"subnet2", "subnet4"},
					},
				},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
				{
					Name:            to.StrPtr("subnet3"),
					AddressPrefixes: []string{"3.3.3.3/24"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.ExtractEmptySubnets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractEmptySubnets() = %v, want %v", got, tt.want)
			}
		})
	}
}
