alter table contributions
    drop column images;
create table images
(
    key             text primary key,
    contribution_id integer not null references contributions (id),
    title           text
);