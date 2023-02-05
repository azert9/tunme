package test

import (
	"github.com/azert9/tunme/pkg/tunme"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFunctional(t *testing.T) {

	if os.Getenv("ENABLE_FUNCTIONAL_TESTS") != "yes" {
		t.SkipNow()
	}

	client, err := tunme.OpenTunnel("tcp-client,localhost:5000")
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		assert.NoError(t, client.Close())
	}()

	server, err := tunme.OpenTunnel("tcp-server,:5000")
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		assert.NoError(t, server.Close())
	}()

	// TODO
}
