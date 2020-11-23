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
	"version": "0.0.1",
	"params": {
		"vms_count": 3,
		"use_public_ip": true,
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"address_prefixes": [
			"10.0.1.0/24"
		],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.1"),
				Params: &Params{
					VmsCount:         to.IntPtr(3),
					UsePublicIP:      to.BooPtr(true),
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AddressPrefixes:  []string{"10.0.1.0/24"},
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown field in main structure",
			args: []byte(`{
	"kind": "azbi",
	"version": "0.0.2",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"vms_count": 3,
		"use_public_ip": true,
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"address_prefixes": [
			"10.0.1.0/24"
		],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.2"),
				Params: &Params{
					VmsCount:         to.IntPtr(3),
					UsePublicIP:      to.BooPtr(true),
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AddressPrefixes:  []string{"10.0.1.0/24"},
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
	"version": "0.0.2",
	"params": {
		"vms_count": 3,
		"extra_inner_field" : "extra_inner_value", 
		"use_public_ip": true,
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"address_prefixes": [
			"10.0.1.0/24"
		],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.2"),
				Params: &Params{
					VmsCount:         to.IntPtr(3),
					UsePublicIP:      to.BooPtr(true),
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AddressPrefixes:  []string{"10.0.1.0/24"},
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
	"version": "0.0.2",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"vms_count": 3,
		"extra_inner_field" : "extra_inner_value", 
		"use_public_ip": true,
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"address_prefixes": [
			"10.0.1.0/24"
		],
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.2"),
				Params: &Params{
					VmsCount:         to.IntPtr(3),
					UsePublicIP:      to.BooPtr(true),
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AddressPrefixes:  []string{"10.0.1.0/24"},
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
	"version": "0.0.1"
}
`),
			want:    nil,
			wantErr: ParamsMissingValidationError,
		},
		{
			name: "just vms_count in params",
			args: []byte(`{
	"kind": "azbi",
	"version": "0.0.1",
	"params": {
		"vms_count": 3
	}
}
`),
			want:    nil,
			wantErr: MinimalParamsValidationError,
		},
		{
			name: "minimal correct json",
			args: []byte(`{
	"kind": "azbi",
	"version": "0.0.1",
	"params": {
		"vms_count": 3,
		"location": "northeurope",
		"name": "epiphany"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.1"),
				Params: &Params{
					VmsCount: to.IntPtr(3),
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "major version mismatch",
			args: []byte(`{
	"kind": "azbi",
	"version": "1.0.0",
	"params": {
		"vms_count": 3,
		"location": "northeurope",
		"name": "epiphany"
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
	"version": "0.1.0",
	"params": {
		"vms_count": 3,
		"location": "northeurope",
		"name": "epiphany"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.1.0"),
				Params: &Params{
					VmsCount: to.IntPtr(3),
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			args: []byte(`{
	"kind": "azbi",
	"version": "0.0.2",
	"params": {
		"vms_count": 3,
		"location": "northeurope",
		"name": "epiphany"
	}
}
`),
			want: &Config{
				Kind:    to.StrPtr("azbi"),
				Version: to.StrPtr("0.0.2"),
				Params: &Params{
					VmsCount: to.IntPtr(3),
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Config{}
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
