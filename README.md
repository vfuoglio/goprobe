# 🚀 goprobe — Concurrent Network Telemetry Engine

A high-performance, asynchronous CLI network utility engineered in Go. The application concurrently audits target server endpoints, parses performance latency metrics down to the millisecond, and handles network telemetry streaming using native concurrency primitives.

## 🛠️ Architectural Design & Engineering Choices

- **Bounded Concurrency (Worker Pool Pattern):** Rather than blindly spawning unlimited threads that risk socket exhaustion, file descriptor leaks, or firewall rate-limiting bans, the engine implements a token-based channel semaphore to strictly cap concurrent workers at `MaxConcurrency`.
- **True Concurrency via Goroutines:** Network I/O operations run inside lightweight, non-blocking green-threads (`Goroutines`), allowing the engine to maximize multi-core CPU efficiency during simultaneous probes.
- **Safe Thread Synchronization:** Utilizes `sync.WaitGroup` to orchestrate deterministic lifecycles for spawned processes, avoiding dangling background routines, while utilizing `defer` rollbacks to ensure tokens clear under network failure modes.
- **Asynchronous Channel Streams:** Leverages a thread-safe Go `chan` to handle multi-producer, single-consumer communication, allowing the main runtime thread to dynamically process telemetry live as jobs finish.
- **Zero-Dependency Core:** Crafted entirely utilizing Go's highly optimized standard library (`net/http`, `sync`, `time`), showcasing an understanding of language foundations without external library bloating.

## 🧳 Tech Stack & Market Skills
- **Language:** Go (Golang)
- **Concurrency Concepts:** Bounded Concurrency, Semaphore Primitives, Goroutines, Channels, Worker Synchronization (`sync.WaitGroup`)
- **Network Engineering:** HTTP Client Tuning, Precise Latency Calculations via `time.Duration`, I/O Streaming

## 📦 Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/vfuoglio/goprobe.git
   cd goprobe
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
go build -o goprobe main.go
./goprobe
```
# 🚀 goprobe — Concurrent Network Telemetry Engine

A high-performance, asynchronous CLI network utility engineered in Go. The application concurrently audits target server endpoints, parses performance latency metrics down to the millisecond, and handles network telemetry streaming using native concurrency primitives.

## 🛠️ Architectural Design & Engineering Choices

- **Bounded Concurrency (Worker Pool Pattern):** Rather than blindly spawning unlimited threads that risk socket exhaustion, file descriptor leaks, or firewall rate-limiting bans, the engine implements a token-based channel semaphore to strictly cap concurrent workers at `MaxConcurrency`.
- **True Concurrency via Goroutines:** Network I/O operations run inside lightweight, non-blocking green-threads (`Goroutines`), allowing the engine to maximize multi-core CPU efficiency during simultaneous probes.
- **Safe Thread Synchronization:** Utilizes `sync.WaitGroup` to orchestrate deterministic lifecycles for spawned processes, avoiding dangling background routines, while utilizing `defer` rollbacks to ensure tokens clear under network failure modes.
- **Asynchronous Channel Streams:** Leverages a thread-safe Go `chan` to handle multi-producer, single-consumer communication, allowing the main runtime thread to dynamically process telemetry live as jobs finish.
- **Zero-Dependency Core:** Crafted entirely utilizing Go's highly optimized standard library (`net/http`, `sync`, `time`), showcasing an understanding of language foundations without external library bloating.

## 🧳 Tech Stack & Market Skills
- **Language:** Go (Golang)
- **Concurrency Concepts:** Bounded Concurrency, Semaphore Primitives, Goroutines, Channels, Worker Synchronization (`sync.WaitGroup`)
- **Network Engineering:** HTTP Client Tuning, Precise Latency Calculations via `time.Duration`, I/O Streaming

## 📦 Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/vfuoglio/goprobe.git
   cd goprobe
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
go build -o goprobe main.go
./goprobe
```

