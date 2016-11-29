package kite

import (
	"bytes"
	"testing"
)

func TestNewSrcCliWithoutUniq(t *testing.T) {
	var buf bytes.Buffer
	srv := NewServer(&buf)
	cli := NewClient(&buf)
	testSendData(srv, cli, t)
}

func TestNewSrcCliWithUniq(t *testing.T) {
	var buf bytes.Buffer
	uniq := Uniq{}
	copy(uniq[:], []byte("kittte"))
	srv := NewServer(&buf)
	cli := NewClient(&buf, uniq)
	testSendData(srv, cli, t)
}

func TestKiteSrvCli(t *testing.T) {
	var buf bytes.Buffer
	srv := NewBaseKite(&buf)
	cli := NewBaseKite(&buf)

	testHandshake(srv, cli, t)
	testSendData(srv, cli, t)
	testHeartBeat(srv, cli, t)
	testClose(srv, cli, t)
}

func testHandshake(srv *baseKite, cli *baseKite, t *testing.T) {
	var err error
	err = cli.sendHandshakeRequest()
	if err != nil {
		t.Fatalf("Client sendHandshakeRequest Err: %s\n", err)
	}

	err = srv.readHandshakeRequest()
	if err != nil {
		t.Fatalf("Server readHandshakeRequest Err: %s\n", err)
	}

	err = srv.sendHandshakeResponse()
	if err != nil {
		t.Fatalf("Server sendHandshakeResponse Err: %s\n", err)
	}

	err = cli.readHandshakeResponse()
	if err != nil {
		t.Fatalf("Client readHandshakeResponse Err: %s\n", err)
	}

	err = cli.sendHandshakeAck()
	if err != nil {
		t.Fatalf("Client sendHandshakeAck Err: %s\n", err)
	}

	err = srv.readHandshakeAck()
	if err != nil {
		t.Fatalf("Server readHandshakeAck Err: %s\n", err)
	}

	err = srv.sendHandshakeEnd()
	if err != nil {
		t.Fatalf("Server sendHandshakeEnd Err: %s\n", err)
	}

	err = cli.readHandshakeEnd()
	if err != nil {
		t.Fatalf("Client readHandshakeEnd Err: %s\n", err)
	}
}

func testSendData(srv Server, cli Client, t *testing.T) {
	var err error

	testMsg := []byte("Can this kite fly?")

	err = cli.SendData(testMsg, D_CHECK|D_ENCRYPTION)
	if err != nil {
		t.Fatalf("Client SendData Err: %s\n", err)
	}

	data, err := srv.ReadData()
	if err != nil {
		t.Fatalf("Server ReadData Err: %s\n", err)
	}
	if !EqualSlice(data, testMsg) {
		t.Fatalf("Server ReadData Not Equal!")
	}

	err = srv.SendDataAck()
	if err != nil {
		t.Fatalf("Server SendDataAck Err: %s\n", err)
	}

	err = cli.ReadDataAck()
	if err != nil {
		t.Fatalf("Client ReadDataAck Err: %s\n", err)
	}
}

func testHeartBeat(srv Server, cli Client, t *testing.T) {
	var err error

	err = cli.SendHeartBeat()
	if err != nil {
		t.Fatalf("Client SendHeartBeat Err: %s\n", err)
	}

	err = srv.ReadHeartBeat()
	if err != nil {
		t.Fatalf("Server ReadHeartBeat Err: %s\n", err)
	}

	err = srv.SendHeartBeatAck()
	if err != nil {
		t.Fatalf("Server SendHeartBeatAck Err: %s\n", err)
	}

	err = cli.ReadHeartBeatAck()
	if err != nil {
		t.Fatalf("Client ReadHeartBeatAck Err: %s\n", err)
	}
}

func testClose(srv Server, cli Client, t *testing.T) {
	var err error
	err = cli.SendClose()
	if err != nil {
		t.Fatalf("Client SendClose Err: %s\n", err)
	}

	err = srv.ReadClose()
	if err != nil {
		t.Fatalf("Server ReadClose Err: %s\n", err)
	}

	err = srv.SendCloseAck()
	if err != nil {
		t.Fatalf("Server SendCloseAck Err: %s\n", err)
	}

	err = cli.ReadCloseAck()
	if err != nil {
		t.Fatalf("Client ReadCloseAck Err: %s\n", err)
	}
}
