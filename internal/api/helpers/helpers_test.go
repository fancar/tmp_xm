package helpers

import (
	"context"
	"net"
	"testing"

	"github.com/fancar/tmp_xm/internal/storage"
	"github.com/fancar/tmp_xm/internal/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

func TestIPaddrCheker(t *testing.T) {
	// ctx := context.Background()
	assert := require.New(t)
	conf := test.GetConfig()
	assert.NoError(storage.Setup(conf))

	assert = require.New(t)
	allowedIP, _ := net.ResolveIPAddr("ip", "104.28.106.39")
	allowed := peer.Peer{
		Addr: allowedIP,
	}
	assert.NoError(IPaddrCheker(peer.NewContext(context.Background(), &allowed)))

	assert = require.New(t)
	deniedIP, _ := net.ResolveIPAddr("ip", "8.8.8.8")
	denied := peer.Peer{
		Addr: deniedIP,
	}

	assert.NoError(IPaddrCheker(peer.NewContext(context.Background(), &denied)))
}
