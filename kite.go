package kite

import (
	"bufio"
	"io"
)

const (
	TokenSize = 16
)

// md5(uniq) as original token
// baseKite implement Kite Client Server interface
type baseKite struct {
	uniq Uniq
	pr   *Proto
	buf  *bufio.ReadWriter
}

type Kite interface {
	ReadPackage() (*ProtoPackage, error)
	WritePackage(KiteType, KiteDefine, []byte) error

	ReadData() ([]byte, error)
	SendData([]byte, KiteDefine) error

	ReadDataAck() error
	SendDataAck() error

	ReadHeartBeat() error
	SendHeartBeat() error

	ReadHeartBeatAck() error
	SendHeartBeatAck() error

	ReadClose() error
	SendClose() error

	ReadCloseAck() error
	SendCloseAck() error
}

type Server interface {
	Kite
	WaitHandshake() error
}

type Client interface {
	Kite
	SendHandshake() error
}

func NewBaseKite(rw io.ReadWriter) *baseKite {
	buf := bufio.NewReadWriter(bufio.NewReader(rw), bufio.NewWriter(rw))
	return &baseKite{
		buf: buf,
		pr:  NewProto(buf),
	}
}

func NewServer(rw io.ReadWriter) Server {
	return NewBaseKite(rw)
}

// TODO: uniq can be a func type return  Uniq
// NewClient return a kite client
func NewClient(rw io.ReadWriter, uniq ...Uniq) Client {
	cli := NewBaseKite(rw)
	if len(uniq) == 0 || len(uniq[0]) != UniqSize {
		cli.uniq = HardwareAddr()
	} else {
		cli.uniq = uniq[0]
	}
	return cli
}

// WaitHandshake called by Server, then wait for shakehanding from client
func (srv *baseKite) WaitHandshake() error {
	err := srv.readHandshakeRequest()
	if err != nil {
		return err
	}

	err = srv.sendHandshakeResponse()
	if err != nil {
		return err
	}

	err = srv.readHandshakeAck()
	if err != nil {
		return err
	}
	err = srv.sendHandshakeEnd()
	return err
}

// readHandshakeRequest get Client's uniq, and set original token as md5(uniq)
func (srv *baseKite) readHandshakeRequest() error {
	pp, err := srv.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HANDSHAKEREQUEST || len(pp.Content()) != UniqSize {
		return ErrUnexpectedPackage
	}
	copy(srv.uniq[:], pp.Content())

	return srv.pr.SetToken(Md5(pp.Content()))
}

// sendHandshakeResponse generate random token, encrypt(uniq+token) with md5(uniq),
// this package should be checked and encrypt
func (srv *baseKite) sendHandshakeResponse() error {
	token := RandBytes(TokenSize)
	data := append(srv.uniq[:], token...)
	err := srv.WritePackage(T_HANDSHAKERESPONSE, D_CHECK|D_ENCRYPTION, data)
	if err != nil {
		return err
	}
	return srv.pr.SetToken(token)
}

// readHandshakeAck get data from Client, the data should be md5(token)
func (srv *baseKite) readHandshakeAck() error {
	pp, err := srv.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HANDSHAKEACK || !EqualSlice(Md5(srv.pr.GetToken()), pp.Content()) {
		return ErrUnexpectedPackage
	}
	return nil
}

// sendHandshakeEnd write package with empty content, finish Handshake.
func (srv *baseKite) sendHandshakeEnd() error {
	err := srv.WritePackage(T_HANDSHAKEEND, D_NONE, nil)
	return err
}

// Sessions are always initiated by Client sending HandshakeRequest
func (cli *baseKite) SendHandshake() error {
	err := cli.sendHandshakeRequest()
	if err != nil {
		return err
	}
	err = cli.readHandshakeResponse()
	if err != nil {
		return err
	}
	err = cli.sendHandshakeAck()
	if err != nil {
		return err
	}
	err = cli.readHandshakeEnd()
	return err
}

// sendHandshakeRequest send uniq to Server, set original token as md5(uniq),
// this package's mark is D_DONE
func (cli *baseKite) sendHandshakeRequest() error {
	cli.pr.SetToken(Md5(cli.uniq[:]))
	return cli.WritePackage(T_HANDSHAKEREQUEST, D_NONE, cli.uniq[:])
}

// TODO: error info should be more detail
// readHandshakeResponse read a package from server, the package's data should be (uniq+token)
func (cli *baseKite) readHandshakeResponse() error {
	pp, err := cli.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HANDSHAKERESPONSE ||
		len(pp.Content()) != UniqSize+TokenSize ||
		!EqualSlice(cli.uniq[:], pp.Content()[0:UniqSize]) {

		return ErrUnexpectedPackage
	}
	cli.pr.SetToken(pp.Content()[UniqSize:])
	return nil
}

// sendHandshakeAck send md5(token) to Server, this package should be check and encryt
func (cli *baseKite) sendHandshakeAck() error {
	data := Md5(cli.pr.GetToken())
	return cli.WritePackage(T_HANDSHAKEACK, D_CHECK|D_ENCRYPTION, data)
}

// readHandshakeEnd read HandshakeEnd package. Handshake finished
func (cli *baseKite) readHandshakeEnd() error {
	pp, err := cli.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HANDSHAKEEND {
		return ErrUnexpectedPackage
	}
	return nil
}

func (bk *baseKite) ReadPackage() (*ProtoPackage, error) {
	return bk.pr.ReadPackage()
}

// WritePackage call Flush() at the end, because it use bufio
func (bk *baseKite) WritePackage(tp KiteType, df KiteDefine, data []byte) error {
	err := bk.pr.WritePackage(tp, df, data)
	if err != nil {
		return err
	}
	return bk.buf.Flush()
}

func (bk *baseKite) SendData(data []byte, df KiteDefine) error {
	return bk.WritePackage(T_DATA, df, data)
}

func (bk *baseKite) SendDataAck() error {
	return bk.WritePackage(T_DATAACK, D_NONE, nil)
}
func (bk *baseKite) SendHeartBeat() error {
	return bk.WritePackage(T_HEARTBEAT, D_NONE, nil)
}

func (bk *baseKite) SendHeartBeatAck() error {
	return bk.WritePackage(T_HEARTBEATACK, D_NONE, nil)
}

func (bk *baseKite) SendClose() error {
	return bk.WritePackage(T_CLOSE, D_NONE, nil)
}

func (bk *baseKite) SendCloseAck() error {
	return bk.WritePackage(T_CLOSEACK, D_NONE, nil)
}

// ReadData read data package else return err
func (bk *baseKite) ReadData() ([]byte, error) {
	pp, err := bk.ReadPackage()
	if err != nil {
		return nil, err
	}
	if pp.PackageType() != T_DATA {
		return nil, ErrUnexpectedPackage
	}
	return pp.Content(), nil
}

// ReadDataAck read a DataAck package else return err
func (bk *baseKite) ReadDataAck() error {
	pp, err := bk.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_DATAACK {
		return ErrUnexpectedPackage
	}
	return nil
}

// ReadHeartBeat read a HeartBeat package else return err
func (bk *baseKite) ReadHeartBeat() error {
	pp, err := bk.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HEARTBEAT {
		return ErrUnexpectedPackage
	}
	return nil
}

// ReadHeartBeatAck read a HeartBeatAck type package else return err
func (bk *baseKite) ReadHeartBeatAck() error {
	pp, err := bk.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_HEARTBEATACK {
		return ErrUnexpectedPackage
	}
	return nil
}

// ReadClose read a Close type package else return err
func (bk *baseKite) ReadClose() error {
	pp, err := bk.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_CLOSE {
		return ErrUnexpectedPackage
	}
	return nil
}

// ReadCloseAck read a CloseAck type package else return err
func (bk *baseKite) ReadCloseAck() error {
	pp, err := bk.ReadPackage()
	if err != nil {
		return err
	}
	if pp.PackageType() != T_CLOSEACK {
		return ErrUnexpectedPackage
	}
	return nil
}
