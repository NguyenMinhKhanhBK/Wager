CREATE TABLE if NOT EXISTS purchase (
    id bigint unsigned not null auto_increment primary key,
    wager_id bigint unsigned not null, 
    buying_price decimal not null,
    bought_at bigint not null,
    foreign key (wager_id) references wagers (id)
)
