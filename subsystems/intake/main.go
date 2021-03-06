package intake

import (
	"fmt"
	"github.com/rambler-digital-solutions/thrustmq/common"
	"github.com/rambler-digital-solutions/thrustmq/config"
	"github.com/rambler-digital-solutions/thrustmq/logging"
	"net"
)

var (
	Stage2CompressorChannel common.IntakeChannel = make(common.IntakeChannel, config.Intake.CompressorBuffer)
	CompressorChannel       common.IntakeChannel = make(common.IntakeChannel, config.Intake.CompressorBuffer)
)

func Init() {
	logging.Debug("Init intake")

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Intake.Port))
	common.FaceIt(err)

	go compressorStage1()
	go compressorStage2()

	var connection net.Conn

	for {
		connection, err = socket.Accept()
		common.FaceIt(err)
		go suck(connection)
	}
}
