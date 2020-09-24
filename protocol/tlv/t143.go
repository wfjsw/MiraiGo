package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T143(arr []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x143)
		w.WriteTlv(arr)
	})
}
