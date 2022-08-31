package rpc

import (
	"context"

	pb "github.com/cut4cut/hezzl-test-work/internal/controller/rpc/proto"
	"github.com/cut4cut/hezzl-test-work/internal/entity"
	"github.com/cut4cut/hezzl-test-work/internal/usecase"
	"github.com/cut4cut/hezzl-test-work/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	pb.UnimplementedServiceUserServer
	u usecase.UserRepoInterface
	l logger.Interface
}

func (s *UserServer) Create(ctx context.Context, in *pb.UserName) (*pb.UserId, error) {
	s.l.Info("call create method")

	var (
		id  int64 = -1
		err error
	)

	if id, err = s.u.Create(ctx, in.GetName()); err != nil {
		return &pb.UserId{Id: id}, err
	}

	return &pb.UserId{Id: id}, nil
}

func (s *UserServer) Delete(ctx context.Context, in *pb.UserId) (*pb.UserId, error) {
	s.l.Info("call delete method")

	var (
		id  int64 = -1
		err error
	)

	if id, err = s.u.Delete(ctx, in.GetId()); err != nil {
		return &pb.UserId{Id: id}, err
	}

	return &pb.UserId{Id: id}, nil
}

func (s *UserServer) GetList(ctx context.Context, in *pb.Pagination) (*pb.UserList, error) {
	s.l.Info("call getList method for page=%v", in.GetPage())

	pag := entity.NewPagination(
		uint64(in.GetPage()), in.GetDescName(), in.GetDescCreated(),
	)

	users, err := s.u.GetList(ctx, pag)
	if err != nil {
		return &pb.UserList{Users: []*pb.User{}}, nil
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{Id: user.Id, CreatedDt: timestamppb.New(user.CreatedDt), Name: user.Name}
	}

	return &pb.UserList{Users: pbUsers}, nil
}

func (s *UserServer) mustEmbedUnimplementedServiceUserServer() {}
