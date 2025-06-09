alter table users
rename column user_email to email;

alter table users
alter column email set data type varchar(100);