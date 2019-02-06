package trace

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// /reqid
var pid = uint32(time.Now().UnixNano() % 4294967291)

// GenReqID creates a decodeable request id
func GenReqID() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func DecodeReqID(reqID string) (output string, err error) {
	var decodedBytes []byte
	var unixNano int64

	decodedBytes, err = base64.URLEncoding.DecodeString(reqID)
	if err != nil || len(decodedBytes) < 4 {
		return
	}
	newBytes := decodedBytes[4:]
	newBytesLen := len(newBytes)
	newStr := ""
	for i := newBytesLen - 1; i >= 0; i-- {
		newStr += fmt.Sprintf("%02X", newBytes[i])
	}
	unixNano, err = strconv.ParseInt(newStr, 16, 64)
	if err != nil {
		return
	}
	dstDate := time.Unix(0, unixNano)
	output = dstDate.Format("2006-01-02/15-04")
	return
}
