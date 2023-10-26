create table customers (
    id uuid primary key,
    user_id uuid not null,
    legal_name varchar(255) not null,
    phone varchar(50),
    email varchar(255),
    birthdate varchar(255),
    address varchar(255),
    photo_profile varchar(255),
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create table merchants (
    id uuid primary key,
    user_id uuid not null,
    name varchar(255),
    address varchar(255),
    phone varchar(50),
    email varchar(255),
    pic_name varchar(255),
    pic_phone varchar(50),
    pic_email varchar(255),
    photo_profile varchar(255),
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create table accounts (
    id uuid primary key,
    owner_id uuid not null,
    account_type int not null,
    balance decimal(18, 2) not null,
    account_no varchar(15) not null,
    pin varchar(255) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create table transactions (
    id bigint primary key,
    account_id uuid not null,
    recipient_id uuid not null,
    trx_type int not null,
    trx_datetime timestamp with time zone not null,
    trx_status smallint not null,
    trx_fee decimal(18, 2) not null,
    nominal decimal(18, 2) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);

create table transaction_types (
    id serial primary key,
    name varchar(255) not null,
    trx_fee_rate decimal(6, 3) not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone
);