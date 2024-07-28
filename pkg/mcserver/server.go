package mcserver

import (
	"fmt"
	"log"
	"net"

	"github.com/m-nny/goinit/pkg/datatypes"
	"github.com/m-nny/goinit/pkg/mcnet"
	"github.com/m-nny/goinit/pkg/packets"
)

type Server struct {
	listener net.Listener
	conns    []*conn
	router   *mcnet.Router
}

func NewServer() *Server {
	return &Server{
		router: getRootRouter(),
	}
}

func getRootRouter() *mcnet.Router {
	router := mcnet.NewRouter()
	router.AddRoute(datatypes.STATE_HANDSHAKING, packets.PACKET_ID_HANDSHAKE, packets.HandshakeHandler)
	// TODO: add routes
	return router
}

func (s *Server) Start(host string, port uint) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	log.Printf("Server started at %s:%d", host, port)

	for {
		rw, err := listener.Accept()
		if err != nil {
			log.Printf("%+v", err)
			return err
		}
		conn := s.newConn(rw)
		s.conns = append(s.conns, conn)
		go conn.Serve()
	}
}

func (s *Server) Close() {
	for _, client := range s.conns {
		if err := client.Close(); err != nil {
			log.Printf("err closing client: %v", err)
			continue
		}
	}
}

func (s *Server) newConn(rwc net.Conn) *conn {
	return &conn{
		rwc:    rwc,
		state:  datatypes.STATE_HANDSHAKING,
		router: s.router,
	}
}
