
-- +migrate Up

create table inventory (
    product_sku varchar(50) not null primary key,
    stock_level int not null,

    -- Ensure stock_level is never negative
    check (stock_level >= 0),

    -- Ensure product_sku only contains lowercase alphanumeric characters, hyphens and underscores
    check (product_sku ~ '^[a-z0-9_-]+$')
);


-- +migrate Down

drop table inventory;