package azbi

import (
	"github.com/mkyc/go-stucts-versioning-tests/pkg/to"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Config
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				b: []byte(`{
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
			},
			wantC: &Config{
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
			wantErr: false,
		},
		{
			name: "unknown field in main structure",
			args: args{
				b: []byte(`{
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
			},
			wantC: &Config{
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
			wantErr: false,
		},
		{
			name: "unknown field in sub structure",
			args: args{
				b: []byte(`{
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
			},
			wantC: &Config{
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
			wantErr: false,
		},
		{
			name: "unknown field in extra fields",
			args: args{
				b: []byte(`{
	"kind": "azbi",
	"version": "0.0.2",
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
	}, 
	"extra": {
		"field-one-key": "field-one-value"
	}
}
`),
			},
			wantC: &Config{
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
				Unused: []string{},
				Extra:  map[string]interface{}{"field-one-key": "field-one-value"},
			},
			wantErr: false,
		},
		{
			name: "unknown fields in all possible places",
			args: args{
				b: []byte(`{
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
	}, 
	"extra": {
		"field-one-key": "field-one-value"
	}
}
`),
			},
			wantC: &Config{
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
				Extra:  map[string]interface{}{"field-one-key": "field-one-value"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := Load(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("Load() gotC = \n\n%#v\n\n, want = \n\n%#v\n\n", gotC, tt.wantC)
			}
		})
	}
}
