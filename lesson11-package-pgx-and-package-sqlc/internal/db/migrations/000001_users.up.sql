create extension if not exists "pgcrypto";

create table if not exists users (
	user_id serial primary key,
	uuid uuid not null default gen_random_uuid(),
	name varchar(50) not null,
	email varchar(100) unique not null,
	created_at timestamptz not null default now()
);