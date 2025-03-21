package gelftcpexporter

import (
	"context"
	"crypto/tls"
	"net"
)

type TLSGatewayEndpoint struct {
	Network  string
	Endpoint string
}

type TLSGateway struct {
	endpoint TLSGatewayEndpoint
	listener *net.Listener
	cancel   context.CancelFunc
	conn     net.Conn
}

func (g *TLSGateway) Addr() net.Addr {
	return (*g.listener).Addr()
}

func (g *TLSGateway) Start() error {
	var err error

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	if g.conn, err = tls.Dial(g.endpoint.Network, g.endpoint.Endpoint, conf); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go g.run(ctx)

	g.cancel = cancel

	return nil
}

func (g *TLSGateway) Shutdown() error {
	g.cancel()
	(*g.listener).Close()
	return nil
}

func (g *TLSGateway) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			g.conn.Close()
			return
		default:
			conn, err := (*g.listener).Accept()

			if err != nil {
				// log error
				continue
			}

			g.forward(conn, g.conn)

			if err := conn.Close(); err != nil {
				// log error
			}
		}
	}

	//for {
	//	conn, err := (*g.listener).Accept()
	//
	//	if err != nil {
	//		// log error
	//		continue
	//	}
	//
	//	g.forward(conn, g.conn)
	//
	//	if err := conn.Close(); err != nil {
	//		// log error
	//	}
	//}
}

//func (g *TLSGateway) Run(wg *sync.WaitGroup) error {
//	var err error
//
//	conf := &tls.Config{
//		InsecureSkipVerify: true,
//	}
//
//	if g.conn, err = tls.Dial(g.endpoint.Network, g.endpoint.Endpoint, conf); err != nil {
//		return err
//	}
//
//	wg.Done()
//
//	ctx, cancel := context.WithCancel(context.Background())
//
//	for {
//		conn, err := (*g.listener).Accept()
//
//		if err != nil {
//			// log error
//			continue
//		}
//
//		g.forward(conn, g.conn)
//
//		if err := conn.Close(); err != nil {
//			// log error
//		}
//	}
//}

//func (g *TLSGateway) Shutdown() error {
//
//}

func (g *TLSGateway) forward(src net.Conn, dst net.Conn) {
	srcChannel := connectionIntoChannel(src)
	dstChannel := connectionIntoChannel(dst)

	for {
		select {
		case b1 := <-srcChannel:
			if b1 == nil {
				return
			} else {
				dst.Write(b1)
			}
		case b2 := <-dstChannel:
			if b2 == nil {
				return
			} else {
				src.Write(b2)
			}
		}
	}
}

func NewTLSGateway(local TLSGatewayEndpoint, remote TLSGatewayEndpoint) (*TLSGateway, error) {
	listener, err := net.Listen(local.Network, local.Endpoint)

	if err != nil {
		return nil, err
	}

	return &TLSGateway{
		endpoint: remote,
		listener: &listener,
	}, nil
}
