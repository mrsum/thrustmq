package common

import (
	"encoding/binary"
	"io"
)

type IndexRecord struct {
	Offset     uint64
	Length     uint64
	Topic      uint64
	Connection uint64
	Ack        byte
}

var IndexSize int64 = 33

func (self IndexRecord) Serialize() []byte {
	buffer := make([]byte, IndexSize)
	binary.LittleEndian.PutUint64(buffer[0:8], self.Offset)
	binary.LittleEndian.PutUint64(buffer[8:16], self.Length)
	binary.LittleEndian.PutUint64(buffer[16:24], self.Topic)
	binary.LittleEndian.PutUint64(buffer[24:32], self.Connection)
	buffer[32] = self.Ack
	return buffer
}

func (self *IndexRecord) Deserialize(reader io.Reader) bool {
	buffer := make([]byte, IndexSize)
	bytesRead, _ := io.ReadFull(reader, buffer)
	if int64(bytesRead) != IndexSize {
		return false
	}
	self.Offset = binary.LittleEndian.Uint64(buffer[0:8])
	self.Length = binary.LittleEndian.Uint64(buffer[8:16])
	self.Topic = binary.LittleEndian.Uint64(buffer[16:24])
	self.Connection = binary.LittleEndian.Uint64(buffer[24:32])
	self.Ack = buffer[32]
	return true
}
