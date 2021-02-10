package v0

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/google/go-cmp/cmp"
)

type TestValidationErrors []TestValidationError

func (e TestValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for _, te := range e {

		buff.WriteString(te.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

type TestValidationError struct {
	Key   string
	Field string
	Tag   string
}

func (e TestValidationError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", e.Key, e.Field, e.Tag)
}

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
			name: "unknown field in main structure",
			args: []byte(`{
	"kind": "azbi",
	"version": "v0.1.0",
	"extra_outer_field" : "extra_outer_value",
	"params": {
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
				Unused: []string{"extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in sub structure",
			args: []byte(`{
	"kind": "azbi",
	"version": "v0.1.0",
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
				Unused: []string{"params.extra_inner_field"},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in all possible places",
			args: []byte(`{
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
			args: []byte(`{}`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Kind",
					Field: "Kind",
					Tag:   "required",
				},
				TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "required",
				},
				TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "just kind field",
			args: []byte(`{
			"kind": "azbi"
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "required",
				},
				TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "just kind and version",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0"
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "just vm_groups in params",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "required",
				},
				TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "required",
				},
			},
		},
		{
			name: "minimal correct json",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("v0.1.0"),
				Params: &Params{
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing subnets list",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty subnets list",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [],
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing subnet name",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty subnet name",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty subnet address prefixes",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "main",
						"address_prefixes": []
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing subnet address prefixes",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "main"
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty subnet address prefixes element",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "main",
						"address_prefixes": [
							""
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes[0]",
					Field: "AddressPrefixes[0]",
					Tag:   "required",
				},
			},
		},
		{
			name: "multiple subnets configuration",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			name: "missing vm_groups parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "main",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups",
					Field: "VmGroups",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			name: "missing vm_groups name parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
				"subnets": [
					{
						"name": "main",
						"address_prefixes": [
							"10.0.1.0/24"
						]
					}
				],
				"vm_groups": [{
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].Name",
					Field: "Name",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing vm_groups vm_count parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "required",
				},
			},
		},
		{
			name: "negative vm_groups vm_count parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					"vm_count": -100,
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm_groups vm_size parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmSize",
					Field: "VmSize",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing vm_groups use_public_ip parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].UsePublicIP",
					Field: "UsePublicIP",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing vm_groups subnet_names parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames",
					Field: "SubnetNames",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups subnet_names list",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames",
					Field: "SubnetNames",
					Tag:   "min",
				},
			},
		},
		{
			name: "vm_groups subnet_names list empty value",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "SubnetNames[0]",
					Tag:   "required",
				},
			},
		},
		{
			name: "vm_groups subnet_names list value not existing in subnets",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "TODO Custom Validator",
					Field: "TODO",
					Tag:   "TODO",
				},
			},
		},
		{
			name: "missing vm_groups vm_image parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage",
					Field: "VmImage",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing vm_groups.vm_image.publisher parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image.publisher parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm_groups.vm_image.offer parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image.offer parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"offer": "",
						"sku": "18.04-LTS",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm_groups.vm_image.sku parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image.sku parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"sku": "",
						"version": "18.04.202006101"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm_groups.vm_image.version parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"sku": "18.04-LTS"
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image.version parameter",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
						"version": ""
					},
					"data_disks": []
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm_groups.data_disks list",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					}
				}]
			}
		}
		`),
			want: nil,
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks",
					Field: "DataDisks",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.data_disks list value",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "required",
				},
			},
		},
		{
			name: "zero vm_groups.data_disks list value",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
			},
		},
		{
			name: "negative vm_groups.data_disks list value",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
			},
		},
		{
			name: "multiple vm_groups configuration",
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			args: []byte(`{
			"kind": "azbi",
			"version": "v0.1.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
		{
			name: "major version mismatch",
			args: []byte(`{
			"kind": "azbi",
			"version": "100.0.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			wantErr: TestValidationErrors{
				TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "semver",
				},
			},
		},
		{
			name: "minor version mismatch",
			args: []byte(`{
			"kind": "azbi",
			"version": "0.100.0",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.100.0"),
				Params: &Params{
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			args: []byte(`{
			"kind": "azbi",
			"version": "0.0.100",
			"params": {
				"location": "northeurope",
				"name": "epiphany",
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
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.100"),
				Params: &Params{
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
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
			got := &Config{}
			err := got.Unmarshal(tt.args)

			if tt.wantErr != nil {
				if err != nil {
					if _, ok := err.(*validator.InvalidValidationError); ok {
						t.Fatal(err)
					}
					errs := err.(validator.ValidationErrors)
					if len(errs) != len(tt.wantErr.(TestValidationErrors)) {
						t.Fatalf("incorrect length of found errors. Got: \n%s\nExpected: \n%s", errs.Error(), tt.wantErr.Error())
					}
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(TestValidationErrors) {
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
