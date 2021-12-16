package postgres

import (
	"time"
	"database/sql"

	"github.com/jmoiron/sqlx"

	pb "github.com/azizshakir/todo/genproto"
)

type taskRepo struct {
	db *sqlx.DB
}

// NewTaskRepo ...
func NewTaskRepo(db *sqlx.DB) *taskRepo {
	return &taskRepo{db: db}
}
var (
	querycreate        = `insert into tasks (assignee,title,summary,deadline,status) values ($1,$2,$3,$4,$5) returning id`
	queryget           = `select id,assignee,title,summary,deadline,status from tasks where id = $1`
	querylist          = `select id,assignee,title,summary,deadline,status from tasks limit $1 offset $2`
	querycount         = `select count(*) from tasks`
	queryupdate        = `update tasks set assignee=$1, title=$2,summary=$3,deadline=$4,status=$5 where id=$6`
	querydel           = `delete from tasks where id = $1`
	querydeadline      = `select id,assignee,title,summary,deadline,status from tasks where deadline > $1`
	querydeadlinecount = `select count(*) from tasks where deadline > $1`
)
func (r *taskRepo) Create(in pb.Task) (pb.Task, error){
	var id int64
	err := r.db.QueryRow(
		querycreate,
		in.Assignee,
		in.Title,
		in.Summary,
		in.Deadline,
		in.Status).Scan(&id)
	if err != nil {
		return pb.Task{},err
	}

	task,err := r.Get(id)
	if err != nil {
		return pb.Task{},err
	}
	
	return task,nil
}
func (r *taskRepo) Get(id int64) (pb.Task, error){
	var task pb.Task

	err := r.db.QueryRow(queryget,id).Scan(
		&task.Id,
		&task.Assignee,
		&task.Title,
		&task.Summary,
		&task.Deadline,
		&task.Status,
	)
	if err != nil {
		return pb.Task{},err
	}
	return task,nil
}

func (r *taskRepo) List(in pb.ListReq) (pb.ListResp, error){
	ofset := (in.Page - 1) * in.Limit

	rows, err := r.db.Queryx(querylist,in.Limit,ofset)
	if err != nil {
		return pb.ListResp{},err
	}
	if err = rows.Err(); err != nil {
		return pb.ListResp{},err
	}
	defer rows.Close()
	
	var list pb.ListResp
	for rows.Next() {
		var task pb.Task
		err := rows.Scan(
			&task.Id,
			&task.Assignee,
			&task.Title,
			&task.Summary,
			&task.Deadline,
			&task.Status,
		)
		if err != nil {
			return pb.ListResp{},err
		}
		list.Tasks = append(list.Tasks, &task)
	}
	err = r.db.QueryRow(querycount).Scan(&list.Count)
	if err != nil {
		return pb.ListResp{},err
	}
	return list,nil
}
func (r *taskRepo) Update(in pb.Task) (pb.Task, error){
	result, err := r.db.Exec(
		queryupdate,
		in.Assignee,
		in.Title,
		in.Summary,
		in.Deadline,
		in.Status,
		in.Id,
	)
	if err != nil {
		return pb.Task{},err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Task{}, sql.ErrNoRows
	}

	task, err := r.Get(in.Id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil

}
func (r *taskRepo) Delete(id int64) error{
	result, err := r.db.Exec(querydel, id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (t *taskRepo) ListOverdue(in pb.OverReq) (pb.ListResp, error) {
	duration, err := time.Parse("2006-01-02", in.Time)
	if err != nil {
		return pb.ListResp{}, err
	}

	rows, err := t.db.Query(querydeadline, duration)
	if err != nil {
		return pb.ListResp{}, nil
	}
	var list pb.ListResp
	for rows.Next() {
		var task pb.Task
		err = rows.Scan(
			&task.Id,
			&task.Assignee,
			&task.Title,
			&task.Summary,
			&task.Deadline,
			&task.Status)
		if err != nil {
			return pb.ListResp{}, nil
		}
		list.Tasks = append(list.Tasks, &task)
	}
	err = t.db.QueryRow(querydeadlinecount, duration).Scan(&list.Count)
	if err != nil {
		return pb.ListResp{}, nil
	}
	return list, nil
}
