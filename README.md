# Running

First of all, you need to launch docker-compose:

```bash
docker compose --env-file ./config/dev.env up --build
```

Then apply mirgations:

```bash
goose -dir ./migrations postgres "host=localhost port=5432 password=1234 user=postgres dbname=postgres sslmode=disable" up
```

Now app is available at **_localhost:3000_**

# PgAdmin

- Username: test@user.com
- Password: 1234

**_The same password is used for entering the server_**

# API Doc

## Attention

I created mock service for getting cars, there are only the next 5 cars:

```go
"c84fde82-6679-4384-af11-7406de3d3e14": {Sign: "Н314ХО123", Name: "Lada Granta", Color: "White"},
"e78b2415-b47c-435a-91e2-655ec5a08023": {Sign: "М265ДЫ123", Name: "Lada Vesta", Color: "Blue"},
"b46cbaa7-d02f-4571-ab6e-1883813715bf": {Sign: "К159ЕК93", Name: "Kia Rio", Color: "White"},
"150bd48b-1671-469a-822c-cc236d670a45": {Sign: "Е358ВА93", Name: "Mitsubishi Lancer", Color: "Red"},
"542af38b-c4b9-4d74-a7fa-e21c2a50a8cf": {Sign: "Х777ХХ123", Name: "Mercedes-Benz C-class", Color: "Black"},
```

Header **'User-ID'**, which is used for creating a report may be random uuid, I didn't do any user_id validation
as I did with cars. For instance, you can pick one of the cars' ids.

### Get cost of the range

```bash
GET /appraise?from=<date_string>&to=<date_string>

params:
  from: date_string #(YYY-MM-DD)
  to: date_string #(YYY-MM-DD)
```

Returns the total booking cost for picked range. Range can't be longer than 30 days.

### Get bookings for the car

```yaml
GET /booked_dates?car_id=<uuid>

params:
  car_id: uuid # uuid of the car, for instance, "/booked_dates?car_id=c84fde82-6679-4384-af11-7406de3d3e14"
```

Returns all the booked date ranges for the passed car. For instance, if car was booked from '20 Aug to 23 Aug', then this method will return '20 Aug to 23 Aug'.

### Get unavaliable for booking ranges for the car

```yaml
GET /unavailable_dates?car_id=<uuid>

params:
  car_id: uuid # uuid of the car, for instance, "/booked_dates?car_id=c84fde82-6679-4384-af11-7406de3d3e14"
```

This method almost the same as preceding, but it adds and removes 3 days after and before range, so it return all the dates which is unavailable for booking. It's pretty handy when frontend will be rendering calendar picker. For instance, if car was booked from '20 Aug to 23 Aug', then this method will return '17 Aug to 26 Aug'.

### Book the car

```yaml
POST /book

headers:
  User-ID: any uuid
body:
  car_id: uuid
  date_from: date_string (YYY-MM-DD)
  date_to: date_string (YYY-MM-DD)
```

Creates booking for the given range. **date_from** and **date_to** can't be weekends and pause between bookings must be at least 3 days. Given range must not be longer than 30 days.

### Unbook the car

```yaml
DELETE /book/<booking_id>

params:
  booking_id: int
```

Removes booking. It's just a handy method for making testing a bit easier.
