-- Active: 1719329026480@@127.0.0.1@5432@booking
DROP TABLE tariffs;
DROP TABLE bookings;

CREATE TABLE tariffs(
  id serial PRIMARY KEY,
  name text UNIQUE,
  price numeric(13, 4)
);

CREATE TABLE bookings(
  id serial PRIMARY KEY,
  user_id uuid,
  car_id uuid,
  from_date date,
  to_date date
);

INSERT INTO tariffs(name, price) VALUES ('basic', 1000.0);
INSERT INTO tariffs(name, price) VALUES ('premium', 800.0);

INSERT INTO bookings(user_id, car_id, from_date, to_date) 
VALUES ('c84fde82-6679-4384-af11-7406de3d3e14', 'c84fde82-6679-4384-af11-7406de3d3e14', '2024-06-01', '2024-06-05');
INSERT INTO bookings(user_id, car_id, from_date, to_date) 
VALUES ('c84fde82-6679-4384-af11-7406de3d3e14', 'c84fde82-6679-4384-af11-7406de3d3e14', '2024-06-09', '2024-06-12');
INSERT INTO bookings(user_id, car_id, from_date, to_date) 
VALUES ('c84fde82-6679-4384-af11-7406de3d3e14', 'c84fde82-6679-4384-af11-7406de3d3e14', '2024-06-19', '2024-06-29');

SELECT * FROM tariffs;
SELECT * FROM bookings;

-- Checks range intersections
SELECT count(*)
FROM bookings
WHERE car_id = 'c84fde82-6679-4384-af11-7406de3d3e14' AND (
  '2024-05-28' >= from_date - 3::integer AND '2024-05-28' <= to_date + 3::integer OR 
  '2024-05-28' >= from_date - 3::integer AND '2024-05-28' <= to_date + 3::integer
);

SELECT to_date, to_date + 3::integer FROM bookings;

DELETE FROM bookings WHERE id = 1;