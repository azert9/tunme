package stream_link

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
)

func testAcceptStream(t *testing.T, roleOrder int) {

	tun1, tun2 := newMockTunPair(roleOrder)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		_, err := tun1.AcceptStream()
		if err != nil {
			t.Logf("error accepting stream: %v", err)
			t.Fail()
		}
	}()

	_, err := tun2.OpenStream()
	if err != nil {
		t.Logf("error opening stream: %v", err)
		t.Fail()
	}
}

func TestAcceptStream(t *testing.T) {

	for roleOrder := 0; roleOrder < 2; roleOrder++ {
		t.Run(fmt.Sprintf("role order %d", roleOrder), func(t *testing.T) {
			testAcceptStream(t, roleOrder)
		})
	}
}

func testSendDataThroughStream(t *testing.T, roleOrder int) {

	tun1, tun2 := newMockTunPair(roleOrder)

	rnd := rand.New(rand.NewSource(0))
	ref := make([]byte, 20000)
	rnd.Read(ref)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tun1.AcceptStream()
		if err != nil {
			t.Logf("error accepting stream: %v", err)
			t.Fail()
		}

		received := make([]byte, len(ref)*2)
		n, err := conn.Read(received)
		if err != nil {
			t.Logf("error receiving data: %v", err)
			t.Fail()
		}

		if !bytes.Equal(received[:n], ref) {
			t.Logf("received and sent not matching")
			t.Fail()
		}
	}()

	conn, err := tun2.OpenStream()
	if err != nil {
		t.Logf("error opening stream: %v", err)
		t.Fail()
	}

	if _, err := conn.Write(ref); err != nil {
		t.Logf("error sending data: %v", err)
		t.Fail()
	}
}

func TestSendDataThroughStream(t *testing.T) {

	for roleOrder := 0; roleOrder < 2; roleOrder++ {
		t.Run(fmt.Sprintf("role order %d", roleOrder), func(t *testing.T) {
			testSendDataThroughStream(t, roleOrder)
		})
	}
}

func testEchoDataThroughStream(t *testing.T, roleOrder int) {

	tun1, tun2 := newMockTunPair(roleOrder)

	rnd := rand.New(rand.NewSource(0))
	ref := make([]byte, 20000)
	rnd.Read(ref)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tun1.AcceptStream()
		if err != nil {
			t.Logf("error accepting stream: %v", err)
			t.Fail()
		}

		received := make([]byte, len(ref)*2)
		n, err := conn.Read(received)
		if err != nil {
			t.Logf("error receiving data: %v", err)
			t.Fail()
		}

		if _, err := conn.Write(received[:n]); err != nil {
			t.Logf("error sending data: %v", err)
			t.Fail()
		}

		if err := conn.Close(); err != nil {
			t.Logf("error closing stream: %v", err)
			t.Fail()
		}
	}()

	conn, err := tun2.OpenStream()
	if err != nil {
		t.Logf("error opening stream: %v", err)
		t.Fail()
	}

	if _, err := conn.Write(ref); err != nil {
		t.Logf("error sending data: %v", err)
		t.Fail()
	}

	received := make([]byte, len(ref)*2)
	n, err := conn.Read(received)
	if err != nil {
		t.Logf("error receiving data: %v", err)
		t.Fail()
	}

	if !bytes.Equal(received[:n], ref) {
		t.Logf("received and sent not matching")
		t.Fail()
	}

	if err := conn.Close(); err != nil {
		t.Logf("error closing stream: %v", err)
		t.Fail()
	}
}

func TestEchoDataThroughStream(t *testing.T) {

	for roleOrder := 0; roleOrder < 2; roleOrder++ {
		t.Run(fmt.Sprintf("role order %d", roleOrder), func(t *testing.T) {
			testEchoDataThroughStream(t, roleOrder)
		})
	}
}

func testSendFragmentedDataThroughStream(t *testing.T, roleOrder int) {

	tun1, tun2 := newMockTunPair(roleOrder)

	rnd := rand.New(rand.NewSource(0))
	ref := make([]byte, 20000)
	rnd.Read(ref)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tun1.AcceptStream()
		if err != nil {
			t.Logf("error accepting stream: %v", err)
			t.Fail()
		}

		received, err := io.ReadAll(conn)
		if err != nil {
			t.Logf("error receiving data: %v", err)
			t.Fail()
		}

		if !bytes.Equal(received, ref) {
			t.Logf("received and sent not matching")
			t.Fail()
		}

		if err := conn.Close(); err != nil {
			t.Logf("error closing stream: %v", err)
			t.Fail()
		}
	}()

	conn, err := tun2.OpenStream()
	if err != nil {
		t.Logf("error opening stream: %v", err)
		t.Fail()
	}

	chunksSize := 10

	for offset := 0; offset < len(ref); offset += chunksSize {

		if _, err := conn.Write(ref[offset : offset+chunksSize]); err != nil {
			t.Logf("error sending data: %v", err)
			t.Fail()
		}
	}

	if err := conn.Close(); err != nil {
		t.Logf("error closing stream: %v", err)
		t.Fail()
	}
}

func TestSendFragmentedDataThroughStream(t *testing.T) {

	for roleOrder := 0; roleOrder < 2; roleOrder++ {
		t.Run(fmt.Sprintf("role order %d", roleOrder), func(t *testing.T) {
			testSendFragmentedDataThroughStream(t, roleOrder)
		})
	}
}

func testStreamEof(t *testing.T, roleOrder int) {

	tun1, tun2 := newMockTunPair(roleOrder)

	rnd := rand.New(rand.NewSource(0))
	ref := make([]byte, 20000)
	rnd.Read(ref)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		conn, err := tun1.AcceptStream()
		if err != nil {
			t.Logf("error accepting stream: %v", err)
			t.Fail()
		}

		_, err = conn.Read(make([]byte, 1))
		if err != io.EOF {
			if err != nil {
				t.Logf("error receiving data: %v", err)
				t.Fail()
			} else {
				t.Logf("expecting EOF but did not get any error")
				t.Fail()
			}
		}

		if err := conn.Close(); err != nil {
			t.Logf("error closing stream: %v", err)
			t.Fail()
		}
	}()

	conn, err := tun2.OpenStream()
	if err != nil {
		t.Logf("error opening stream: %v", err)
		t.Fail()
	}

	if err := conn.Close(); err != nil {
		t.Logf("error closing stream: %v", err)
		t.Fail()
	}
}

func TestStreamEof(t *testing.T) {

	for roleOrder := 0; roleOrder < 2; roleOrder++ {
		t.Run(fmt.Sprintf("role order %d", roleOrder), func(t *testing.T) {
			testStreamEof(t, roleOrder)
		})
	}
}
