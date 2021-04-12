package tcpConnPool

import (
	"errors"
	"github.com/Grishameister/Coursach/configs/config"
	"github.com/enriquebris/goconcurrentqueue"
	"log"
	"net"
)

type TcpPool struct {
	q      *goconcurrentqueue.FixedFIFO
	isInit bool
}

func NewPool(cap int) *TcpPool {
	return &TcpPool{
		q:      goconcurrentqueue.NewFixedFIFO(100),
		isInit: false,
	}
}

func InitPool(size int) (*TcpPool, error) {
	if size > 100 {
		return &TcpPool{}, errors.New("too many size of connections")
	}
	p := NewPool(size)
	for i := 0; i < size; i++ {
		conn, err := net.Dial("tcp", config.Conf.Stats.Server.Address+":"+config.Conf.Stats.Server.Port)
		if err != nil {
			log.Println(err.Error())
			return p, err
		}

		if err := p.Push(conn); err != nil {
			log.Println(err.Error())
			return p, err
		}
	}

	p.isInit = true
	return p, nil
}

func (p *TcpPool) Push(conn net.Conn) error {
	return p.q.Enqueue(conn)
}

func (p *TcpPool) Pop() (net.Conn, error) {
	conn, err := p.q.DequeueOrWaitForNextElement()
	if err != nil {
		return nil, err
	}
	return conn.(net.Conn), nil
}
