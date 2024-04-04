alter table config_mapping add column timeout int not null default 86400 comment 'Time out of template in seconds (default 24h as 86400 seconds)';
