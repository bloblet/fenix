package packets

// A Packet is a base interface for all supported packets, that works with all supported versions.
type Packet interface {
	CopyToNewPacket(Packet)
	FromJSON(map[string]interface{}) (Packet, error)
	SetID(id interface{})
}