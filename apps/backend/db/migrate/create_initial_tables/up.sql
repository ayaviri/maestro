create table if not exists user (
    id integer primary key autoincrement,
    username text not null unique,
    password_hash text not null,
    email text unique,
    created_at datetime default current_timestamp,
    last_login datetime
);

create table if not exists video (
    id integer primary key autoincrement,
    youtube_id text not null unique,
    title text not null,
    channel_title text not null,
    description text not null,
    published_at datetime not null,
    youtube_link text not null,
    duration_seconds integer not null,
    view_count integer not null
);

create table if not exists cart (
    id integer primary key autoincrement,
    user_id integer not null unique,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp,
    foreign key (user_id) references user(id)
);

create table if not exists cart_item (
    id integer primary key autoincrement,
    cart_id integer not null,
    video_id integer not null,
    added_at datetime default current_timestamp,
    foreign key (cart_id) references cart(id),
    foreign key (video_id) references video(id),
    unique(cart_id, video_id)
);

create table if not exists search (
    id integer primary key autoincrement,
    query text not null,
    user_id integer not null,
    executed_at datetime default current_timestamp,
    foreign key (user_id) references user(id)
);

create table if not exists search_result (
    id integer primary key autoincrement,
    search_id integer not null,
    video_id integer not null,
    -- rank of search result for maintain order is inferred from 
    -- insertion order
    foreign key (search_id) references search(id),
    foreign key (video_id) references video(id),
    unique (search_id, video_id)
);
