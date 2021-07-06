package v0

import (
	"github.com/epiphany-platform/e-structures/globals"
)

type State struct {
	Status globals.Status `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *Config        `json:"config" validate:"omitempty"`
	Output *Output        `json:"output" validate:"omitempty"`
}

type Output struct {
	RgName   *string         `json:"rg_name"`
	VnetName *string         `json:"vnet_name"`
	VmGroups []OutputVmGroup `json:"vm_groups"`
}

type OutputVmGroup struct {
	Name *string    `json:"vm_group_name"`
	Vms  []OutputVm `json:"vms"`
}

type OutputVm struct {
	Name       *string          `json:"vm_name"`
	PrivateIps []string         `json:"private_ips"`
	PublicIp   *string          `json:"public_ip"`
	DataDisks  []OutputDataDisk `json:"data_disks"`
}

type OutputDataDisk struct {
	Size *int `json:"size"`
	Lun  *int `json:"lun"`
}
