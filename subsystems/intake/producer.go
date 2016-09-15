package intake

import (
	"bufio"
	"net"
	"thrust/config"
	"thrust/logging"
)

type messageStruct struct {
	AckChannel chan bool
	Payload    []byte
}

func serve(connection net.Conn, turbineChannel chan<- messageStruct) {
	logging.NewProducer(connection.RemoteAddr())
	defer logging.LostProducer(connection.RemoteAddr())

	ackChannel := make(chan bool)
	reader := bufio.NewReader(connection)

	for {
		payload, err := reader.ReadSlice('\n')
		if err != nil {
			return
		}

		logging.WatchCapacity("dumper", len(turbineChannel), config.Config.Intake.CompressorBlades)

		turbineChannel <- messageStruct{AckChannel: ackChannel, Payload: payload}

		<-ackChannel // recieve acknowledgement

		connection.Write([]byte{'y'})
	}
}