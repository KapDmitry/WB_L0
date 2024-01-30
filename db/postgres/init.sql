

\c orders;

DROP TABLE IF EXISTS orderslist;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS items;

CREATE TABLE orderslist (
    id SERIAL PRIMARY KEY,
    order_uid text unique, 
    track_number text NOT NULL,
    order_entry text NOT NULL,
    locale text NOT NULL,
    internal_signature text NOT NULL,
    customer_id text NOT NULL,
    delivery_service text NOT NULL,
    shardkey text NOT NULL,
    sm_id bigint NOT NULL,
    date_created text NOT NULL,
    oof_shard text NULL
);


CREATE TABLE delivery (
    id SERIAL PRIMARY KEY, 
    order_id text NOT NULL unique,
    delivery_name text NOT NULL,
    phone text NOT NULL,
    zip text NOT NULL,
    city text NOT NULL,
    delivery_address text NOT NULL,
    region text NOT NULL,
    email text NOT NULL,
    CONSTRAINT fk_delivery FOREIGN KEY(order_id)
        REFERENCES orderslist(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    order_id text NOT NULL unique,
    payment_transaction text NOT NULL unique,
    request_id text NOT NULL unique,
    currency text NOT NULL,
    payment_provider text NOT NULL,
    amount integer NOT NULL,
    payment_dt bigint NOT NULL,
    bank text NOT NULL,
    delivery_cost integer NOT NULL,
    goods_total integer NOT NULL,
    custom_fee integer NOT NULL,
     CONSTRAINT fk_payment FOREIGN KEY(order_id)
        REFERENCES orderslist(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY, 
    order_id text NOT NULL, 
    chrt_id bigint NOT NULL unique,
    track_number text NOT NULL,
    price integer NOT NULL,
    rid text NOT NULL,
    name_item text NOT NULL,
    sale integer NOT NULL,
    size_item text NOT NULL,
    total_price integer NOT NULL,
    nm_id bigint NOT NULL unique,
    brand text NOT NULL,
    status_item integer NOT NULL, 
    CONSTRAINT fk_items FOREIGN KEY(order_id)
        REFERENCES orderslist(order_uid)
        ON DELETE CASCADE
);

