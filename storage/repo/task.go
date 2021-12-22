package repo

import (
	pb "github.com/azizshakir/todo/genproto"
)

// TaskStorageI ...
type TaskStorageI interface {
	Create(pb.Task) (pb.Task, error)
	Get(id string) (pb.Task, error)
	List(pb.ListReq) (pb.ListResp, error)
	Update(pb.Task) (pb.Task, error)
	Delete(string) error
	ListOverdue(pb.OverReq) (pb.ListResp, error)
}
