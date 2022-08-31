package rpc

import (
	pb "github.com/cut4cut/hezzl-test-work/internal/controller/rpc/proto"
	"github.com/cut4cut/hezzl-test-work/internal/usecase"
	"github.com/cut4cut/hezzl-test-work/pkg/logger"
	"google.golang.org/grpc"
)

func Register(s *grpc.Server, u usecase.UserRepoInterface, l logger.Interface) {
	userServ := UserServer{u: u, l: l}
	pb.RegisterServiceUserServer(s, &userServ)
}
