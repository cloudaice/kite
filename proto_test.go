package kite

import (
	"bytes"
	"testing"
)

func TestWritePackage(t *testing.T) {
	var buf bytes.Buffer
	var err error
	pr := NewProto(&buf)
	err = pr.SetToken([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatalf("TestWritePackage SetToken Error: %s", err.Error())
	}

	err = pr.SetToken([]byte("token"))
	if err == nil {
		t.Fatalf("TestWritePackage SetToken Error: %s", err.Error())
	}

	pr.token = []byte{}
	err = pr.WritePackage(T_HANDSHAKEREQUEST, D_CHECK|D_ENCRYPTION, []byte("lovely day"))
	if err == nil {
		t.Fatal("TestWritePackage WritePackage Error")
	}
}

func TestReadPackage(t *testing.T) {
	buf := bytes.NewBuffer([]byte("test buffer"))
	pr := NewProto(buf)
	_, err := pr.ReadPackage()
	if err == nil {
		t.Fatalf("TestReadPackage Error")
	}
	buf = bytes.NewBuffer([]byte("head"))
	pr = NewProto(buf)
	_, err = pr.ReadPackage()
	if err == nil {
		t.Fatal("TestReadPackage Error")
	}

	buf.Reset()
	pr = NewProto(buf)
	err = pr.WritePackage(T_HEARTBEAT, D_CHECK|D_COMPRESSION|D_ENCRYPTION, nil)
	if err != nil {
		t.Fatalf("TestReadPackage Error")
	}
	err = pr.WritePackage(T_HEARTBEAT, D_NONE, nil)
	if err != nil {
		t.Fatalf("TestReadPackage Error")
	}

	_, err = pr.ReadPackage()
	if err != nil {
		t.Fatalf("TestReadPackage Error")
	}

	_, err = pr.ReadPackage()
	if err != nil {
		t.Fatalf("TestReadPackage Error")
	}
}

func TestProto(t *testing.T) {
	tds := mkdata()

	for _, td := range tds {
		var buf bytes.Buffer
		p := NewProto(&buf)

		err := p.WritePackage(td.tp, td.df, td.content)
		if err != nil {
			t.Fatalf("TestProto: WritePackage Error %s\n", err.Error())
		}
		pack, err := p.ReadPackage()
		assertWriteRead(
			td, pack,
			err, t,
		)
	}
}

type testData struct {
	tp      KiteType
	df      KiteDefine
	content []byte
}

func mkdata() []testData {
	var tds []testData

	tds = append(tds, testData{
		tp:      T_HANDSHAKEREQUEST,
		df:      D_CHECK,
		content: []byte("handshake"),
	})

	tds = append(tds, testData{
		tp:      T_HEARTBEAT,
		df:      D_CHECK | D_ENCRYPTION,
		content: []byte("heartbeat"),
	})

	tds = append(tds, testData{
		tp:      T_DATA,
		df:      D_CHECK | D_ENCRYPTION | D_COMPRESSION,
		content: []byte("command"),
	})

	tds = append(tds, testData{
		tp:      T_DATAACK,
		df:      D_CHECK | D_ENCRYPTION | D_COMPRESSION,
		content: []byte("response"),
	})

	tds = append(tds, testData{
		tp:      T_DATAACK,
		df:      D_NONE,
		content: nil,
	})

	return tds
}

func assertWriteRead(
	td testData,
	pack *ProtoPackage,
	err error,
	t *testing.T,
) {
	t.Logf("Assert %d %s\n", td.tp, string(td.content))
	if err != nil {
		t.Fatalf("TestProto ReadPackage Error %s", err)
	}

	if td.tp != pack.PackageType() {
		t.Fatalf("TestProto tppe Error")
	}

	if td.df&D_CHECK == D_CHECK {
		if !pack.IfCheck() {
			t.Fatalf("TestProto check Error")
		}
		if !EqualSlice(
			pack.Checksum(),
			Md5(pack.Content()),
		) {
			t.Fatalf("TestProto checksum Error")
		}
	}

	if td.df&D_ENCRYPTION == D_ENCRYPTION && !pack.IfEncrypt() {
		t.Fatalf("TestProto encrypt Error")
	}

	if td.df&D_COMPRESSION == D_COMPRESSION && !pack.IfCompress() {
		t.Fatalf("TestProto compress Error")
	}

	if string(td.content) != string(pack.Content()) {
		t.Fatalf("TestProto Read data Error %s", err)
	}
}
