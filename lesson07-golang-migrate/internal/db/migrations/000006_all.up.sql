create table if not exists categories (
	category_id serial primary key,
	name varchar(150) not null
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