package dashboard

import (
	"fmt"
	"thrust/common"
	"thrust/config"
	"thrust/subsystems/exhaust"
	"thrust/subsystems/intake"
	"thrust/subsystems/oplog"
	"time"
)

func Init() {
	for {
		time.Sleep(time.Second)
		if config.Config.Debug {
			fmt.Printf("\r %6d ->msg/sec %6d msg/sec->", oplog.IntakeThroughput, oplog.ExhaustThroughput)
			fmt.Printf(" | %4d->compressor %4d->combustor %4d->turbine", len(intake.CompressorChannel), len(exhaust.CombustorChannel), len(exhaust.TurbineChannel))
			fmt.Printf(" | r %d conn_id: %d", oplog.Requeued, exhaust.State.ConnectionId)
			fmt.Printf(" | h %d t %d span: %d capacity: %.2f", exhaust.State.Head, exhaust.State.Tail, (exhaust.State.Head-exhaust.State.Tail)/uint64(common.IndexSize), exhaust.State.Capacity)
			for _, connectionStruct := range exhaust.ConnectionsMap {
				fmt.Printf(" [%4d]", len(connectionStruct.Channel))
			}
			fmt.Printf(" | errata %d", oplog.ExhaustTotal-oplog.IntakeTotal)
			fmt.Printf("               ")
			oplog.IntakeTotal += oplog.IntakeThroughput
			oplog.ExhaustTotal += oplog.ExhaustThroughput
			oplog.IntakeThroughput = 0
			oplog.ExhaustThroughput = 0
		}
	}
}
