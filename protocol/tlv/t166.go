package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T166(imageType byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x166)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteByte(imageType)
		}))
	})
}
