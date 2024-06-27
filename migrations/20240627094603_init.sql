-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tariffs;
DROP TABLE bookings;
-- +goose StatementEnd
