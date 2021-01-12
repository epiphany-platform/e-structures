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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
			name:    "empty json",
			args:    []byte(`{}`),
			want:    nil,
			wantErr: KindMissingValidationError,
		},
		{
			name: "just kind field",
			args: []byte(`{
	"kind": "azbi"
}
`),
			want:    nil,
			wantErr: VersionMissingValidationError,
		},
		{
			name: "just kind and version",
			args: []byte(`{
	"kind": "azbi",
	"version": "v0.1.0"
}
`),
			want:    nil,
			wantErr: ParamsMissingValidationError,
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'name' parameter missing"},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'subnets' list parameter missing or is 0 length"},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'subnets' list parameter missing or is 0 length"},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"one of subnets is missing 'name' field or name is empty"},
		},
		{
			name: "0 length subnet name",
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"one of subnets is missing 'name' field or name is empty"},
		},
		{
			name: "0 length subnet address prefixes",
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'address_prefixes' list parameter in one of subnets missing or is 0 length"},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: &MinimalParamsValidationError{"'address_prefixes' list parameter in one of subnets missing or is 0 length"},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
		}]
	}
}
`),
			want:    nil,
			wantErr: MajorVersionMismatchError,
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
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
			"image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}
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
							Image: &Image{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
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
