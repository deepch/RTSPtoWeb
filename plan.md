# Plan to Replace RTSPtoWeb Frontend with React + Vite

## 1. Project Overview & Backend Analysis

The current project, **RTSPtoWeb**, is a high-performance RTSP stream converter written in Go. It consumes RTSP streams and restreams them via:
- **HLS (HTTP Live Streaming)**
- **MSE (Media Source Extensions)** via WebSockets
- **WebRTC**

### Backend Architecture
- **Server**: Uses `gin-gonic` for the HTTP server.
- **API**: Provides a RESTful JSON API for managing streams and channels.
- **Authentication**: Supports Basic Auth and a custom cookie-based session (`RTSP_SESSION`).
- **CORS**: The backend already includes a `CrossOrigin` middleware, allowing requests from different origins (crucial for our separate React dev server).
- **Configuration**: Uses `config.json` for server and stream settings.

### How to use the Backend
We will use the existing Go application as a headless API server.
1.  **Run the Go server**: `go run *.go` (or use the binary/Docker). It listens on port `:8083` by default.
2.  **API Endpoints**: The React frontend will consume endpoints like:
    -   `POST /login`: For authentication.
    -   `GET /streams`: To list available streams.
    -   `POST /stream/:uuid/add`: To create streams.
    -   `GET /stream/:uuid/channel/:channel/webrtc`: To initiate WebRTC playback.
    -   (And others documented in `docs/api.md`)

## 2. Frontend Implementation Plan

We will build a modern Single Page Application (SPA) using **React** and **Vite**.

### Technology Stack
-   **Framework**: React 18+
-   **Build Tool**: Vite
-   **Styling**: Tailwind CSS
-   **UI Library**: Shadcn UI (for a premium, polished look)
-   **Icons**: Lucide React
-   **Routing**: React Router DOM
-   **State Management**: React Context or Zustand (for auth and stream state)
-   **Video Players**:
    -   **HLS**: `hls.js`
    -   **WebRTC**: Native `RTCPeerConnection` (custom hook/component)
    -   **MSE**: Native `MediaSource` + WebSocket (custom hook/component)

### Key Features & Components

1.  **Authentication**:
    -   Login Page with username/password.
    -   Auth Context to handle session persistence (checking `RTSP_SESSION` cookie or local state).
    -   Role-based access control (Admin vs. Viewer).

2.  **Dashboard (Stream List)**:
    -   Grid view of available streams.
    -   Live thumbnails (if possible) or status indicators (Online/Offline).
    -   **Admin Controls**: Add, Edit, Delete buttons (visible only to Admins).

3.  **Stream Player Component**:
    -   A robust `Player` component that accepts a `streamId`, `channelId`, and `protocol` (WebRTC/MSE/HLS).
    -   **WebRTC Logic**: Ported from `RtspToWeb.js`. Handles signaling via the `/webrtc` endpoint.
    -   **MSE Logic**: Ported from `RtspToWeb.js`. Connects to WebSocket and feeds `SourceBuffer`.
    -   **HLS Logic**: Uses `hls.js` to play the `.m3u8` stream.

4.  **Stream Management (Admin)**:
    -   Forms to Add/Edit streams (Name, RTSP URL, On-demand toggle, etc.).
    -   Validation for RTSP URLs.

## 3. Implementation Steps

1.  **Initialize Project**: Create `frontend` folder with Vite.
2.  **Setup Tailwind & Shadcn**: Configure the design system.
3.  **Develop API Client**: Create an `api.ts` helper to handle requests to `http://localhost:8083`.
4.  **Build Layouts**: Create a main layout with a sidebar/navbar.
5.  **Implement Login**: Connect to `/login` API.
6.  **Implement Dashboard**: Fetch and display streams from `/streams`.
7.  **Build Player**: Implement the complex streaming logic (WebRTC/MSE/HLS).
8.  **Add Admin Features**: Create forms for stream management.

---

## 4. Prompt for Gemini

Use the following prompt to generate the project structure and code.

```markdown
I want to replace the frontend of an existing Go project (RTSPtoWeb) with a new React application using Vite.
The Go backend runs on `http://localhost:8083` and exposes a JSON API.

**Project Requirements:**
1.  **Tech Stack**: React, Vite, TypeScript, Tailwind CSS, Shadcn UI, React Router DOM, Lucide React.
2.  **Folder Structure**: Create a `frontend` directory in the root.
3.  **Authentication**:
    -   Create a Login page.
    -   Use the `POST /login` endpoint (accepts JSON `{username, password}`).
    -   Handle the `RTSP_SESSION` cookie.
    -   Implement an `AuthProvider` to manage user state and roles ('admin' vs 'user').
4.  **Dashboard**:
    -   Fetch streams from `GET /streams`.
    -   Display streams in a responsive grid.
    -   Show status (online/offline).
5.  **Streaming Player**:
    -   Create a `Player` component that supports **WebRTC**, **MSE**, and **HLS**.
    -   **WebRTC**: Implement the logic to create an `RTCPeerConnection`, send the offer to `POST /stream/:uuid/channel/:channel/webrtc`, and set the remote answer. Use a STUN server (e.g., `stun:stun.l.google.com:19302`).
    -   **MSE**: Connect to WebSocket at `ws://localhost:8083/stream/:uuid/channel/:channel/mse` and feed the `MediaSource`.
    -   **HLS**: Use `hls.js` to play `http://localhost:8083/stream/:uuid/channel/:channel/hls/live/index.m3u8`.
6.  **Admin Features**:
    -   Allow Admins to **Add**, **Edit**, and **Delete** streams via the API.
    -   Regular users should only be able to view streams.
7.  **UI/UX**:
    -   Use a dark, premium theme with Shadcn UI components.
    -   Responsive sidebar navigation.

**Specific API Details**:
-   **List Streams**: `GET /streams` -> returns JSON object with stream details.
-   **Add Stream**: `POST /stream/:uuid/add` -> body `{name, url, on_demand, ...}`.
-   **WebRTC Signaling**: POST SDP offer to `/stream/:uuid/channel/:channel/webrtc` -> returns SDP answer (base64 encoded).

Please generate the initial project setup, including:
-   `vite.config.ts` (with proxy to localhost:8083 to avoid CORS issues if needed, though backend supports it).
-   `src/api/client.ts` (Axios or Fetch wrapper).
-   `src/components/Player.tsx` (The core streaming logic).
-   `src/pages/Dashboard.tsx`.
-   `src/pages/Login.tsx`.
-   `src/App.tsx` with routing.
```
