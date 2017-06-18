package crypto

import (
	"testing"
)

var KEY = []byte{
	0x2c, 0x88, 0x25, 0x1a,
	0xaa, 0xae, 0xc2, 0xa2,
	0xaf, 0xe7, 0x84, 0x8a,
	0x10, 0xcf, 0xe3, 0x2a}

func Test_EncryptId(t *testing.T) {
	c, _ := New(KEY)
	var id int64 = 123123123123
	encodedId := c.EncryptId(id)
	pureId := c.DecryptId(encodedId)
	if id != pureId {
		t.Errorf("id want: %d, got %d", encodedId, pureId)
		t.FailNow()
	}
}

func Test_Encrypt(t *testing.T) {
	c, _ := New(KEY)
	var id int64 = 123123123123
	var domain int64 = 321321
	encodedId, err := c.Encrypt(domain, id)
	if err != nil {
		t.Errorf("failed: %s", err)
		t.Fail()
	}
	pureDomain, pureId := c.Decrypt(encodedId)
	if domain != pureDomain {
		t.Errorf("domain got %d, want %d", pureDomain, domain)
		t.Fail()
	}
	if id != pureId {
		t.Errorf("id got %d, want %d", pureId, id)
		t.Fail()
	}
}
