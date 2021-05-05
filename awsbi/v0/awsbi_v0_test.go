package v0

import (
	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
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
	"kind": "awsbi",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 1,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [
				{
					"name": "first_private_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.1.0/24"
				}
			],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [
			{
				"name": "default_sg",
				"rules": {
					"ingress": [
						{
							"protocol": "-1",
							"from_port": 0,
							"to_port": 0,
							"cidr_blocks": [
								"10.1.0.0/20"
							]
						},
						{
							"protocol": "tcp",
							"from_port": 22,
							"to_port": 22,
							"cidr_blocks": [
								"0.0.0.0/0"
							]
						}
					],
					"egress": [
						{
							"protocol": "-1",
							"from_port": 0,
							"to_port": 0,
							"cidr_blocks": [
								"0.0.0.0/0"
							]
						}
					]
				}
			}
		],
		"vm_groups": [
			{
				"name": "vm-group0",
				"vm_count": 1,
				"vm_size": "t3.medium",
				"use_public_ip": false,
				"subnet_names": [
					"first_private_subnet"
				],
				"sg_names": [
					"default_sg"
				],
				"vm_image": {
					"ami": "RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2",
					"owner": "309956199498"
				},
				"root_volume_size": 30,
				"data_disks": [
					{
						"device_name": "/dev/sdf",
						"disk_size_gb": 16,
						"type": "gp2"
					}
				]
			}
		]
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					Name:                  to.StrPtr("epiphany"),
					Region:                to.StrPtr("eu-central-1"),
					NatGatewayCount:       to.IntPtr(1),
					VirtualPrivateGateway: to.BooPtr(false),
					RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
					VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
					Subnets: &Subnets{
						Private: []Subnet{
							{
								Name:             to.StrPtr("first_private_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.1.0/24"),
							},
						},
						Public: []Subnet{
							{
								Name:             to.StrPtr("first_public_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
							},
						},
					},
					SecurityGroups: []SecurityGroup{
						{
							Name: to.StrPtr("default_sg"),
							Rules: &Rules{
								Ingress: []SecurityRule{
									{
										Protocol:   to.StrPtr("-1"),
										FromPort:   to.IntPtr(0),
										ToPort:     to.IntPtr(0),
										CidrBlocks: []string{"10.1.0.0/20"},
									},
									{
										Protocol:   to.StrPtr("tcp"),
										FromPort:   to.IntPtr(22),
										ToPort:     to.IntPtr(22),
										CidrBlocks: []string{"0.0.0.0/0"},
									},
								},
								Egress: []SecurityRule{
									{
										Protocol:   to.StrPtr("-1"),
										FromPort:   to.IntPtr(0),
										ToPort:     to.IntPtr(0),
										CidrBlocks: []string{"0.0.0.0/0"},
									},
								},
							},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:               to.StrPtr("vm-group0"),
							VmCount:            to.IntPtr(1),
							VmSize:             to.StrPtr("t3.medium"),
							UsePublicIp:        to.BooPtr(false),
							SubnetNames:        []string{"first_private_subnet"},
							SecurityGroupNames: []string{"default_sg"},
							VmImage: &VmImage{
								AMI:   to.StrPtr("RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2"),
								Owner: to.StrPtr("309956199498"),
							},
							RootVolumeGbSize: to.IntPtr(30),
							DataDisks: []DataDisk{
								{
									DeviceName: to.StrPtr("/dev/sdf"),
									GbSize:     to.IntPtr(16),
									Type:       to.StrPtr("gp2"),
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
	"kind": "awsbi",
	"version": "v0.0.1",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"extra_inner_field" : "extra_inner_value",
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 1,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [
				{
					"name": "first_private_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.1.0/24"
				}
			],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [
			{
				"name": "default_sg",
				"rules": {
					"ingress": [
						{
							"protocol": "-1",
							"from_port": 0,
							"to_port": 0,
							"cidr_blocks": [
								"10.1.0.0/20"
							]
						},
						{
							"protocol": "tcp",
							"from_port": 22,
							"to_port": 22,
							"cidr_blocks": [
								"0.0.0.0/0"
							]
						}
					],
					"egress": [
						{
							"protocol": "-1",
							"from_port": 0,
							"to_port": 0,
							"cidr_blocks": [
								"0.0.0.0/0"
							]
						}
					]
				}
			}
		],
		"vm_groups": [
			{
				"name": "vm-group0",
				"vm_count": 1,
				"vm_size": "t3.medium",
				"use_public_ip": false,
				"subnet_names": [
					"first_private_subnet"
				],
				"sg_names": [
					"default_sg"
				],
				"vm_image": {
					"ami": "RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2",
					"owner": "309956199498"
				},
				"root_volume_size": 30,
				"data_disks": [
					{
						"device_name": "/dev/sdf",
						"disk_size_gb": 16,
						"type": "gp2"
					}
				]
			}
		]
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					Name:                  to.StrPtr("epiphany"),
					Region:                to.StrPtr("eu-central-1"),
					NatGatewayCount:       to.IntPtr(1),
					VirtualPrivateGateway: to.BooPtr(false),
					RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
					VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
					Subnets: &Subnets{
						Private: []Subnet{
							{
								Name:             to.StrPtr("first_private_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.1.0/24"),
							},
						},
						Public: []Subnet{
							{
								Name:             to.StrPtr("first_public_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
							},
						},
					},
					SecurityGroups: []SecurityGroup{
						{
							Name: to.StrPtr("default_sg"),
							Rules: &Rules{
								Ingress: []SecurityRule{
									{
										Protocol:   to.StrPtr("-1"),
										FromPort:   to.IntPtr(0),
										ToPort:     to.IntPtr(0),
										CidrBlocks: []string{"10.1.0.0/20"},
									},
									{
										Protocol:   to.StrPtr("tcp"),
										FromPort:   to.IntPtr(22),
										ToPort:     to.IntPtr(22),
										CidrBlocks: []string{"0.0.0.0/0"},
									},
								},
								Egress: []SecurityRule{
									{
										Protocol:   to.StrPtr("-1"),
										FromPort:   to.IntPtr(0),
										ToPort:     to.IntPtr(0),
										CidrBlocks: []string{"0.0.0.0/0"},
									},
								},
							},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:               to.StrPtr("vm-group0"),
							VmCount:            to.IntPtr(1),
							VmSize:             to.StrPtr("t3.medium"),
							UsePublicIp:        to.BooPtr(false),
							SubnetNames:        []string{"first_private_subnet"},
							SecurityGroupNames: []string{"default_sg"},
							VmImage: &VmImage{
								AMI:   to.StrPtr("RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2"),
								Owner: to.StrPtr("309956199498"),
							},
							RootVolumeGbSize: to.IntPtr(30),
							DataDisks: []DataDisk{
								{
									DeviceName: to.StrPtr("/dev/sdf"),
									GbSize:     to.IntPtr(16),
									Type:       to.StrPtr("gp2"),
								},
							},
						},
					},
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
	"kind": "awsbi",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 0,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [],
		"vm_groups": []
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					Name:                  to.StrPtr("epiphany"),
					Region:                to.StrPtr("eu-central-1"),
					NatGatewayCount:       to.IntPtr(0),
					VirtualPrivateGateway: to.BooPtr(false),
					RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
					VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
					Subnets: &Subnets{
						Private: []Subnet{},
						Public: []Subnet{
							{
								Name:             to.StrPtr("first_public_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
							},
						},
					},
					SecurityGroups: []SecurityGroup{},
					VmGroups:       []VmGroup{},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "major version mismatch",
			json: []byte(`{
	"kind": "awsbi",
	"version": "v100.0.1",
	"params": {
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 0,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [],
		"vm_groups": []
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
	"kind": "awsbi",
	"version": "v0.100.1",
	"params": {
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 0,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [],
		"vm_groups": []
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr("v0.100.1"),
				Params: &Params{
					Name:                  to.StrPtr("epiphany"),
					Region:                to.StrPtr("eu-central-1"),
					NatGatewayCount:       to.IntPtr(0),
					VirtualPrivateGateway: to.BooPtr(false),
					RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
					VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
					Subnets: &Subnets{
						Private: []Subnet{},
						Public: []Subnet{
							{
								Name:             to.StrPtr("first_public_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
							},
						},
					},
					SecurityGroups: []SecurityGroup{},
					VmGroups:       []VmGroup{},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			json: []byte(`{
	"kind": "awsbi",
	"version": "v0.0.100",
	"params": {
		"name": "epiphany",
		"region": "eu-central-1",
		"nat_gateway_count": 0,
		"virtual_private_gateway": false,
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vpc_address_space": "10.1.0.0/20",
		"subnets": {
			"private": [],
			"public": [
				{
					"name": "first_public_subnet",
					"availability_zone": "any",
					"address_prefixes": "10.1.2.0/24"
				}
			]
		},
		"security_groups": [],
		"vm_groups": []
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr("v0.0.100"),
				Params: &Params{
					Name:                  to.StrPtr("epiphany"),
					Region:                to.StrPtr("eu-central-1"),
					NatGatewayCount:       to.IntPtr(0),
					VirtualPrivateGateway: to.BooPtr(false),
					RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
					VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
					Subnets: &Subnets{
						Private: []Subnet{},
						Public: []Subnet{
							{
								Name:             to.StrPtr("first_public_subnet"),
								AvailabilityZone: to.StrPtr("any"),
								AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
							},
						},
					},
					SecurityGroups: []SecurityGroup{},
					VmGroups:       []VmGroup{},
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
