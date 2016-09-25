package intake

import (
	// "fmt"
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"thrust/common"
	"thrust/logging"
	"thrust/subsystems/exhaust"
)

const headerSize = 12
const headerFieldSize = 4
const headerFieldTopic = 8

func parseHeader(reader *bufio.Reader) (int64, int32) {
	buffer := make([]byte, headerSize)
	bytesRead, _ := io.ReadFull(reader, buffer)
	if bytesRead != headerSize {
		return 0, 0
	}
	topic := int64(binary.LittleEndian.Uint64(buffer[:headerFieldTopic]))
	size := int32(binary.LittleEndian.Uint32(buffer[headerFieldTopic:]))
	return topic, size
}

func suck(connection net.Conn) {
	logging.NewProducer(connection.RemoteAddr())
	defer logging.LostProducer(connection.RemoteAddr())

	ackChannel := make(chan bool, 1)
	reader := bufio.NewReader(connection)

	for {
		topic, size := parseHeader(reader)
		if size == 0 {
			return
		}

		buffer := make([]byte, size)
		_, err := io.ReadFull(reader, buffer)
		common.FaceIt(err)

		message := common.MessageStruct{AckChannel: ackChannel, Payload: buffer, Topic: topic}
		CompressorChannel <- message
		<-ackChannel // receive acknowledgement, then move forward
		select {
		case exhaust.CombustorChannel <- message:
		}
		connection.Write([]byte{'y'})
	}
}
