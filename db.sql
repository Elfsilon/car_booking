-- Active: 1719329026480@@127.0.0.1@5432@booking
DROP TABLE tariffs;

CREATE TABLE tariffs(
  id serial PRIMARY KEY,
  name UNIQUE text,
  price numeric(13, 4)
);

INSERT INTO tariffs(name, price) VALUES ('basic', 1000.0);
INSERT INTO tariffs(name, price) VALUES ('premium', 800.0);

CREATE TABLE bookings(
  id serial PRIMARY KEY,
  user_id uuid,
  car_id uuid,
  from_date date,
  to_date date
);

SELECT * FROM tariffs;