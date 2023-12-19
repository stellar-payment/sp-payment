alter table accounts alter column account_no TYPE bytea using account_no::bytea;
alter table accounts add column account_no_hash bytea not null default ''::bytea;
alter table accounts add column row_hash bytea not null default ''::bytea;