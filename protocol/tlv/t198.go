package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T198() []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x198)
		w.WriteTlv([]byte{0})
	})
}
