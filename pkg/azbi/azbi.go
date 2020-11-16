package azbi

import (
	"encoding/json"
	"errors"
	maps "github.com/mitchellh/mapstructure"
)

type Params struct {
	VmsCount         *int     `json:"vms_count"`
	UsePublicIP      *bool    `json:"use_public_ip"`
	Location         *string  `json:"location"`
	Name             *string  `json:"name"`
	AddressSpace     []string `json:"address_space"`
	AddressPrefixes  []string `json:"address_prefixes"`
	RsaPublicKeyPath *string  `json:"rsa_pub_path"`
}

type Config struct {
	Kind    *string                `json:"kind"`
	Version *string                `json:"version"`
	Params  *Params                `json:"params"`
	Extra   map[string]interface{} `json:"extra"`
	Unused  []string               `json:"-"`
}

func (c Config) Save() (b []byte, err error) {
	return json.MarshalIndent(c, "", "\t")
}

func Load(b []byte) (c *Config, err error) {
	var input map[string]interface{}
	if err = json.Unmarshal(b, &input); err != nil {
		return
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &c,
	})
	if err != nil {
		return
	}
	err = d.Decode(input)
	if err != nil {
		return
	}
	c.Unused = md.Unused
	err = c.isValid()
	return
}

//TODO implement more interesting validation
func (c Config) isValid() error {
	if c.Version == nil {
		return errors.New("field 'Version' cannot be nil")
	}
	return nil
}
