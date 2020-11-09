package azbi

const (
	moduleShortName = "azbi"
)

type Params struct {
	VmsCount         int      `json:"vms_count"`
	UsePublicIP      bool     `json:"use_public_ip"`
	Location         string   `json:"location"`
	Name             string   `json:"name"`
	AddressSpace     []string `json:"address_space"`
	AddressPrefixes  []string `json:"address_prefixes"`
	RsaPublicKeyPath string   `json:"rsa_pub_path"`
}

type Config struct {
	Kind   string `json:"kind"`
	Params Params `json:"params"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Kind: moduleShortName,
		Params: Params{
			VmsCount:         3,
			UsePublicIP:      true,
			Location:         "northeurope",
			Name:             "epiphany",
			AddressSpace:     []string{"10.0.0.0/16"},
			AddressPrefixes:  []string{"10.0.1.0/24"},
			RsaPublicKeyPath: "/shared/vms_rsa.pub",
		},
	}
}
