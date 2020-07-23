package packets

import "github.com/mitchellh/mapstructure"
import "fenix/errors"
import packets "fenix/packets"

// Version packet
type Version struct {
	id interface{}
}

// FromJSON instanciates a Version packet from a JSON map
func (p Version) FromJSON(data map[string]interface{}) (packets.Packet, error) {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnused: true, Result: p})
	err := decoder.Decode(data["d"])

	if err != nil {
		return nil, errors.DataError{}
	}

	p.id = data["id"]

	return &p, nil
}

// SetID sets the ID for this packet.
func (p *Version) SetID(id interface{}) {
	p.id = id
}

// CopyToNewPacket copies all user essential data to a new packet
func (p *Version) CopyToNewPacket(new packets.Packet) {
	new.SetID(p.id)
}
