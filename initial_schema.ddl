create table todo (
  id serial primary key,
  text text not null default '',
  completed boolean default false
);
