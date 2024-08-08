create table if not exists job (
    id integer primary key autoincrement,
    status text not null,
    created_at datetime default current_timestamp
);
