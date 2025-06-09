alter table users
rename column email to user_email;

alter table users
alter column user_email set data type text;