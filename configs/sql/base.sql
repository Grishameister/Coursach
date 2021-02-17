create database images with owner postgres encoding 'utf8' LC_COLLATE = 'ru_RU.UTF-8'LC_CTYPE = 'ru_RU.UTF-8' TABLESPACE = pg_default;

CREATE EXTENSION btree_gist;

create table frames(
	id bigserial primary key,
	bytes bytea not null,
	reg_date timestamp
);

CREATE INDEX ON frames USING gist(reg_date);
