package server

import (
	"context"

	"cskyzn.com/pkg/bimgserver/rpc"
	"github.com/h2non/bimg"
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (s *Server) Thumbnail(ctx context.Context, req *rpc.ThumbnailReq) (*rpc.ContentResp, error) {
	var (
		img  = bimg.NewImage(req.Content)
		resp rpc.ContentResp
		err  error
	)

	if resp.Content, err = img.Thumbnail(int(req.Pixels)); err != nil {
		return nil, err
	}

	return &resp, nil
}
