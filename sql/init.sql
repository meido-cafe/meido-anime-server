create table if not exists video
(
    id            integer
    constraint table_name_pk
    primary key autoincrement,
    bangumi_id    integer default 0 not null,
    origin        integer default 0 not null,
    category      integer default 0 not null,
    title         text    default '' not null,
    season        integer default 0 not null,
    cover         text    default '' not null,
    total         integer default 0 not null,
    rss_url       text    default '' not null,
    play_time     integer default 0 not null,
    source_dir    text    default '' not null,
    create_time   integer default 0 not null,
    update_time   integer default 0 not null
);


create table if not exists video_item
(
    id          integer
    constraint video_item_pk
    primary key autoincrement,
    pid         integer default 0 not null,
    episode     integer default 0 not null,
    source_path text    default '' not null,
    link_path   text    default '' not null,
    create_time integer default 0 not null,
    update_time integer default 0 not null
);