package tests

import (
	"github.com/rambler-digital-solutions/thrustmq/common"
	"github.com/rambler-digital-solutions/thrustmq/config"
	"os"
	"testing"
)

func TestRecordSerialization(t *testing.T) {
	indexFile, err := os.OpenFile(config.Base.Index+"_test", os.O_RDWR|os.O_CREATE, 0666)
	common.FaceIt(err)
	defer indexFile.Close()

	record := common.Record{}

	slots := record.Slots()
	for i := 0; i < len(slots); i++ {
		*slots[i] = uint64(i + 1)
	}

	indexFile.Seek(0, os.SEEK_SET)
	indexFile.Write(record.Serialize())

	indexFile.Seek(0, os.SEEK_SET)
	readRecord := common.Record{}
	readRecord.Deserialize(indexFile)

	readSlots := readRecord.Slots()
	for i := 0; i < len(slots); i++ {
		if *slots[i] != *readSlots[i] {
			t.Fatalf("deserialized field %d ne expected %d", *readSlots[i], *slots[i])
		}
	}
}
