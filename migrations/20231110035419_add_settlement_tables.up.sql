create table beneficiaries (
    id bigint not null,
    merchant_id uuid not null,
    amount decimal(18, 2) not null default 0,
    withdrawal_date timestamp with time zone,
    status smallint not null default 0,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    deleted_at timestamp with time zone
);

create table settlements (
    id bigint not null,
    transaction_id bigint not null,
    merchant_id uuid not null,
    beneficiary_id bigint not null,
    amount decimal(18, 2) not null,
    settlement_date timestamp with time zone not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    deleted_at timestamp with time zone
);