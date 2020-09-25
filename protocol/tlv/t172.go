package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T172(rollbackSig []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x172)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.Write(rollbackSig)
		}))
	})
}
