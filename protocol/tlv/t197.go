package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T197() []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x197)
		w.WriteTlv([]byte{0})
	})
}
