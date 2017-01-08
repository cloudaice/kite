/*
   |    2 bytes     |   4 bytes   |
   |  command define |  data length |
*/

package kite

import (
	"io"

	"github.com/cloudaice/kite/binaryproto"
)

type KiteType uint8
type KiteDefine uint8

const (
	T_HANDSHAKEREQUEST  KiteType = 0x01
	T_HANDSHAKERESPONSE KiteType = 0x02
	T_HANDSHAKEACK      KiteType = 0x03
	T_HANDSHAKEEND      KiteType = 0x04
	T_HEARTBEAT         KiteType = 0x05
	T_HEARTBEATACK      KiteType = 0x06
	T_DATA              KiteType = 0x07
	T_DATAACK           KiteType = 0x08
	T_CLOSE             KiteType = 0x09
	T_CLOSEACK          KiteType = 0x0a

	D_COMPRESSION KiteDefine = 0x04
	D_ENCRYPTION  KiteDefine = 0x02
	D_CHECK       KiteDefine = 0x01
	D_NONE        KiteDefine = 0x00
)

var (
	IV = []byte("abcdef0123456789")
)

type ProtoPackage struct {
	kiteType    KiteType
	compression bool
	encryption  bool
	check       bool
	checksum    []byte
	content     []byte
}

func (pp *ProtoPackage) PackageType() KiteType {
	return pp.kiteType
}

func (pp *ProtoPackage) IfCompress() bool {
	return pp.compression
}

func (pp *ProtoPackage) IfEncrypt() bool {
	return pp.encryption
}

func (pp *ProtoPackage) IfCheck() bool {
	return pp.check
}

func (pp *ProtoPackage) Content() []byte {
	return pp.content
}

func (pp *ProtoPackage) Checksum() []byte {
	return pp.checksum
}

type Proto struct {
	br    *binaryproto.Reader
	bw    *binaryproto.Writer
	token []byte
}

func NewProto(rw io.ReadWriter) *Proto {
	return &Proto{
		br:    binaryproto.NewReader(rw),
		bw:    binaryproto.NewWriter(rw),
		token: []byte("0123456789abcdef"),
	}
}

func (pr *Proto) GetToken() []byte {
	return pr.token
}

func (pr *Proto) SetToken(token []byte) error {
	if len(token) != 16 {
		return PError("Invailid Token Size")
	}
	pr.token = token
	return nil
}

// check define, the last three bits are compression, encryption, check
func (pr *Proto) checkDefine(define KiteDefine) (compression, encryption, check bool) {
	compression, encryption, check = false, false, false

	if define&D_CHECK == D_CHECK {
		check = true
	}
	if define&D_ENCRYPTION == D_ENCRYPTION {
		encryption = true
	}

	if define&D_COMPRESSION == D_COMPRESSION {
		compression = true
	}

	return
}

func (pr *Proto) ReadPackage() (*ProtoPackage, error) {
	header, err := pr.br.ReadHeader()
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, PError("ReadPackage Fail: ").E(err)
	}

	compression, encryption, check := pr.checkDefine(KiteDefine(header.Define))

	body, err := pr.br.ReadBody(header.Length, check)
	if err != nil {
		return nil, err
	}

	if compression {
		body.Data, err = UnGzip(body.Data)
		if err != nil {
			return nil, ErrCompression.E(err)
		}
	}

	if encryption {
		body.Data, err = AesDecrypt(body.Data, pr.token, IV)
		if err != nil {
			return nil, ErrEncryption.E(err)
		}
	}

	if check {
		if !EqualSlice(body.CheckSum, Md5(body.Data)) {
			return nil, ErrCheck
		}
	}

	return &ProtoPackage{
		kiteType:    KiteType(header.Type),
		compression: compression,
		encryption:  encryption,
		check:       check,
		checksum:    body.CheckSum,
		content:     body.Data,
	}, nil
}

func (pr *Proto) WritePackage(tp KiteType, df KiteDefine, data []byte) error {

	var checksum []byte
	var err error
	compression, encryption, check := pr.checkDefine(df)
	if check {
		checksum = Md5(data)
	}
	if encryption {
		data, err = AesEncrypt(data, pr.token, IV)
		if err != nil {
			return ErrEncryption.E(err)
		}
	}
	if compression {
		data, err = Gzip(data)
		if err != nil {
			return ErrCompression.E(err)
		}
	}

	body := &binaryproto.Body{
		CheckSum: checksum,
		Data:     data,
	}
	header := binaryproto.NewHeader(uint8(tp), uint8(df), body.Length())

	err = pr.write(header.BinaryCode())
	if err != nil {
		return err
	}
	return pr.write(body.BinaryCode())
}

func (pr *Proto) write(data []byte) error {
	err := pr.bw.Write(data)
	if err != nil {
		if err == io.EOF {
			return err
		}
		return PError("WritePackage Fail: ").E(err)
	}
	return nil
}
