package common

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type ConnectionStruct struct {
	Connection net.Conn
	Bucket     uint64
	ClientId   uint64
	BatchSize  uint32
	Id         uint64
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	Channel    MessageChannel
}

type ConnectionsMap map[uint64]ConnectionStruct
type BucketsMap map[uint64]([]uint64)

var ConnectionHeaderSize = 20

func (self *ConnectionStruct) DeserializeHeader() {
	buffer := make([]byte, ConnectionHeaderSize)
	_, err := io.ReadFull(self.Reader, buffer)
	FaceIt(err)

	self.ClientId = binary.LittleEndian.Uint64(buffer[0:8])
	self.Bucket = binary.LittleEndian.Uint64(buffer[8:16])
	self.BatchSize = binary.LittleEndian.Uint32(buffer[16:20])
}

func (self *ConnectionStruct) NextBatchSize() int {
	return min(int(self.BatchSize), len(self.Channel))
}

func (self *ConnectionStruct) SendActualBatchSize(batchSize int) {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(batchSize))
	self.Writer.Write(buffer)
}

func (self *ConnectionStruct) SendMessage(message MessageStruct) error {
	bytes := message.Serialize()
	_, err := self.Writer.Write(bytes)
	return err
}

func (self *ConnectionStruct) GetAcks(batchSize int) ([]byte, error) {
	buffer := make([]byte, batchSize)
	_, err := io.ReadFull(self.Reader, buffer)
	return buffer, err
}

func (self *ConnectionStruct) Ping() {
	self.SendActualBatchSize(1)
	message := MessageStruct{}
	self.SendMessage(message)
	self.Writer.Flush()
	self.GetAcks(1)
}
