package crypto

import "github.com/wfjsw/MiraiGo/binary"

type EncryptSession struct {
	T133 []byte
}

func NewEncryptSession(t133 []byte) *EncryptSession {
	return &EncryptSession{T133: t133}
}

func (e *EncryptSession) DoEncrypt(d, k []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		encrypt := binary.NewTeaCipher(k).Encrypt(d)
		w.WriteUInt16(uint16(len(e.T133)))
		w.Write(e.T133)
		w.Write(encrypt)
	})
}

func (e *EncryptSession) Id() byte {
	return 69
}
