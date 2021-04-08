drop table system_data;
create table system_data
(
    key        varchar(255) primary key,
    value      text not null,
    type       varchar(100),
    updated_at timestamptz
);
insert into system_data
values ('term_and_condition', 'this is example Term and Condition', 'document', now());