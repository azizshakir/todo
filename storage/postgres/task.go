package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	pb"github.com/azizshakir/todo/genproto"
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
	err := r.db.QueryRow(querycreate,in.)
}
func (r *taskRepo) Get(id int64) (pb.Task, error)
func (r *taskRepo) List(pb.ListReq) (pb.ListResp,)
func (r *taskRepo) Update(pb.Task) (pb.Task, error)
func (r *taskRepo) Delete(int64) error
func (r *taskRepo) ListOverdue(pb.OverReq) (pb.ListResp,error)










func (r *taskRepo) Create(task pb.Task) (pb.Task, error) {
	var id int64
	err := r.db.QueryRow(`
        INSERT INTO users(first_name, last_name)
        VALUES ($1,$2) returning id`, task.FirstName, task.LastName).Scan(&id)
	if err != nil {
		return pb.Task{}, err
	}

	task, err = r.Get(id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Get(id int64) (pb.Task, error) {
	var user pb.Task
	err := r.db.QueryRow(`
        SELECT id, first_name, last_name FROM users
        WHERE id=$1`, id).Scan(&user.Id, &user.FirstName, &user.LastName)
	if err != nil {
		return pb.TAsk{}, err
	}

	return user, nil
}

func (r *taskRepo) List(page, limit int64) ([]*pb.Task, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT id, first_name, last_name FROM users LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		tasks []*pb.Task
		task  pb.Task
		count int64
	)
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.FirstName, &task.LastName)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, &task)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM users`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

func (r *taskRepo) Update(task pb.Task) (pb.Task, error) {
	result, err := r.db.Exec(`UPDATE users SET first_name=$1, last_name=$2 WHERE id=$3`,
		task.FirstName, task.LastName, task.Id)
	if err != nil {
		return pb.Task{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Task{}, sql.ErrNoRows
	}

	task, err = r.Get(task.Id)
	if err != nil {
		return pb.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
