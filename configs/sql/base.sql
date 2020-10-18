create database images with owner postgres encoding 'utf8' LC_COLLATE = 'ru_RU.UTF-8'LC_CTYPE = 'ru_RU.UTF-8' TABLESPACE = pg_default;

create table frames(
	id bigserial primary key,
	bytes bytea not null,
	reg_date timestamp
);