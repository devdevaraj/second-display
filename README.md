# Virtual Display Management System

A production-grade system to manage and stream headless Wayland displays. This system allows you to create virtual displays using Sway, capture them using `wf-recorder` and `ffmpeg`, and stream them seamlessly to your browser using **WebRTC**.

## Architecture

```text
                        +-----------------------------------------------+
                        |                 Go Backend                    |
+--------------+        |  +-------------+    +----------------------+  |
| Next.js App  | <--API--> | HTTP Server | -> | Session Manager      |  |
| (Dashboard & |        |  +-------------+    +----------------------+  |
|  WebRTC      |        |                           |                   |
|  Player)     | <----------------+                 | (Process Supv)    |
+--------------+       (WebRTC)   |                 v                   |
                                  |  +-------------------------------+  |
                                  |  | swaymsg create_output HEADLESS|  |
                                  |  +-------------------------------+  |
                                  |  | wf-recorder (capture wlroots) |  |
                                  |  +-------------------------------+  |
                                  |  | ffmpeg (encode to h264 RTP)   |  |
                                  |  +-------------------------------+  |
                                  |                 |                   |
                                  +--<-- Pion WebRTC Handler <-<--<--+  |
                                  +-------------------------------------+
```

- **Domain/Service Layer**: Clean Go architecture, handles state and orchestrates Linux processes.
- **Sway Integration**: Dynamically adds/removes HEADLESS outputs and configures their resolution.
- **Capture Pipeline**: `wf-recorder` captures the frame buffer and pipes raw video to `ffmpeg`, which encodes it with `libx264` (ultrafast) into an RTP stream aimed at localhost.
- **Streaming Pipeline**: Pion WebRTC receives localhost RTP packets on isolated ports per session and forwards them cleanly to connected Next.js WebRTC clients. No ICE Trickling complex external dependencies required!

## System Prerequisites

This project is built for **Debian 13 (Wayland)**. It strictly depends on `wf-recorder` and `sway` (wlroots).

1. `setup.sh` is provided in the repository root to install exact dependencies.
2. Go 1.22+
3. Node.js 18+

## Quick Start

### 1. Install Dependencies

```bash
./setup.sh
```

### 2. Start Go Backend

```bash
cd backend
go run ./cmd/vdisplay
```

*The backend will run on `http://localhost:8080`*

### 3. Start Next.js Frontend

```bash
cd frontend
npm run dev
```

*The frontend will run on `http://localhost:3000`*

## Troubleshooting

- **No Output/Black Screen**: Make sure you are running the backend *inside* the active Sway session (so the `WAYLAND_DISPLAY` and `SWAYSOCK` env vars are available to `exec.Command`).
- **WebRTC Connection Failed**: Ensure ICE candidates can resolve (usually not an issue locally), and confirm `ffmpeg` is definitely emitting the local RTP stream.
- **Zombie processes**: The Go supervisor (`capture.Pipeline`) cleans up `ffmpeg` and `wf-recorder` gracefully using context cancellations. Double-check using `ps aux | grep wf-recorder` if terminated improperly.

## Scaling Strategy

If deploying in a multi-user "Mini Cloud Desktop" environment:

1. **Containerization**: Do NOT run sway on bare-metal for isolated users. Run a headless `sway` instance inside a Docker container using `xwayland-run` or similar wrappers.
2. **GPU Encoding**: Replace `libx264` with `h264_vaapi` or `h264_nvenc` in the `-c:v` arguments dynamically if underlying hardware is available, heavily dropping CPU load per session.
3. **WebRTC SFU**: If multiple users need to watch the *same* virtual display, the current PeerConnection 1:1 map should be upgraded using something like LiveKit or a selective forwarding scheme via Pion.
