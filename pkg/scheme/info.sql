create table users (
                       id       uuid primary key default gen_random_uuid(),
                       username text not null,
                       password text not null,
                       created_at  timestamptz not null default current_timestamp,
                       updated_at  timestamptz,
                       deleted_at  timestamptz
);

create table tokens (
                        id    uuid primary key default gen_random_uuid(),
                        token text not null
);

create table accounts (
                          id      uuid primary key default gen_random_uuid(),
                          number  text        not null,
                          user_id uuid        not null references users on delete cascade,
                          balance decimal     not null default 0.0,
                          created_at timestamptz not null default current_timestamp,
                          updated_at timestamptz,
                          deleted_at timestamptz
);

create table transactions (
                              id         uuid primary key default gen_random_uuid(),
                              account_id uuid not null references accounts on delete cascade,
                              type       text not null,
                              amount     decimal not null default 0.0,
                              created_at    timestamptz not null default current_timestamp,
                              updated_at    timestamptz,
                              deleted_at    timestamptz
);