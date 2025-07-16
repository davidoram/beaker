
-- +migrate Up

create table inventory (
    product_sku varchar(50) not null primary key,
    stock_level int not null,

    -- Ensure stock_level is never negative
    check (stock_level >= 0)
);


-- +migrate Down

drop table inventory;