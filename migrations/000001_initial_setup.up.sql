CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL Primary Key,
    assignee VARCHAR(64),
    title VARCHAR(64),
    summary VARCHAR(128),
    deadline timestamp not null, 
    status VARCHAR(32)
);
