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
  "kind": "hi",
  "version": "v0.0.1",
  "params": {
    "vm_groups": [
      {
        "name": "vm-group0",
        "admin_user": "operations",
        "hosts": [
          {
            "name": "epiphany-vm-group0-1",
            "ip": "11.22.33.44"
          }
        ],
        "mount_point": [
          {
            "lun": 10,
            "path": "/data/lun10"
          }
        ]
      }
    ],
    "rsa_private_path": "/shared/vms_rsa"
  }
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					VmGroups: []VmGroup{
						{
							Name:      to.StrPtr("vm-group0"),
							AdminUser: to.StrPtr("operations"),
							Hosts: []Host{
								{
									Name: to.StrPtr("epiphany-vm-group0-1"),
									Ip:   to.StrPtr("11.22.33.44"),
								},
							},
							MountPoints: []MountPoint{
								{
									Lun:  to.IntPtr(10),
									Path: to.StrPtr("/data/lun10"),
								},
							},
						},
					},
					RsaPrivateKeyPath: to.StrPtr("/shared/vms_rsa"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
  "kind": "hi",
  "version": "v0.0.1",
  "extra_outer_field" : "extra_outer_value",
  "params": {
    "extra_inner_field" : "extra_inner_value",
    "vm_groups": [
      {
        "name": "vm-group0",
        "admin_user": "operations",
        "hosts": [
          {
            "name": "epiphany-vm-group0-1",
            "ip": "11.22.33.44"
          }
        ],
        "mount_point": [
          {
            "lun": 10,
            "path": "/data/lun10"
          }
        ]
      }
    ],
    "rsa_private_path": "/shared/vms_rsa"
  }
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					VmGroups: []VmGroup{
						{
							Name:      to.StrPtr("vm-group0"),
							AdminUser: to.StrPtr("operations"),
							Hosts: []Host{
								{
									Name: to.StrPtr("epiphany-vm-group0-1"),
									Ip:   to.StrPtr("11.22.33.44"),
								},
							},
							MountPoints: []MountPoint{
								{
									Lun:  to.IntPtr(10),
									Path: to.StrPtr("/data/lun10"),
								},
							},
						},
					},
					RsaPrivateKeyPath: to.StrPtr("/shared/vms_rsa"),
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
  "kind": "hi",
  "version": "v0.0.1",
  "params": {
    "vm_groups": [],
    "rsa_private_path": "/shared/vms_rsa"
  }
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr(version),
				Params: &Params{
					VmGroups:          []VmGroup{},
					RsaPrivateKeyPath: to.StrPtr("/shared/vms_rsa"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "major version mismatch",
			json: []byte(`{
  "kind": "hi",
  "version": "v100.0.1",
  "params": {
    "vm_groups": [],
    "rsa_private_path": "/shared/vms_rsa"
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
  "kind": "hi",
  "version": "v0.100.1",
  "params": {
    "vm_groups": [],
    "rsa_private_path": "/shared/vms_rsa"
  }
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr("v0.100.1"),
				Params: &Params{
					VmGroups:          []VmGroup{},
					RsaPrivateKeyPath: to.StrPtr("/shared/vms_rsa"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			json: []byte(`{
  "kind": "hi",
  "version": "v0.0.100",
  "params": {
    "vm_groups": [],
    "rsa_private_path": "/shared/vms_rsa"
  }
}
`),
			want: &Config{
				Kind:    to.StrPtr(kind),
				Version: to.StrPtr("v0.0.100"),
				Params: &Params{
					VmGroups:          []VmGroup{},
					RsaPrivateKeyPath: to.StrPtr("/shared/vms_rsa"),
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
