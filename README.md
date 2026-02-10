# 🎤 NOTEUP - Voice Notes Evolved

> *Speak your thoughts. AI writes them down. You get the credit.* ✨

```
  ███╗   ██╗ ██████╗ ████████╗███████╗██╗   ██╗██████╗ 
  ████╗  ██║██╔═══██╗╚══██╔══╝██╔════╝██║   ██║██╔══██╗
  ██╔██╗ ██║██║   ██║   ██║   █████╗  ██║   ██║██████╔╝
  ██║╚██╗██║██║   ██║   ██║   ██╔══╝  ██║   ██║██╔═══╝ 
  ██║ ╚████║╚██████╔╝   ██║   ███████╗╚██████╔╝██║     
  ╚═╝  ╚═══╝ ╚═════╝    ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     
  
  🎙️ Voice → 🧠 AI Processing → 📝 Beautiful Notes
  Powered by Groq • Built in Go • Scaled on Kafka
```

---

## ✨ What's Noteup? (The Vibe Check)

Ever had an amazing idea at 3 AM but couldn't be bothered to type it out? **That's where Noteup comes in.** 

Just hit record. Speak naturally. Our AI transcribes, structures, and beautifies your rambling thoughts into polished, professional notes. No more voice memos that nobody reads. No more lost ideas. Just pure thought-to-text magic. ✨

Think of it as having a personal note-taking assistant who's *always* listening, *never* judging, and actually knows how to write.

---

## 🏗️ Architecture (How the Magic Works)

```
┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  📱 CLIENT (Mobile/Web)                                     │
│  ├─ Records audio                                            │
│  ├─ Uploads to Imagekit                                     │
│  └─ Gets note ID (processing status)                        │
│                                                               │
│  ↓                                                            │
│                                                               │
│  🔄 GO BACKEND (Port 3000)                                  │
│  ├─ Validates user (Firebase JWT)                           │
│  ├─ Deducts coins (⚙️ Rate limiting)                         │
│  ├─ Stores note (PostgreSQL)                                │
│  └─ Publishes to Kafka "note.created"                       │
│                                                               │
│  ↓                                                            │
│                                                               │
│  🎙️ KAFKA EVENT QUEUE                                       │
│  └─ "note.created" → Audio URL + Note ID                    │
│                                                               │
│  ↓                                                            │
│                                                               │
│  🧠 PYTHON AI SERVICE (Groq Consumer)                       │
│  ├─ 🎵 Whisper: Transcribe audio (fast)                     │
│  ├─ 🤖 Llama 3.1: Polish transcript → Beautiful notes       │
│  └─ 📤 Publishes to Kafka "aiwork.done"                     │
│                                                               │
│  ↓                                                            │
│                                                               │
│  🔄 GO BACKEND (AI Result Consumer)                         │
│  ├─ Receives title + transcript                              │
│  ├─ Updates note status → "completed"                       │
│  └─ Invalidates cache (Redis)                               │
│                                                               │
│  ↓                                                            │
│                                                               │
│  📬 CLIENT (Real-time update via polling/webhooks)          │
│  └─ "Your note is ready!" ✅                                │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack (The Weapons)

| Component | Tech | Why | ⚡ |
|-----------|------|-----|-----|
| **Backend** | Go 1.25 + Chi Router | Concurrency monster, JSON handling, REST API | ✅ |
| **AI Processing** | Python 3.11 + Groq | Whisper-large-v3 (transcribe) + Llama 3.1 (rewrite) | ✨ |
| **Message Queue** | Apache Kafka 7.7 | Async tasks, event sourcing, decoupled services | 🚀 |
| **Database** | PostgreSQL 17 | Full-text search, TSVECTOR, reliable & fast | 🔒 |
| **Cache Layer** | Redis 8.4 | Ultra-fast note retrieval, user data | ⚡ |
| **Auth** | Firebase Auth | Zero-friction Google login + JWT tokens | 🔐 |
| **File Storage** | Imagekit | Image upload, auto-optimization | 📸 |
| **Container** | Docker + Compose | One-command setup, production-ready | 🐳 |

---

## 🚀 Quick Start (5 Minutes to Greatness)

### Prerequisites
- Docker + Docker Compose (easiest way)
- OR: Go 1.25.3, Python 3.11, PostgreSQL 17, Redis 8.4, Kafka 7.7

### Option 1: Docker Compose (Recommended 🎯)

```bash
# Clone the repo
git clone <repo-url>
cd noteup

# Copy env files
cp go-service/.env.example go-service/.env
cp ai-service/.env.example ai-service/.env

# Add your secrets to .env files:
# - GROQ_API_KEY=your_groq_key
# - FIREBASE_URL=path/to/firebase.json
# - IMAGEKIT_PRIVATE_KEY=xxx

# Start everything (Kafka, Postgres, Redis, Go service, Python AI)
docker-compose up -d

# Check services
docker-compose ps

# View logs
docker-compose logs -f go_backend
docker-compose logs -f python_ai
```



### Option 2: Local Development

```bash
# 1. Start Kafka, Postgres, Redis
docker-compose up -d kafka postgres redis

# 2. Set up database
export DB_DEV_URL="postgresql://postgres:postgres@localhost:5432/notedbv2?sslmode=disable"
migrate -path go-service/migrations -database "$DB_DEV_URL" up

# 3. Start Go backend
cd go-service
go run cmd/main.go

# 4. Start Python AI (new terminal)
cd ai-service
python main.py

# 5. Test it
curl http://localhost:3000/api/v1/auth/
# Should return: {"message":"Health ok","statusCode":200,"data":null}
```

---

## 📚 API Endpoints (The Good Stuff)

### 🔐 Authentication (No Auth Needed)

```http
POST /api/v1/auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure123"
}
```

```http
POST /api/v1/auth/google
Content-Type: application/json

{
  "idToken": "firebase_id_token_here"
}
```

### 📝 Notes (🔒 Protected - Requires Bearer Token)

```http
# Create note (upload audio)
POST /api/v1/note/audio
Authorization: Bearer ACCESS_TOKEN
Content-Type: multipart/form-data

audio: <audio.mp3>
duration_seconds: 45
```

**Response (202 Accepted):**
```json
{
  "message": "processing started",
  "statusCode": 202,
  "data": {
    "note_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing"
  }
}
```

```http
# Get all your notes
GET /api/v1/note/
Authorization: Bearer ACCESS_TOKEN
```

```http
# Get specific note
GET /api/v1/note/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer ACCESS_TOKEN
```

```http
# Update note (edit title/transcript)
PUT /api/v1/note/{noteId}
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json

{
  "title": "New title",
  "transcript": "Edited content"
}
```

```http
# Delete note
DELETE /api/v1/note/{noteId}
Authorization: Bearer ACCESS_TOKEN
```

### 🔍 Search (Public)

```http
GET /api/v1/search/?q=golang
```

### 💰 Coins (Monetization)

```http
# Get coin packs
GET /api/v1/coin_packs/

# Get user balance
GET /api/v1/user/coins
Authorization: Bearer ACCESS_TOKEN
```

### ✅ Tasks (Bonus Feature)

```http
# Create task
POST /api/v1/tasks/
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json

{
  "title": "Learn Rust",
  "description": "Complete the book",
  "priority": "high",
  "due_at": "2026-03-01T00:00:00Z"
}

# Get all tasks
GET /api/v1/tasks/
Authorization: Bearer ACCESS_TOKEN

# Update task
PUT /api/v1/tasks/{taskId}
Authorization: Bearer ACCESS_TOKEN

# Mark complete
PATCH /api/v1/tasks/{taskId}/complete
Authorization: Bearer ACCESS_TOKEN

# Delete task
DELETE /api/v1/tasks/{taskId}
Authorization: Bearer ACCESS_TOKEN
```

---

## 🧠 How AI Processing Works

### The Flow

1. **User uploads audio** → Note created with status `"processing"`
2. **Kafka event published** → `note.created` topic
3. **Python AI consumer** receives message:
   - 🎵 **Whisper Transcription** (Groq): Converts audio → raw transcript
   - 🤖 **Llama 3.1 Rewrite** (Groq): Polishes into beautiful notes
   - Publishes `aiwork.done` event with title + transcript
4. **Go backend consumer** updates note:
   - Status → `"completed"`
   - Stores title + transcript + word count
   - Invalidates Redis cache
5. **Client polls** → Gets completed note with full text ✅

### Why Async?

- Audio processing can take **5-30 seconds** depending on file size
- Client gets **immediate response** (202 Accepted) with note ID
- Uses **Kafka** to decouple services (Go doesn't wait for AI)
- **Redis cache** for instant retrieval once processed
- **Fire and forget** = responsive UI 🚀

---

## 💰 Coin System (The Economics)

### How It Works

```
📊 Coin Balance = Audio Processing Credits

Seconds Per Coin: 30 seconds
Max Balance: Depends on plan

🔴 FREE PLAN
├─ 2 coins (60 seconds of audio)
└─ Refresh daily? TBD

🟢 PRO PLAN
├─ 200+ coins (6000+ seconds)
├─ Unlimited monthly
└─ Costs ₹199/month (via RevenueCat)
```

### Coin Deduction

```
When user uploads audio:
1. Calculate: coins_needed = (duration_seconds + 29) / 30
2. Check balance >= coins_needed
3. If NO → Return 402 Payment Required (buy more coins!)
4. If YES → Deduct coins + Process audio
5. Log transaction in coin_transactions table
```

### Purchase Flow

```
GET /api/v1/coin_packs/     ← Available packs
{
  "id": "pack-1",
  "coin_value": 120,
  "price": 19900,           ← ₹199 (paise)
  "popular": true
}

Tap "Buy" → RevenueCat → Webhook → Add coins
```

---

## 🗄️ Database Schema (Quick Ref)

### Users
```sql
id (UUID)          → Firebase UID
email (TEXT)       → Unique
google_id (TEXT)   → For OAuth
password (TEXT)    → Hashed (bcrypt)
picture (TEXT)     → Profile pic URL
plan (TEXT)        → "free" or "pro"
is_active (BOOL)   → Can log in?
created_at         → Auto timestamp
```

### Notes
```sql
id (UUID)                      → Note ID
user_id (UUID)                 → Owner
audio_url (TEXT)               → Imagekit URL
audio_duration_seconds (INT)   → 30, 60, 120...
audio_file_size_mb (INT)       → File size
transcript (TEXT)              → AI-generated
title (TEXT)                   → AI-generated
word_count (INT)               → Character count
status (TEXT)                  → pending|processing|completed
search_vector (TSVECTOR)       → Full-text search index
```

### UserCoin
```sql
user_id (UUID)    → Foreign key
balance (INT)     → Current coins
updated_at        → Last transaction
```

### Tasks
```sql
id (UUID)
user_id (UUID)
title (TEXT)
description (TEXT)
status (TEXT)     → pending|completed|cancelled
priority (TEXT)   → low|medium|high
due_at (TIMESTAMP)
completed_at (TIMESTAMP)  → Auto-set when status=completed
```

---

## ⚙️ Configuration (Env Variables)

### Go Service (.env)

```bash
# Database
DB_DEV_URL=postgresql://postgres:postgres@localhost:5432/notedbv2?sslmode=disable
DB_URL=postgresql://user:pass@prod-db:5432/notedbv2  # Production

# Redis
REDIS_HOST=localhost:6379
REDIS_TTL=10  # Cache TTL in seconds

# Firebase
FIREBASE_URL=./internals/config/noteon_back.json
FIREBASE_PROJECT_ID=notev2-34951

# Kafka
KAFKA_BOOTSTRAP_SERVERS=kafka:9092

# JWT Tokens
ACCESS_TOKEN_KEY=your-secret-key-min-32-chars
REFRESH_TOKEN_KEY=your-secret-key-min-32-chars

# Server
PORT=:3000
APP_ENV=dev  # or "prod"
```

### Python Service (.env)

```bash
# Groq API (free tier = 30 req/min)
GROQ_API_KEY=gsk_xxxxx

# Kafka
KAFKA_BOOTSTRAP_SERVERS=kafka:9092

# Logging
LOG_LEVEL=INFO
```

---

## 🧪 Testing (Make Sure It Works)

### Health Check

```bash
curl http://localhost:3000/api/v1/auth/
# Should return: {"message":"Health ok","statusCode":200,"data":null}
```

### Full Flow Test

```bash
# 1. Sign up
curl -X POST http://localhost:3000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Response:
{
  "message": "Successfully authenticated",
  "statusCode": 201,
  "data": {
    "tokens": {
      "access_token": "eyJ0...",
      "refresh_token": "eyJ0..."
    },
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com",
    "plan": "free",
    "created_at": "2026-02-10T12:00:00Z"
  }
}

# 2. Get user coins
TOKEN="eyJ0..."
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/user/coins

# Response:
{
  "message": "success",
  "statusCode": 200,
  "data": {
    "balance": 2,
    "seconds_per_coin": 30,
    "max_seconds": 60
  }
}

# 3. Upload audio
curl -X POST http://localhost:3000/api/v1/note/audio \
  -H "Authorization: Bearer $TOKEN" \
  -F "audio=@sample.mp3" \
  -F "duration_seconds=45"

# Response (202 = Processing):
{
  "message": "processing started",
  "statusCode": 202,
  "data": {
    "note_id": "550e8400-e29b-41d4-a716-446655440001",
    "status": "processing"
  }
}

# 4. Poll for result (wait 5-10 seconds)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/note/550e8400-e29b-41d4-a716-446655440001

# Response (when ready):
{
  "message": "Fetched",
  "statusCode": 200,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Meeting Notes from Q1 Planning",
    "transcript": "During our Q1 planning session, we discussed...",
    "status": "completed",
    "word_count": 342,
    "created_at": "2026-02-10T12:00:00Z"
  }
}
```

---

## 📊 Performance (How Fast Is It?)

| Operation | Time | Notes |
|-----------|------|-------|
| Audio upload | ~500ms | Imagekit upload |
| Transcription (60s) | 3-5s | Groq Whisper |
| Rewrite to notes (60s) | 2-4s | Groq Llama 3.1 |
| **Total end-to-end** | **6-10s** | Async → client gets 202 immediately |
| Note retrieval (cached) | **<50ms** | Redis FTW |
| Database query (indexed) | ~100-200ms | PostgreSQL TSVECTOR index |
| Coin deduction | ~50ms | Single UPDATE |

---

## 🔐 Security (Don't Worry, We Got You)

✅ **JWT Token Validation**
- Firebase ID tokens verified server-side
- Access tokens expire in 30 mins
- Refresh tokens expire in 7 days
- Always use HTTPS in production

✅ **Database Security**
- Passwords hashed with bcrypt (cost=10)
- SQL injection prevented with parameterized queries (pgx)
- Row-level security with user_id checks

✅ **Rate Limiting**
- Coin system prevents audio spam
- Each user limited by their balance
- Free tier = 2 coins (60 seconds/day)

✅ **File Validation**
- Max audio file: 25MB
- Max free audio: 10MB
- Max pro audio: 80MB
- Viruses scanned by Imagekit

---

## 🚀 Deployment (Ship It!)

### Docker (Production)

```bash
# Build images
docker-compose build

# Run with production env
docker-compose -f docker-compose.yml up -d

# Check health
docker-compose ps
docker-compose logs -f go_backend
```

### Environment for Production

```bash
# .env (production)
APP_ENV=prod
DB_URL=postgresql://user:password@prod-db.example.com/notedbv2
REDIS_HOST=redis-prod.example.com:6379
KAFKA_BOOTSTRAP_SERVERS=kafka-broker-1:9092,kafka-broker-2:9092,kafka-broker-3:9092
ACCESS_TOKEN_KEY=<64-char-random-string>
REFRESH_TOKEN_KEY=<64-char-random-string>
PORT=:3000
```

### Deployment Checklist

- [ ] Postgres backups enabled (daily)
- [ ] Redis memory limit set (`maxmemory-policy allkeys-lru`)
- [ ] Kafka replicas ≥ 3
- [ ] SSL/TLS enabled for all services
- [ ] Monitoring set up (Prometheus, Grafana)
- [ ] Logs centralized (ELK stack or similar)
- [ ] Error tracking (Sentry)
- [ ] Rate limiting configured
- [ ] CORS configured for frontend domain
- [ ] Load balancer in front of Go service

---

## 🐛 Troubleshooting

### "Kafka connection failed"
```bash
✅ Check: docker-compose ps | grep kafka
✅ Logs: docker-compose logs kafka
✅ Reset: docker-compose down -v && docker-compose up
```

### "Database connection refused"
```bash
✅ Postgres running? docker-compose ps postgres
✅ Migrations run? migrate -path migrations -database $DB_URL version
✅ Credentials correct in .env?
```

### "AI processing stuck"
```bash
✅ Python service running? docker-compose logs python_ai
✅ Groq API key valid? (check in .env)
✅ Kafka topic exists? kafka-topics --list --bootstrap-server kafka:9092
```

### "Notes not appearing in list"
```bash
✅ Status = "completed"? (Query db directly)
✅ Cache key deleted? redis-cli DEL user:notes:v1:*
✅ Note belongs to user? Check user_id in database
```

---

## 📈 What's Next? (Roadmap)

- [ ] **Webhooks** - Real-time note completion notifications
- [ ] **OCR** - Extract text from images + notes
- [ ] **Folders** - Organize notes by topic
- [ ] **Collaboration** - Share notes with team
- [ ] **Export** - Download as PDF/Markdown
- [ ] **Mobile Apps** - iOS + Android native apps
- [ ] **Offline mode** - Record without internet
- [ ] **Voice commands** - "Send to Slack", "Create task"
- [ ] **Multiple languages** - Beyond English
- [ ] **Analytics** - Track writing patterns

---

## 💪 Built With

- **Go** - The systems language (fast AF)
- **Python** - AI glue (Groq magic)
- **PostgreSQL** - The database god
- **Redis** - Cache that doesn't lie
- **Kafka** - Async = happy users
- **Docker** - Containers 🐳
- **Firebase** - Auth done right
- **Groq** - Free AI (whisper + llama)
- **Imagekit** - Image CDN
- **Love** - Lots of it ❤️

---



```
╔══════════════════════════════════════════════════════╗
║                                                      ║
║  Made with ❤️  by farhan khan              ║
║  "Your thoughts deserve better than voice memos"    ║
║                                                      ║
║  🎙️ Speak → 🧠 AI → 📝 Notes         ║
║                                                      ║
╚══════════════════════════════════════════════════════╝

Ready to take note? Let's go. 🚀
```
