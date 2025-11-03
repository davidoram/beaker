-- +migrate Up

create table inventory (
    product_sku text not null primary key,
    stock_level int  not null,

    -- Ensure stock_level is never negative
    constraint inventory_stock_level_nonnegative
        check (stock_level >= 0),

    -- Ensure product_sku only contains lowercase alphanumeric characters, hyphens and underscores
    constraint inventory_product_sku_format
        check (product_sku ~ '^[a-z0-9_-]{1,64}$')
);


-- +migrate Down

drop table inventory;