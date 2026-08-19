# Hoard

A read-later app that lets you one-click save links from anywhere (browser, phone) and browse them in a card-feed reading experience.

## Tech Stack

- **Frontend**: Vite + React + TypeScript + Tailwind CSS
- **Backend**: Go (gin + pgx + golang-jwt)
- **Database**: Supabase PostgreSQL (with Auth)
- **Chrome Extension**: MV3 (no popup, one-click save)
- **iOS**: Phase 2 (WKWebView + Share Extension)

## Project Structure

```
hoard/
├── backend/             # Go API server
│   ├── cmd/server/      # Entry point
│   ├── internal/
│   │   ├── config/      # Environment config
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # JWT auth middleware
│   │   ├── model/       # Data models
│   │   ├── repository/  # Database access
│   │   └── service/     # Business logic
│   ├── Dockerfile       # Container build
│   └── railway.toml     # Railway config
├── frontend/            # React SPA
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── contexts/    # React contexts
│   │   ├── lib/         # API client, mock data
│   │   └── pages/       # Page components
│   └── vercel.json      # Vercel config
├── extension/           # Chrome Extension (MV3)
├── migrations/          # SQL migrations
├── fly.toml             # Fly.io config
└── .env.example         # Environment variables template
```

## Quick Start (Local Development)

### Prerequisites

- Go 1.25+
- Node.js 22+
- Supabase account (free tier)

### 1. Database Setup

1. Create a Supabase project at [supabase.com](https://supabase.com)
2. Go to **SQL Editor**, paste and run `migrations/001_init.sql`
3. Copy your project URL, anon key, and JWT secret from **Settings > API**

### 2. Environment Variables

```bash
cp .env.example .env
# Fill in your Supabase credentials
```

### 3. Backend

```bash
cd backend
GOCACHE=/tmp/go-cache go mod download
GOCACHE=/tmp/go-cache go run ./cmd/server
```

### 4. Frontend

```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:5173

## Deployment

### Frontend → Vercel

```bash
cd frontend
npm run build
vercel --prod
```

Environment variables to set on Vercel:
- `VITE_SUPABASE_URL` — your Supabase project URL
- `VITE_SUPABASE_ANON_KEY` — your Supabase anon key

### Backend → Fly.io

```bash
cd backend
fly launch --dockerfile Dockerfile
fly secrets set DATABASE_URL=postgres://... SUPABASE_JWT_SECRET=... SUPABASE_URL=https://...
fly deploy
```

### Backend → Railway

1. Create a new Railway project from `backend/`
2. Railway auto-detects the Dockerfile
3. Set environment variables:
   - `DATABASE_URL`
   - `SUPABASE_JWT_SECRET`
   - `SUPABASE_URL`
   - `PORT=8080`

### Chrome Extension

```bash
cd extension
./build.sh
```

Load `hoard-extension.zip` or the unpacked `extension/` folder in Chrome → Extensions → Developer mode → Load unpacked.

## Tickets

See [tickets/000-map.md](tickets/000-map.md) for the full roadmap.

## License

MIT
