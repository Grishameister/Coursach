package poolDb

import (
	"bytes"
	"compress/gzip"
	images "github.com/Grishameister/Coursach/internal/images"
	"log"
)

type Pool struct {
	in chan []byte
	uc images.IUsecase
}

func NewPool(in chan []byte, uc images.IUsecase) *Pool {
	return &Pool{
		in: in,
		uc: uc,
	}
}

func (p *Pool) ZipToDb() {
	for im := range p.in {
		var zipbuf bytes.Buffer
		gz := gzip.NewWriter(&zipbuf)

		if _, err := gz.Write(im); err != nil {
			log.Println(err)
			return
		}
		if err := gz.Close(); err != nil {
			log.Println(err)
			return
		}
		if err := p.uc.PostFrame(zipbuf.Bytes()); err != nil {
			log.Println(err)
			return
		}
		log.Println("SUCCESS")
	}
}

func (p *Pool) Run(counter int) {
	for i := 0; i < counter; i++ {
		go p.ZipToDb()
	}
}
