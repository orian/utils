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
	encodedID := c.EncryptId(id)
	pureID := c.DecryptId(encodedID)
	if id != pureID {
		t.Errorf("id want: %s, got %d", encodedID, pureID)
		t.FailNow()
	}
}

func Test_Encrypt(t *testing.T) {
	c, _ := New(KEY)
	var id int64 = 123123123123
	var domain int64 = 321321
	encodedID, err := c.Encrypt(domain, id)
	if err != nil {
		t.Errorf("failed: %s", err)
		t.Fail()
	}
	pureDomain, pureID := c.Decrypt(encodedID)
	if domain != pureDomain {
		t.Errorf("domain got %d, want %d", pureDomain, domain)
		t.Fail()
	}
	if id != pureID {
		t.Errorf("id got %d, want %d", pureID, id)
		t.Fail()
	}
}
