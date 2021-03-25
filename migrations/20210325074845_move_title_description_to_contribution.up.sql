alter table articles
    drop column title;
alter table articles
    drop column description;
alter table contributions
    add column title text not null default 'default title';
alter table contributions
    add column description text;