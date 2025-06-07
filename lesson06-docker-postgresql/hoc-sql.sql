-- Create a new database
create database hoc_golang;

-- Drop a database
drop database testing;

-- Create a new schema
create schema school;

-- Drop a schema
drop schema school cascade;

-- One - One
-- Create users table
create table if not  exists users (
	user_id serial primary key,
	name varchar(50) not null,
	email varchar(100) unique not null
);

-- Create profiles table
create table if not exists profiles (
	profile_id serial primary key,
	user_id int unique not null,
	phone varchar(10),
	address varchar(100),
	constraint fk_user foreign key (user_id) references users(user_id) on delete cascade
);

-- Drop table (dangerours)
drop table if exists profiles;
drop table if exists users;

-- One - Many
-- Create categories table
create table if not exists categories (
	category_id serial primary key,
	name varchar(50) not null
);

-- Create products table
create table if not exists products (
	product_id serial primary key,
	category_id int not null,
	name varchar(100) not null,
	price int not null check (price > 0),
	image varchar(255),
	status int not null check (status in (1,2)),
	constraint fk_category foreign key (category_id) references categories (category_id) on delete restrict
);

-- Drop table (dangerours)
drop table if exists products;
drop table if exists categories;


-- Many - Many
-- Create students table
create table if not exists students (
	student_id serial primary key,
	name varchar(50) not null
);

-- Create courses table
create table if not exists courses (
	course_id serial primary key,
	name varchar(50) not null
);

-- Create students_courses table
create table if not exists students_courses (
	student_id int not null,
	course_id int not null,
	primary key (student_id, course_id),
	constraint fk_student foreign key (student_id) references students(student_id) on delete cascade,
	constraint fk_course foreign key (course_id) references courses(course_id) on delete cascade
)

-- Drop table (dangerours)
drop table if exists students_courses;
drop table if exists courses;
drop table if exists students;


------------------------- Các cấu truy vấn (SQL) hay sử dụng -----------------------------------
-- Thêm dữ liệu: INSERT INTO table (col1, col2) VALUES (val1, val2)
insert into users (name, email) values ('Vu Quoc Tuan', 'contact.quoctuan@gmail.com');
insert into users (name, email) values ('Toney Teo', 'contact.teo@gmail.com');
insert into users (name, email) values ('Le Van Tung', 'contact.tung@gmail.com');

insert into profiles (user_id, phone, address) values (1, '0901234567', '123 Cach mang thang 8');

insert into categories (name) values ('Dien thoai'), ('Laptop');

insert into products (category_id, name, price, image, status) values
(3, 'iPhone 18 Pro Max', 10, 'images/iphone-18-pro-max.jpg', 1),
(4, 'iPhone 17 Pro Max', 30000000, 'images/iphone-17-pro-max.jpg', 1);

-- Cập nhật dữ liệu: UPDATE table SET col1 = value1, col2 = val2 WHERE condition
update users set email = 'tuan@quoctuan.com', name = 'Mr.Tuan' where user_id = 1;
update profiles set phone = '0906784312' where user_id = 2;

-- Xóa dữ liệu: DELETE FROM table WHERE condition
delete from users where user_id = 3;
delete from products;
delete from categories where category_id = 1;

-- Lấy dữ liệu: SELECT * FROM table WHERE condition ORDER BY col [DESC/ASC] LIMIT ... OFFSET ...
select * from products;
select name, price from products;
select count(*) as total_rows from products;

select * from products where price >= 400000 and price <= 1000000;

select * from products order by price desc;
select * from products order by price asc;

select * from products limit 3 offset 4;

select name, price from products 
where price >= 400000 and price <= 30000000
order by price desc
limit 3;

select category_id, count(*) from products
group by category_id
having count(*) > 2;