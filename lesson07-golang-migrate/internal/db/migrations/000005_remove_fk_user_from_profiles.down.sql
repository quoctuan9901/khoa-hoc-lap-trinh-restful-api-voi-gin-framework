alter table profiles
add column user_id int unique not null;

alter table profiles
add constraint fk_user foreign key (user_id) references users(user_id) on delete cascade;