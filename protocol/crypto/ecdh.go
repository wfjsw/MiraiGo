package crypto

import "github.com/wfjsw/MiraiGo/binary"

type EncryptECDH struct {
	InitialShareKey []byte
	PublicKey       []byte
}

func (e *EncryptECDH) DoEncrypt(d, k []byte) []byte {
	w := binary.NewWriter()
	w.WriteByte(0x01)
	w.WriteByte(0x01)
	w.Write(k)
	w.WriteUInt16(258)
	w.WriteUInt16(uint16(len(ECDH.PublicKey)))
	w.Write(ECDH.PublicKey)
	w.EncryptAndWrite(ECDH.InitialShareKey, d)
	return w.Bytes()
}

func (e *EncryptECDH) Id() byte {
	return 7
}
