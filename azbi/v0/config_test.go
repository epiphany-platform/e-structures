package v0

import (
	"fmt"
	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"reflect"
	"testing"
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
					"disk_size_gb": 10, 
					"storage_type": "Premium_LRS"
				}
			]
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
									GbSize:      to.IntPtr(10),
									StorageType: to.StrPtr("Premium_LRS"),
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
					"disk_size_gb": 10, 
					"storage_type": "Premium_LRS"
				}
			]
		}],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
									GbSize:      to.IntPtr(10),
									StorageType: to.StrPtr("Premium_LRS"),
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
					Key:   "Config.Meta",
					Field: "Meta",
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
			name: "meta missing",
			json: []byte(`{
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
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
			},
		},
		{
			name: "major version mismatch",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v100.0.0",
		"module_version": "v0.0.1"
	},
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
					Key:   "Config.Meta.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v0.100.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.100.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.0.100",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.0.100"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
		{
			name: "empty params.name",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "",
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "required",
				},
			},
		},
		{
			name: "incorrect vm_groups.data_disks list value",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
					"disk_size_gb": 0, 
					"storage_type": "incorrect"
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
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "eq=Standard_LRS|eq=Premium_LRS|eq=StandardSSD_LRS|eq=UltraSSD_LRS",
				},
			},
		},
		{
			name: "another incorrect vm_groups.data_disks list value",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
					"disk_size_gb": -1, 
					"storage_type": ""
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
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "eq=Standard_LRS|eq=Premium_LRS|eq=StandardSSD_LRS|eq=UltraSSD_LRS",
				},
			},
		},
		{
			name: "multiple vm_groups configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
						"disk_size_gb": 10, 
						"storage_type": "Premium_LRS"
					},
					{
						"disk_size_gb": 20, 
						"storage_type": "Standard_LRS"
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
						"disk_size_gb": 30, 
						"storage_type": "StandardSSD_LRS"
					},
					{
						"disk_size_gb": 40, 
						"storage_type": "UltraSSD_LRS"
					}
				]
			}
		]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbi"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
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
									GbSize:      to.IntPtr(10),
									StorageType: to.StrPtr("Premium_LRS"),
								},
								{
									GbSize:      to.IntPtr(20),
									StorageType: to.StrPtr("Standard_LRS"),
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
									GbSize:      to.IntPtr(30),
									StorageType: to.StrPtr("StandardSSD_LRS"),
								},
								{
									GbSize:      to.IntPtr(40),
									StorageType: to.StrPtr("UltraSSD_LRS"),
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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
	"meta": {
		"kind": "azbi",
		"version": "v0.1.0",
		"module_version": "v0.0.1"
	},
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

func configLoadTestingBody(t *testing.T, json []byte, want *Config, wantErr error) {
	a := assert.New(t)
	r := require.New(t)
	p, err := createTempDocumentFile("azbi-load", json)
	r.NoError(err)
	got := &Config{}
	err = got.Load(p)

	if wantErr != nil {
		r.Error(err)
		_, ok := err.(*validator.InvalidValidationError)
		r.Equal(false, ok)
		_, ok = err.(validator.ValidationErrors)
		r.Equal(true, ok)
		errs := err.(validator.ValidationErrors)
		a.Equal(len(wantErr.(test.TestValidationErrors)), len(errs))

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
		a.NoError(err)
		a.Equal(want, got)
	}
}

func createTempDocumentFile(name string, document []byte) (string, error) {
	p, err := ioutil.TempDir("", fmt.Sprintf("e-structures-%s-*", name))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path.Join(p, "file.json"), document, 0644)
	return path.Join(p, "file.json"), err
}
