create table if not exists app_user (
    id text primary key, -- uuid string
    username text not null unique,
    password_hash text not null,
    token text unique,
    token_created_at timestamp,
    created_at timestamp default current_timestamp,
    last_login timestamp
);

create unique index if not exists app_user_token_index on app_user(token);

create table if not exists video (
    youtube_id text primary key,
    title text not null,
    channel_title text not null,
    description text not null,
    published_at timestamp not null,
    youtube_link text not null,
    duration_seconds integer not null,
    view_count integer not null
);

create unique index if not exists video_youtube_id_index on video(youtube_id);

create table if not exists cart_item (
    id text primary key, -- uuid string
    app_user_id text not null,
    video_youtube_id text not null,
    added_at timestamp default current_timestamp,
    foreign key (app_user_id) references app_user(id) on delete cascade,
    foreign key (video_youtube_id) references video(youtube_id) on delete cascade,
    unique(app_user_id, video_youtube_id)
);

create index if not exists cart_item_app_user_id on cart_item(app_user_id);

create table if not exists search (
    id text primary key, -- uuid string
    query text not null,
    app_user_id text not null,
    executed_at timestamp default current_timestamp,
    foreign key (app_user_id) references app_user(id) on delete cascade
);

create index if not exists search_query on search(query);

create table if not exists search_result (
    id text primary key,
    search_id text not null,
    video_youtube_id text not null,
    -- rank of search result for maintain order is inferred from 
    -- insertion order
    foreign key (search_id) references search(id) on delete cascade,
    foreign key (video_youtube_id) references video(youtube_id) on delete cascade,
    unique (search_id, video_youtube_id)
);

create table if not exists job (
    id text primary key, -- uuid string status text not null,
    response_payload text, -- a json string
    created_at timestamp default current_timestamp
);
