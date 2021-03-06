package exhaust

import (
	"bufio"
	"github.com/rambler-digital-solutions/thrustmq/common"
	"github.com/rambler-digital-solutions/thrustmq/config"
	"github.com/rambler-digital-solutions/thrustmq/logging"
	"github.com/rambler-digital-solutions/thrustmq/subsystems/oplog"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"
)

func registerConnect(connection net.Conn) *common.ConnectionStruct {
	common.State.ConnectionId++

	connectionStruct := &common.ConnectionStruct{}
	connectionStruct.Connection = connection
	connectionStruct.Id = common.State.ConnectionId
	connectionStruct.Reader = bufio.NewReaderSize(connection, config.Base.NetworkBuffer)
	connectionStruct.Writer = bufio.NewWriterSize(connection, config.Base.NetworkBuffer)
	connectionStruct.Channel = make(common.RecordPipe, config.Exhaust.NozzleBuffer)

	ConnectionsMap[connectionStruct.Id] = connectionStruct
	logging.NewConsumer(connectionStruct, len(ConnectionsMap))

	connectionStruct.DeserializeHeader()
	logging.NewConsumerHeader(connectionStruct)

	return connectionStruct
}

func registerDisconnect(connectionStruct *common.ConnectionStruct) {
	delete(ConnectionsMap, connectionStruct.Id)
	logging.LostConsumer(connectionStruct.Connection.RemoteAddr(), len(ConnectionsMap))
}

func sendBatch(client *common.ConnectionStruct, batch []*common.Record) {
	client.SendActualBatchSize(len(batch))
	for i := 0; i < len(batch); i++ {
		record := <-client.Channel
		record.Sent = common.TimestampUint64()
		err := client.SendMessage(record)
		if err != nil {
			log.Print(err)
			return
		}
		batch[i] = record
		oplog.ExhaustThroughput++
	}
	client.Writer.Flush()
}

func recieveAcks(client *common.ConnectionStruct, batch []*common.Record) {
	acks, _ := client.GetAcks(len(batch))
	for i := 0; i < len(batch); i++ {
		if acks[i] == 1 {
			batch[i].Delivered = common.TimestampUint64()
		} else {
			log.Print("returning record to combustor")
			log.Print(acks[i])
		}
	}
}

func blow(connection net.Conn) {
	client := registerConnect(connection)
	defer registerDisconnect(client)

	for {
		batchSize := common.Min(int(client.BatchSize), len(client.Channel))

		if batchSize > 0 {
			batch := make([]*common.Record, batchSize)
			sendBatch(client, batch)
			recieveAcks(client, batch)
		} else {
			logging.Debug("Trying to ping client", strconv.FormatInt(int64(client.Id), 4), "...")
			time.Sleep(time.Duration(config.Exhaust.HeartbeatRate) * time.Nanosecond)
			runtime.Gosched()
			if !client.Ping() {
				return
			}
		}
	}
}
