package postgres

import (
	"database/sql"
	"time"

	pb "github.com/azizshakir/todo/genproto"

	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	db *sqlx.DB
}

// NewTaskRepo ...
func NewTaskRepo(db *sqlx.DB) *taskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(in pb.Task) (pb.Task, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO tasks (assignee, title, summary, deadline, status) 
		VALUES ($1,$2,$3,$4,$5) returning id`,
		in.Assignee,
		in.Title,
		in.Summary,
		in.Deadline,
		in.Status).Scan(&id)
	if err != nil {
		return pb.Task{}, err
	}

	task, err := r.Get(id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Get(id int64) (pb.Task, error) {
	var task pb.Task

	err := r.db.QueryRow(`SELECT id, assignee, title, summary, deadline, status 
		FROM tasks WHERE id = $1`, id).Scan(
		&task.Id,
		&task.Assignee,
		&task.Title,
		&task.Summary,
		&task.Deadline,
		&task.Status,
	)
	if err != nil {
		return pb.Task{}, err
	}
	return task, nil
}

func (r *taskRepo) List(in pb.ListReq) (pb.ListResp, error) {
	ofset := (in.Page - 1) * in.Limit

	rows, err := r.db.Queryx(`SELECT id, assignee, title, summary, deadline, status 
	FROM tasks limit $1 offset $2`, in.Limit, ofset)
	if err != nil {
		return pb.ListResp{}, err
	}
	if err = rows.Err(); err != nil {
		return pb.ListResp{}, err
	}
	defer rows.Close()

	var list pb.ListResp
	for rows.Next() {
		var task pb.Task
		err = rows.Scan(
			&task.Id,
			&task.Assignee,
			&task.Title,
			&task.Summary,
			&task.Deadline,
			&task.Status,
		)
		if err != nil {
			return pb.ListResp{}, err
		}
		list.Tasks = append(list.Tasks, &task)
	}
	err = r.db.QueryRow(`SELECT count(*) FROM tasks`).Scan(&list.Count)
	if err != nil {
		return pb.ListResp{}, err
	}
	return list, nil
}

func (r *taskRepo) Update(in pb.Task) (pb.Task, error) {
	result, err := r.db.Exec(
		`UPDATE tasks set assignee=$1, title=$2, summary=$3, deadline=$4, status=$5 WHERE id=$6`,
		in.Assignee,
		in.Title,
		in.Summary,
		in.Deadline,
		in.Status,
		in.Id,
	)
	if err != nil {
		return pb.Task{}, err
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

func (r *taskRepo) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
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

	rows, err := t.db.Query(`SELECT id, assignee, title, summary, deadline, status FROM tasks WHERE deadline > $1`, duration)
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
	err = t.db.QueryRow(`SELECT count(*) FROM tasks WHERE deadline > $1`, duration).Scan(&list.Count)
	if err != nil {
		return pb.ListResp{}, nil
	}
	return list, nil
}
