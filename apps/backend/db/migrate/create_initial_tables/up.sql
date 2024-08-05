create table if not exists user (
    id integer primary key autoincrement,
    username text not null unique,
    password_hash text not null,
    token text unique,
    email text unique,
    token_created_at datetime,
    created_at datetime default current_timestamp,
    last_login datetime
);

create unique index if not exists user_token_index on user(token);

create table if not exists video (
    youtube_id text primary key,
    title text not null,
    channel_title text not null,
    description text not null,
    published_at datetime not null,
    youtube_link text not null,
    duration_seconds integer not null,
    view_count integer not null
);

create unique index if not exists video_youtube_id_index on video(youtube_id);

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
    video_youtube_id integer not null,
    added_at datetime default current_timestamp,
    foreign key (cart_id) references cart(id),
    foreign key (video_youtube_id) references video(youtube_id),
    unique(cart_id, video_youtube_id)
);

-- TODO: Add index on cart ID

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
    video_youtube_id text not null,
    -- rank of search result for maintain order is inferred from 
    -- insertion order
    foreign key (search_id) references search(id),
    foreign key (video_youtube_id) references video(youtube_id),
    unique (search_id, video_youtube_id)
);
