package kite

import (
	"testing"
)

func TestHardwareAddr(t *testing.T) {
	addr := HardwareAddr()
	t.Logf("HardwareAddr: %v\n", addr)
}

func TestRandBytes(t *testing.T) {
	b0 := RandBytes(20)
	b1 := RandBytes(20)
	if EqualSlice(b0, b1) {
		t.Errorf("RandBytes Not Random")
	}
}

func TestMd5(t *testing.T) {
	data := []byte("testdata")
	md := Md5(data)
	if len(md) == 0 {
		t.Error("Empty Md5")
	}
}

func TestAesCrypt(t *testing.T) {
	data := []byte("testdata")
	key := []byte("0123456789abcdef")
	iv := []byte("abcdef0123456789")

	endata, err := AesEncrypt(data, key, iv)
	if err != nil {
		t.Fatalf("AesEncrypt Err: %s\n", err)
	}

	dedata, err := AesDecrypt(endata, key, iv)
	if err != nil {
		t.Fatalf("AesDecrypt Err: %s\n", err)
	}

	if !EqualSlice(data, dedata) {
		t.Errorf("AesCrypt unequal!")
	}
}

func TestEqualSlice(t *testing.T) {
	if !EqualSlice([]byte("123"), []byte("123")) {
		t.Error("EqualSlice Fail!")
	}
	if EqualSlice([]byte("123\n"), []byte("122")) {
		t.Error("EqualSlice Fail!")
	}
}

func TestGzip(t *testing.T) {
	data := []byte("testdata")

	gzd, err := Gzip(data)
	if err != nil {
		t.Fatalf("Gzip err: %s\n", err)
	}

	ugzd, err := UnGzip(gzd)
	if err != nil {
		t.Fatalf("UnGzip err: %s\n", err)
	}

	if !EqualSlice(data, ugzd) {
		t.Errorf("Gzip unequal!")
	}
}

func TestRsa(t *testing.T) {
	data := []byte("test rsa data")
	pubKey, privKey, err := GenRsaKey(1024)
	if err != nil {
		t.Fatalf("GenRsaKey Fail: %s\n", err)
	}

	ciphertext, err := RsaEncrypt(data, pubKey)
	if err != nil {
		t.Fatalf("RsaEncrypt Fail: %s\n", err)
	}

	final, err := RsaDecrypt(ciphertext, privKey)
	if err != nil {
		t.Fatalf("RsaDecrypt Fail: %s\n", err)
	}

	if !EqualSlice(final, data) {
		t.Error("Not Equal!")
	}
}
