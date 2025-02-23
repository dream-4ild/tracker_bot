CREATE TYPE task_status AS ENUM ('active', 'close', 'backlog');

create table tasks(
    id serial primary key,
    user_id integer,
    project text,
    task text,
    status task_status,
    deadline timestamp,
    updated_at timestamp default now()
);

create index user_index on tasks using btree (user_id);

create index status_index on tasks using btree (status);
