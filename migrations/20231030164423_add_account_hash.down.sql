alter table accounts alter column account_no TYPE varchar(15) not null;
alter table accounts drop column account_no_hash;
alter table accounts drop column row_hash;