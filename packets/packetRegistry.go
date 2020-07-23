package packets
import e "fenix/errors"

// SupportedVersions is a map of current packet versions
var SupportedVersions = map[interface{}]bool {"6.0": true}

// PacketRegistry is a map of current packet versions
var PacketRegistry = map[interface{}]Packet {}

// ValidatePacket validates a recieved and unparsed packet
func ValidatePacket(packet map[string]interface{}) (Packet, error) {
	if _, validType := PacketRegistry[packet["t"]]; !validType {
		return nil, e.TypeError{}
	}
	if _, ok := packet["d"]; !ok {
		return nil, e.DataError{}
	}
	
	if _, ok := packet["id"]; !ok {
		return nil, e.DataError{}
	}
	
	return PacketRegistry[packet["t"]].FromJSON(packet)
}

