package packets

import "github.com/mitchellh/mapstructure"
import "fenix/server/errors"
import packets "fenix/server/packets"

type Version struct {
	id interface{}
}

func (p Version) FromJSON(data map[string]interface{}) (packets.Packet, error) {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: p})
	err := decoder.Decode(data["d"])

	if err != nil {
		return nil, errors.DataError{}
	}

	p.id = data["id"]

	return &p, nil
}
func (p *Version) SetID(id interface{}) {
	p.id = id
}

// CopyToNewPacket copies all user essential data to a new packet
func (p *Version) CopyToNewPacket(new packets.Packet) {
	new.SetID(p.id)
}
