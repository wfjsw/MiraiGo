package tlv

import "github.com/wfjsw/MiraiGo/binary"

func T147(appId uint32, apkVersionName, apkSignatureMd5 []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x147)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(appId)
			w.WriteTlvLimitedSize(apkVersionName, 32)
			w.WriteTlvLimitedSize(apkSignatureMd5, 32)
		}))
	})
}
