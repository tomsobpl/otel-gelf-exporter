package tlsgateway

import (
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"net"
)

type Endpoint struct {
	Network  string
	Endpoint string
}

type TLSGateway struct {
	cancel   context.CancelFunc
	conn     net.Conn
	endpoint Endpoint
	listener net.Listener
	logger   *zap.Logger
}

func (g *TLSGateway) Addr() net.Addr {
	return g.listener.Addr()
}

func (g *TLSGateway) Start(config *tls.Config) error {
	var err error

	if g.conn, err = tls.Dial(g.endpoint.Network, g.endpoint.Endpoint, config); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go g.run(ctx)

	g.cancel = cancel

	return nil
}

func (g *TLSGateway) Shutdown() error {
	g.cancel()
	return g.listener.Close()
}

func (g *TLSGateway) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if err := g.conn.Close(); err != nil {
				g.logger.Error("failed to close local connection", zap.Error(err))
			}

			return
		default:
			conn, err := g.listener.Accept()

			if err != nil {
				g.logger.Error("failed to accept connection", zap.Error(err))
			}

			g.forward(conn, g.conn)

			if err := conn.Close(); err != nil {
				g.logger.Error("failed to close remote connection", zap.Error(err))
			}
		}
	}
}

func (g *TLSGateway) forward(src net.Conn, dst net.Conn) {
	srcChannel := connectionIntoChannel(src)
	dstChannel := connectionIntoChannel(dst)

	for {
		select {
		case b1 := <-srcChannel:
			if b1 == nil {
				return
			} else {
				if _, err := dst.Write(b1); err != nil {
					g.logger.Warn("failed to write to destination", zap.Error(err))
				}
			}
		case b2 := <-dstChannel:
			if b2 == nil {
				return
			} else {
				if _, err := src.Write(b2); err != nil {
					g.logger.Warn("failed receiving from destination", zap.Error(err))
				}
			}
		}
	}
}

func NewTLSGateway(local Endpoint, remote Endpoint, logger *zap.Logger) (*TLSGateway, error) {
	listener, err := net.Listen(local.Network, local.Endpoint)

	if err != nil {
		return nil, err
	}

	return &TLSGateway{
		endpoint: remote,
		listener: listener,
		logger:   logger,
	}, nil
}
