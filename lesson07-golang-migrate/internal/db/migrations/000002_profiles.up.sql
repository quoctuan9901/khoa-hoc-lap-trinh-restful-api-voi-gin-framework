create table if not exists profiles (
	profile_id serial primary key,
	user_id int unique not null,
	phone varchar(10),
	address varchar(100),
	constraint fk_user foreign key (user_id) references users(user_id) on delete cascade
);