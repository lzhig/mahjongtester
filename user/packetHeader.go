package user

import "io"

type MahjongPacketHeader struct {
}

func (this *MahjongPacketHeader) BuildHeader(bodyLen int, data []byte) error {
	if len(data) < this.GetHeaderLen() {
		return io.EOF
	}
	bodyLen += this.GetHeaderLen()
	ndx := 0
	data[ndx] = byte((bodyLen & 0xFF00) >> 8)
	data[ndx+1] = byte(bodyLen & 0xFF)
	return nil
}

func (this *MahjongPacketHeader) ParsePacketHeader(data []byte) (ok bool, headerLen int32, packetLen int32, err error) {
	totalLen := this.GetHeaderLen()
	if len(data) < totalLen {
		return false, int32(totalLen), 0, nil
	}

	packetLen = ((int32(data[0]) << 8) + int32(data[1])) - int32(totalLen)
	return true, int32(totalLen), packetLen, nil
}

func (this *MahjongPacketHeader) GetHeaderLen() int {
	return 2
}
