package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"log"
)

var padding = []byte{
	0x3c, 0x3c, 0x3c, 0x3c,
	0x3c, 0x3c, 0x3c, 0x3c}

type Coder struct {
	b cipher.Block
}

func New(key []byte) (*Coder, error) {
	b, err := aes.NewCipher(key)
	return &Coder{b}, err
}

func (c *Coder) Encrypt(domain int64, id int64) (string, error) {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, id); err != nil {
		return "", err
	}
	if err := binary.Write(buf, binary.LittleEndian, domain); err != nil {
		return "", err
	}
	in := buf.Bytes()
	out := make([]byte, len(in))
	c.b.Encrypt(out, in)
	// base64 encode
	bbUrl := &bytes.Buffer{}
	encoderUrl := base64.NewEncoder(base64.URLEncoding, bbUrl)
	encoderUrl.Write(out)
	encoderUrl.Close()
	s := bbUrl.String()
	//fmt.Println(id, in, out, s)
	return s[:len(s)-2], nil
}

func (c *Coder) Decrypt(enc string) (domain int64, key int64) {
	//  if (len(stringId) != 14) {
	//    fmt.Println("bug, stringId:", stringId)
	//  }
	// base64 decode
	buf := bytes.NewBufferString(enc + "==")
	decoderUrl := base64.NewDecoder(base64.URLEncoding, buf)
	decoded := make([]byte, 20)
	n, _ := decoderUrl.Read(decoded)
	//fmt.Println(decoded)
	decoded = decoded[:n]
	// decrypt
	decrypted := make([]byte, 16)
	c.b.Decrypt(decrypted, decoded[:n])
	decryptedKey := decrypted[:8]
	decryptedDomain := decrypted[8:]
	//fmt.Println(decoded[:n], decrypted)
	reader := bytes.NewReader(decryptedDomain)
	binary.Read(reader, binary.LittleEndian, &domain)
	reader = bytes.NewReader(decryptedKey)
	binary.Read(reader, binary.LittleEndian, &key)
	//fmt.Println(id)
	return domain, key
}

func (c *Coder) EncryptId(id int64) string {
	// encode as byte array
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, id); err != nil {
		log.Fatalf("binary.Write faild:", err)
	}
	buf.Write(padding)
	in := buf.Bytes()
	// encrypt id
	out := make([]byte, len(in))
	//fmt.Println(len(in), len(out))
	c.b.Encrypt(out, in)
	// base64 encode
	bbUrl := &bytes.Buffer{}
	encoderUrl := base64.NewEncoder(base64.URLEncoding, bbUrl)
	encoderUrl.Write(out)
	encoderUrl.Close()
	s := bbUrl.String()
	//fmt.Println(id, in, out, s)
	return s[:len(s)-2]
}

func (c *Coder) DecryptId(stringId string) int64 {
	//  if (len(stringId) != 14) {
	//    fmt.Println("bug, stringId:", stringId)
	//  }
	// base64 decode
	buf := bytes.NewBufferString(stringId + "==")
	decoderUrl := base64.NewDecoder(base64.URLEncoding, buf)
	decoded := make([]byte, 20)
	n, err := decoderUrl.Read(decoded)
	if err != nil {
		log.Fatalf("cannot decode url: %s", err)
	}
	//fmt.Println(decoded)
	decoded = decoded[:n]
	// decrypt
	decrypted := make([]byte, 16)
	c.b.Decrypt(decrypted, decoded[:n])
	decrypted = decrypted[:8]
	//fmt.Println(decoded[:n], decrypted)
	var id int64
	reader := bytes.NewReader(decrypted)
	binary.Read(reader, binary.LittleEndian, &id)
	//fmt.Println(id)
	return id
}
