package row

import (
	"encoding/binary"
)

type Row struct {
	ID    int
	Value string
}

// Serialize 将 Row 序列化为字节流
func (r *Row) Serialize() []byte {
	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, uint32(r.ID))

	valueBytes := []byte(r.Value)
	valueLength := len(valueBytes)

	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(valueLength))

	return append(append(idBytes, lengthBytes...), valueBytes...)
}
