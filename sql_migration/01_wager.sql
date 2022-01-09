create table if not exists wagers (
    id bigint unsigned not null auto_increment primary key,
    total_wager_value int unsigned not null,
    odds int unsigned not null,
    selling_percentage int unsigned not null,
    selling_price decimal not null,
    current_selling_price decimal not null,
    percentage_sold int unsigned,
    amount_sold decimal,
    place_at datetime not null
)
