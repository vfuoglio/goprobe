# 🚀 Concurrent Web Core-Metrics Ping Engine

An asynchronous CLI network utility engineered in Go. The application concurrently audits target server endpoints, parses performance latency metrics down to the millisecond, and handles network telemetry streaming using native concurrency primitives.

## 🛠️ Architectural Design & Engineering Choices

- **True Concurrency via Goroutines:** Rather than blocking execution on synchronous HTTP network operations, this engine abstracts each URL check into its own decoupled `Goroutine` green-thread, maximizing system CPU core usage.
- **Safe Thread Synchronization:** Utilizes `sync.WaitGroup` to orchestrate deterministic lifecycles for spawned processes, avoiding dangling goroutines and memory leaks.
- **Asynchronous Channel Streams:** Leverages a thread-safe Go `chan` to handle multi-producer, single-consumer communication, allowing the main engine thread to dynamically aggregate telemetry as network tasks resolve.
- **Zero-Dependency Core:** Crafted entirely utilizing Go's highly optimized standard library (`net/http`, `sync`, `time`), showcasing an understanding of language foundations without external library bloating.

## 🧳 Tech Stack & Market Skills
- **Language:** Go (Golang)
- **Concurrency Concepts:** Goroutines, Channels, Data Race Avoidance, Worker Synchronization (`sync.WaitGroup`)
- **Network Engineering:** HTTP Client Tuning, Precise Latency Calculations via `time.Duration`, I/O Streaming

## 📦 Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/vfuoglio/goprobe.git 
   cd go-ping-engine
   ```

2. **Verify your local Go environment configuration:**
   ```bash
   go version
   ```

## 💻 Usage

1. Populate your target endpoints inside the `urls.txt` file (one URL per line).
2. Execute the engine compiler directly:
   ```bash
   go run main.go
   ```

To compile a native standalone distribution binary executable:
```bash
go build -o ping-engine main.go
./ping-engine
```

