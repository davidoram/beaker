
-- +migrate Up

create table inventory (
    product_sku varchar(50) not null primary key,
    stock_level int not null
);

-- +migrate Down

drop table inventory;