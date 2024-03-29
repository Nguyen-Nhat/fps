alter table config_task add column is_async tinyint(1) not null default 0 comment 'Is async task';
