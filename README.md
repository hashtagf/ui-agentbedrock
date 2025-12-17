# 🤖 AgentBedrock UI

A beautiful and easy-to-use Chat UI for testing AWS AgentBedrock Team Agents. An alternative to the complex AWS Console.

![UI Preview](https://via.placeholder.com/800x400?text=AgentBedrock+UI)

## ✨ Features

- 💬 **ChatGPT-style Interface** - Familiar and intuitive chat experience
- 🌊 **Real-time Streaming** - See responses as they're generated (HTTP SSE)
- 🧠 **AI Thinking Status** - Know what the AI is processing
- 🔄 **Agent Call Tracking** - See which agents are being invoked step-by-step
- 📊 **Trace Viewer** - Detailed execution traces for debugging
- ⚠️ **Error Display** - Clear Lambda error messages and stack traces
- 💾 **Session Management** - Save and continue conversations
- 🌓 **Theme Support** - Dark, Light, and System themes

## 🛠 Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Nuxt 4 + Vue 3 |
| Styling | TailwindCSS |
| Backend | Golang + Gin |
| Database | MongoDB |
| AWS | AgentBedrock SDK |
| Streaming | Server-Sent Events (SSE) |
| Deployment | Docker Compose |

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- AWS Account with AgentBedrock setup
- Agent ID and Agent Alias ID

### 1. Clone and Configure

```bash
# Clone the repository
git clone <your-repo-url>
cd ui-agentbedrock

# Copy environment file
cp .env.example .env

# Edit .env with your AWS credentials
nano .env
```

### 2. Configure Environment Variables

```env
# Required: AWS AgentBedrock
AGENT_ID=your-agent-id
AGENT_ALIAS=your-agent-alias-id
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

### 3. Start with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### 4. Access the UI

Open [http://localhost:3000](http://localhost:3000) in your browser.

## 📁 Project Structure

```
ui-agentbedrock/
├── frontend/                 # Nuxt 4 Application
│   ├── app/
│   │   ├── components/       # Vue components
│   │   ├── composables/      # Vue composables
│   │   └── pages/            # Page components
│   ├── nuxt.config.ts
│   └── Dockerfile
├── backend/                  # Golang Gin API
│   ├── cmd/server/           # Main entry point
│   ├── internal/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── services/         # Business logic
│   │   ├── repository/       # Data access
│   │   └── models/           # Data models
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

## 🔧 Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev
```

## 📡 API Endpoints

### Sessions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sessions` | List all sessions |
| POST | `/api/sessions` | Create new session |
| GET | `/api/sessions/:id` | Get session with messages |
| PUT | `/api/sessions/:id` | Update session title |
| DELETE | `/api/sessions/:id` | Delete session |

### Chat

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/chat/stream` | Send message (SSE streaming) |

### SSE Events

```typescript
// Event types from /api/chat/stream
event: thinking    // AI is processing
event: agent_step  // Agent invocation step
event: content     // Response chunk
event: trace       // Execution trace
event: error       // Error occurred
event: done        // Stream complete
```

## 🎨 Themes

The UI supports three theme modes:

- **Light** ☀️ - Clean white interface
- **Dark** 🌙 - Easy on the eyes (default)
- **System** 💻 - Follows OS preference

Toggle themes using the button in the header.

## 🐛 Troubleshooting

### Connection Refused

```bash
# Check if services are running
docker-compose ps

# Restart services
docker-compose restart
```

### AWS Credentials Error

Ensure your `.env` file has valid AWS credentials with permissions for:
- `bedrock-agent-runtime:InvokeAgent`

### MongoDB Connection Error

```bash
# Check MongoDB logs
docker-compose logs mongodb

# Reset MongoDB data
docker-compose down -v
docker-compose up -d
```

## 📄 License

MIT License - feel free to use this for your projects!

---

Built with ❤️ for easier AgentBedrock testing

