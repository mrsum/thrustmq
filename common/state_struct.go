package common

import (
	"encoding/gob"
	"github.com/rambler-digital-solutions/thrustmq/config"
	"os"
	"time"
)

type StateStruct struct {
	Tail         uint64
	Head         uint64
	Capacity     float32
	ConnectionId uint64
}

var State StateStruct = loadState()

func loadState() StateStruct {
	if _, err := os.Stat(config.Exhaust.Chamber); err == nil {
		file, err := os.OpenFile(config.Exhaust.Chamber, os.O_RDONLY|os.O_CREATE, 0666)
		FaceIt(err)
		dec := gob.NewDecoder(file)
		result := StateStruct{}
		err = dec.Decode(&result)
		FaceIt(err)
		file.Close()
		return result
	}
	return StateStruct{Tail: 0, ConnectionId: 1}
}

func SaveState() {
	for {
		time.Sleep(1e9)
		file, err := os.OpenFile(config.Exhaust.Chamber, os.O_WRONLY|os.O_CREATE, 0666)
		FaceIt(err)
		enc := gob.NewEncoder(file)
		err = enc.Encode(State)
		FaceIt(err)
		file.Sync()
		file.Close()
	}
}

func (self *StateStruct) Distance() uint64 {
	return self.Head - self.Tail
}
