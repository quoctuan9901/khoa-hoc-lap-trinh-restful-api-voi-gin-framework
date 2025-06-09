alter table profiles
drop constraint fk_user;

alter table profiles
drop column user_id;