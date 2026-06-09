# TravelSphere

A destination discovery and trip planner built with Go and Beego v2. Explore countries, view attractions, and manage a personal travel wishlist — no database required.

## Features

- Country explorer with SSR rendering and AJAX search/filter
- Destination detail pages with OpenTripMap attractions
- In-memory wishlist CRUD via JSON API
- Auth middleware (header-based `Username`)
- Logging filter on all requests
- Unit tests across services, utils, and middleware layers

## Tech Stack

- **Go 1.22+** · **Beego v2.1.0**
- REST Countries API · OpenTripMap API
- In-memory store with `sync.RWMutex`

## Prerequisites

- Go 1.22 or higher
- Beego CLI (`bee`) is required
- An OpenTripMap API key

## Setup

**1. Clone the repository**

```bash
git clone https://github.com/200215-Moynul-Islam/TravelSphere-Moynul.git
cd TravelSphere-Moynul
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Configure the app**

```bash
cp conf/app.conf.example conf/app.conf
```

Edit `conf/app.conf` and set your API key:

```ini
opentripmap_api_key = your_api_key_here
```

**4. Run the server**

```bash
bee run
```

The app starts at `http://localhost:8080`.

## Configuration

| Key                       | Description                                |
| ------------------------- | ------------------------------------------ |
| `opentripmap_api_key`     | Required. Your OpenTripMap API key         |
| `restcountries_base_url`  | REST Countries base URL (has default)      |
| `opentripmap_geoname_url` | OpenTripMap geoname endpoint (has default) |
| `opentripmap_radius_url`  | OpenTripMap radius endpoint (has default)  |
| `httpport`                | Server port (default: `8080`)              |

## Routes

### SSR Pages

| Method | Route              | Description                                |
| ------ | ------------------ | ------------------------------------------ |
| GET    | `/`                | Home — featured countries, search          |
| GET    | `/countries`       | Country explorer with search/filter        |
| GET    | `/countries/:code` | Destination detail (e.g. `/countries/USA`) |
| GET    | `/wishlist`        | Wishlist page _(not yet implemented)_      |
| GET    | `/dashboard`       | Dashboard page _(not yet implemented)_     |

> The `:code` slug uses the **CCA3 country code** (e.g. `USA`, `BGD`, `FRA`).

### JSON API

| Method | Route                    | Description                                      |
| ------ | ------------------------ | ------------------------------------------------ |
| GET    | `/api/countries`         | Country list; supports `?search=` and `?region=` |
| GET    | `/api/wishlist`          | Get wishlist entries                             |
| POST   | `/api/wishlist`          | Add wishlist entry                               |
| PUT    | `/api/wishlist/:id`      | Update note/status                               |
| DELETE | `/api/wishlist/:id`      | Delete entry                                     |
| GET    | `/api/dashboard/summary` | Wishlist stats                                   |

Wishlist and dashboard API routes are protected — include a `Username` header:

```
Username: your_username
```

## Sample API Requests

**Get all countries**

```bash
curl http://localhost:8080/api/countries
```

**Search countries**

```bash
curl "http://localhost:8080/api/countries?search=france&region=europe"
```

**Add to wishlist**

```bash
curl -X POST http://localhost:8080/api/wishlist \
  -H "Username: moynul" \
  -H "Content-Type: application/json" \
  -d '{"country_name": "Japan", "note": "Cherry blossom season", "status": "Planned"}'
```

**Get wishlist**

```bash
curl http://localhost:8080/api/wishlist \
  -H "Username: moynul"
```

**Update wishlist entry**

```bash
curl -X PUT http://localhost:8080/api/wishlist/{id} \
  -H "Username: moynul" \
  -H "Content-Type: application/json" \
  -d '{"note": "Updated note", "status": "Visited"}'
```

**Delete wishlist entry**

```bash
curl -X DELETE http://localhost:8080/api/wishlist/{id} \
  -H "Username: moynul"
```

**Get dashboard summary**

```bash
curl http://localhost:8080/api/dashboard/summary \
  -H "Username: moynul"
```

## Wishlist Storage

Uses an **in-memory store** (`map[string][]WishlistEntry`) protected by `sync.RWMutex`. Data is scoped per username and does not persist across server restarts.

Wishlist entry fields:

| Field          | Type      | Notes                  |
| -------------- | --------- | ---------------------- |
| `id`           | string    | UUID, auto-generated   |
| `country_name` | string    | Required               |
| `note`         | string    | Optional               |
| `status`       | string    | `Planned` or `Visited` |
| `created_at`   | time.Time | Auto-set               |

## Running Tests

```bash
go test ./...
```

With coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Project Structure

```
TravelSphere/
├── conf/               # app.conf and app.conf.example
├── constants/          # Nav items, status values, field lists
├── data/               # In-memory WishlistStore and RWMutex
├── controllers/        # SSR controllers + api/ sub-package
├── middlewares/        # Auth middleware
├── models/             # Country, WishlistEntry, Attraction
├── routers/            # SSR and API route registration
├── services/           # Business logic layer
├── static/             # CSS and JS assets
├── utils/              # API clients, DTOs, mappers, validation
└── views/              # Beego templates (layouts, partials, pages)
```

## Known Limitations

- Wishlist (`/wishlist`) and Dashboard (`/dashboard`) SSR pages are not yet implemented
- Add-to-wishlist button on destination page is wired but AJAX handler is pending
- Login/logout UI is present in the header partial but not functional
- Data resets on server restart (in-memory only)
