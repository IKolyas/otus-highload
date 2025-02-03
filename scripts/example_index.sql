select count(*) from users;

select min_val, max_val from pg_settings where name='max_connections';

show max_connections;

drop index if exists idx_first_name;
CREATE INDEX idx_first_name ON users (first_name);
drop index if exists idx_second_name;
CREATE INDEX idx_second_name ON users (second_name);
explain analyse select * from users where first_name LIKE 'Ni%' and second_name LIKE 'Pe%' order by id;


drop index if exists idx_first_name_second_name;
CREATE INDEX idx_first_name_second_name ON users (first_name text_pattern_ops, second_name text_pattern_ops);

explain analyse select * from users where first_name LIKE 'Sa%' and second_name LIKE 'Si%' order by id;