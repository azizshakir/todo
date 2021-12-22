package service

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/azizshakir/todo/genproto"
	l "github.com/azizshakir/todo/pkg/logger"
	"github.com/azizshakir/todo/storage"
)

// TaskService ...
type TaskService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewTaskService ...
func NewTaskService(db *sqlx.DB, log l.Logger) *TaskService {
	return &TaskService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *TaskService) Create(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	id, err := uuid.NewV4()

	if err != nil {
		s.logger.Error("failed to create uuid", l.Error(err))
		return &pb.Task{},err
	}
	req.Id = id.String()
	task, err := s.storage.Task().Create(*req)
	if err != nil {
		s.logger.Error("failed to create task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	return &task, nil
}

func (s *TaskService) Get(ctx context.Context, req *pb.ByIdReq) (*pb.Task, error) {
	task, err := s.storage.Task().Get(req.GetId())
	if err != nil {
		s.logger.Error("failed to get task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	return &task, nil
}

func (s *TaskService) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	list, err := s.storage.Task().List(*req)
	if err != nil {
		s.logger.Error("failed to list tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list tasks")
	}

	return &list, nil
}

func (s *TaskService) Update(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	task, err := s.storage.Task().Update(*req)
	if err != nil {
		s.logger.Error("failed to update task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	return &task, nil
}

func (s *TaskService) Delete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyResp, error) {
	err := s.storage.Task().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete task", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	return &pb.EmptyResp{}, nil
}

func (s *TaskService) ListOverdue(ctx context.Context, req *pb.OverReq) (*pb.ListResp, error) {
	list, err := s.storage.Task().ListOverdue(*req)
	if err != nil {
		s.logger.Error("failed to listOverdue tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to listOverdue tasks")
	}

	return &list, nil
}
