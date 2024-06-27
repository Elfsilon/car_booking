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
  to_date date,
  UNIQUE(car_id, from_date, to_date)
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
SELECT * FROM bookings ORDER BY to_date;

-- Checks range intersections
SELECT count(*)
FROM bookings
WHERE car_id = 'c84fde82-6679-4384-af11-7406de3d3e14' AND (
  '2024-05-28' >= from_date - 3::integer AND '2024-05-28' <= to_date + 3::integer OR 
  '2024-05-28' >= from_date - 3::integer AND '2024-05-28' <= to_date + 3::integer
);

-- Report
WITH start_end_month_dates AS (
  SELECT 
    date_trunc('month', '2024-06-21'::date)::date as month_start_date,
    (date_trunc('month', '2024-06-21'::date) + interval '1 month')::date as month_end_date
), days_in_month AS (
  SELECT EXTRACT(epoch FROM (month_end_date::timestamp - month_start_date::timestamp)) / 86400 as days_count
  FROM start_end_month_dates
), computed_booked_days_count AS (
  SELECT 
    car_id,
    CASE WHEN to_date >= month_end_date THEN month_end_date - 1::integer ELSE to_date END - 
      CASE WHEN from_date < month_start_date THEN month_start_date ELSE from_date END + 1::integer 
      as booked_days
  FROM bookings, start_end_month_dates
  WHERE from_date < month_end_date AND to_date >= month_start_date
)

SELECT car_id, ROUND(sum(booked_days) / days_count * 100) as booked_percent
FROM computed_booked_days_count, days_in_month
GROUP BY car_id, days_count;

SELECT to_date, to_date + 3::integer FROM bookings;

DELETE FROM bookings WHERE id = 1;

SELECT date_trunc('month', '2024-06-21'::date);
SELECT date_trunc('month', '2024-06-21'::date) + interval '1 month';