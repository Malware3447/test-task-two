create table currency(
    uuid UUID primary key,
    amount int
);

insert into currency (uuid, amount)
values ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 10000);