# WebSocket Integration Guide for Mobile Applications

**Target Audience**: Mobile development teams (Flutter, React Native, Native iOS/Android)  
**Backend Version**: 1.0.0  
**Last Updated**: November 12, 2025

---

## Table of Contents

- [Introduction](#introduction)
- [System Overview](#system-overview)
- [Connection Requirements](#connection-requirements)
- [Authentication Flow](#authentication-flow)
- [Message Protocol](#message-protocol)
- [Event Reference](#event-reference)
- [Integration Workflow](#integration-workflow)
- [Location Management](#location-management)
- [Notification Handling](#notification-handling)
- [Error Scenarios](#error-scenarios)
- [Best Practices](#best-practices)
- [Testing Environment](#testing-environment)
- [Appendix](#appendix)

---

## Introduction

The Risk Place Angola backend provides a real-time notification system via WebSockets that enables mobile applications to:

- **Receive instant alerts** about risks and dangers in the user's vicinity
- **Get report notifications** about incidents near the user's location
- **Update user location** for proximity-based notifications
- **Maintain persistent connections** for real-time updates
- **Handle reconnection** automatically on network changes

This document explains how the backend WebSocket system works and what mobile applications need to implement for successful integration.

---

## System Overview

### How It Works

```
┌─────────────────┐
│  Mobile App     │
│  (Flutter)      │
└────────┬────────┘
         │
         │ 1. HTTP POST /auth/login
         ▼
┌─────────────────┐
│   Backend API   │
│   Returns JWT   │
└────────┬────────┘
         │
         │ 2. WebSocket Connect
         │    ws://host/ws/alerts
         │    Header: Authorization: Bearer <JWT>
         ▼
┌─────────────────┐
│  WebSocket Hub  │
│  Validates JWT  │
└────────┬────────┘
         │
         │ 3. Register Client
         ▼
┌─────────────────┐
│  Active Session │
│  Send/Receive   │
└─────────────────┘
```

### Key Components

1. **WebSocket Hub**: Central manager for all active connections
2. **Location Store**: Redis-based geospatial index for user positions
3. **Event Dispatcher**: Triggers notifications based on domain events
4. **Push Notification Service**: FCM fallback for offline users

### Architecture Benefits

- **Real-time**: Sub-second notification delivery
- **Scalable**: Supports thousands of concurrent connections
- **Reliable**: Automatic reconnection and offline handling
- **Efficient**: Only users in proximity receive notifications

---

## Connection Requirements

### Endpoint Information

| Environment | WebSocket URL | Protocol |
|-------------|---------------|----------|
| **Development** | `ws://localhost:8000/ws/alerts` | ws:// |
| **Staging** | `ws://risk-place-angola-904a.onrender.com/ws/alerts` | ws:// |
| **Production** | `wss://example.riskplace.com/ws/alerts` | wss:// (TLS) |

> ⚠️ **Important**: Production uses `wss://` (WebSocket Secure) with TLS encryption.

### Required Headers

```
Authorization: Bearer <JWT_TOKEN>
```

The JWT token must be obtained from the login endpoint before connecting.

### Connection Parameters

- **Protocol Version**: WebSocket (RFC 6455)
- **Subprotocols**: None
- **Compression**: Not required
- **Heartbeat**: Client-initiated (recommended every 30s)

---

## Authentication Flow

### Step 1: Login to Get JWT

**HTTP Request:**
```http
POST /api/v1/auth/login HTTP/1.1
Host: example.riskplace.com
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "userPassword123"
}
```

**HTTP Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3YWE5YzBjMC0xNGQzLTRhZjktYTYzMS05NTZjMTJhMWYxMDAiLCJleHAiOjE3MzE0NTYwMDB9.abc123...",
  "expires_in": 1762442336,
  "refresh_token": "def456...",
  "token_type": "Bearer",
  "user": {
    "id": "7aa9c0c0-14d3-4af9-a631-956c12a1f100",
    "name": "Lopes Estevão",
    "email": "lopes@example.com",
    "phone": "+244923111111",
    "roles": ["citizen"]
  }
}
```

**Store the JWT token** securely on the device (e.g., secure storage, keychain).

### Step 2: Connect WebSocket with JWT

The mobile app must include the JWT token in the `Authorization` header when establishing the WebSocket connection.

**Connection Example (Conceptual):**
```javascript
// Conceptual example - adapt to your platform
const ws = new WebSocket('wss://example.riskplace.com/ws/alerts', {
  headers: {
    'Authorization': 'Bearer ' + jwtToken
  }
});
```

### Step 3: Connection Validation

The backend validates the JWT on connection:

**Success:**
- Connection upgraded to WebSocket
- Client registered in the hub
- Ready to send/receive messages

**Failure:**
- Connection rejected with HTTP 401 Unauthorized
- Client must re-authenticate

### Token Expiration

- **JWT Validity**: Typically 24 hours (check `exp` claim)
- **What Happens on Expiry**: WebSocket connection remains active
- **Recommended**: Refresh token proactively before expiration
- **On Disconnect**: Re-authenticate and reconnect with new token

---

## Message Protocol

### Message Format

All messages (sent and received) use JSON format:

```json
{
  "event": "event_name",
  "data": {
    // Event-specific payload
  }
}
```

### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event` | string | Yes | Event type identifier (see [Event Reference](#event-reference)) |
| `data` | object | Yes | Event payload (varies by event type) |

### Message Size Limits

- **Max message size**: 256 KB
- **Recommended**: Keep messages small (<10 KB)
- **Large data**: Use HTTP API endpoints instead

---

## Event Reference

### Client → Server Events

These are events that the mobile app sends to the backend.

#### 1. `update_location`

**Purpose**: Update the user's current GPS coordinates so they can receive proximity-based notifications.

**When to Send**:
- When app launches (after WebSocket connects)
- Every 30-60 seconds while app is active
- When user's location changes significantly (e.g., >100 meters)
- Before entering background (optional)

**Request Format**:
```json
{
  "event": "update_location",
  "data": {
    "latitude": -8.903290,
    "longitude": 13.312540
  }
}
```

**Fields**:
| Field | Type | Required | Range | Description |
|-------|------|----------|-------|-------------|
| `latitude` | float64 | Yes | -90 to 90 | GPS latitude coordinate |
| `longitude` | float64 | Yes | -180 to 180 | GPS longitude coordinate |

**Response**:
```json
{
  "event": "location_updated",
  "data": {
    "status": "ok"
  }
}
```

**What Happens on Backend**:
1. Validates coordinates
2. Stores location in Redis: `GEOADD user_locations <lon> <lat> <user_id>`
3. User becomes eligible for proximity notifications
4. Sends confirmation back to client

**Error Handling**:
- Invalid coordinates: No response (check backend logs)
- Redis error: No response (location not updated)

**Example Sequence**:
```
Mobile App                    Backend
    |                            |
    |---update_location--------->|
    |  {lat: -8.90, lon: 13.31}  |
    |                            |
    |<--location_updated---------|
    |  {status: "ok"}            |
    |                            |
```

---

### Server → Client Events

These are events that the backend sends to the mobile app.

#### 1. `location_updated`

**Purpose**: Confirms that the location update was successful.

**Payload**:
```json
{
  "event": "location_updated",
  "data": {
    "status": "ok"
  }
}
```

**What to Do**:
- Log success (optional)
- Update UI indicator if needed
- Continue periodic location updates

---

#### 2. `new_alert`

**Purpose**: Notifies the user of a new emergency alert in their proximity.

**When Received**:
- A government authority or verified user creates an alert
- The user's current location is within the alert's radius
- The user has an active WebSocket connection

**Payload**:
```json
{
  "event": "new_alert",
  "data": {
    "alert_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Tiroteio reportado na área do Morro Bento",
    "latitude": -8.842560,
    "longitude": 13.300120,
    "radius": 5000
  }
}
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| `alert_id` | UUID | Unique identifier for the alert |
| `message` | string | Human-readable alert description |
| `latitude` | float64 | Alert location latitude |
| `longitude` | float64 | Alert location longitude |
| `radius` | float64 | Alert broadcast radius in meters |

**What to Do**:
1. **Show notification**: Display high-priority alert to user
2. **Calculate distance**: Compute distance from user's location to alert
3. **Update map**: Show alert marker on map if applicable
4. **Log event**: Store alert for later retrieval
5. **Fetch details**: Call `GET /api/v1/alerts/{alert_id}` for full details

**Notification Behavior**:
- **Foreground**: Show in-app notification
- **Background**: System notification (via FCM)
- **Offline**: Push notification only (via FCM)

**Example UI Actions**:
```
┌─────────────────────────────┐
│  🚨 ALERTA DE RISCO         │
├─────────────────────────────┤
│  Tiroteio reportado na área │
│  do Morro Bento             │
│                             │
│  Distância: 1.2 km          │
│                             │
│  [Ver no Mapa] [Ignorar]    │
└─────────────────────────────┘
```

---

#### 3. `report_created`

**Purpose**: Notifies the user that someone reported a risk near their location.

**When Received**:
- Any user submits a report via the app
- The user's current location is within the report's notification radius
- The user has an active WebSocket connection

**Payload**:
```json
{
  "event": "report_created",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Buraco grande na estrada principal do Gamek",
    "latitude": -8.828765,
    "longitude": 13.247865
  }
}
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| `report_id` | UUID | Unique identifier for the report |
| `message` | string | Report description |
| `latitude` | float64 | Report location latitude |
| `longitude` | float64 | Report location longitude |

**What to Do**:
1. **Show notification**: Display report to user (lower priority than alerts)
2. **Update map**: Add report marker to map
3. **Store locally**: Cache for offline viewing
4. **Fetch details**: Call `GET /api/v1/reports/{report_id}` for full information

**Notification Priority**:
- **Lower** than `new_alert` (use different notification channel)
- Can be silent/non-intrusive if user preference is set

---

#### 4. `report_verified`

**Purpose**: Informs the user that their submitted report was verified by a moderator.

**When Received**:
- Only sent to the **original report creator**
- When a moderator verifies the report via `/api/v1/reports/{id}/verify`

**Payload**:
```json
{
  "event": "report_verified",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Seu relatório foi verificado."
  }
}
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| `report_id` | UUID | The verified report ID |
| `message` | string | Verification message (localized by backend) |

**What to Do**:
1. **Show notification**: "Your report was verified"
2. **Update report status**: Change status in local database
3. **Update UI**: Show verified badge on report
4. **Gamification**: Award points/badges if applicable

**Example UI**:
```
✅ Seu relatório foi verificado
   Obrigado por contribuir para a segurança da comunidade!
```

---

#### 5. `report_resolved`

**Purpose**: Notifies users that a previously reported issue has been resolved.

**When Received**:
- A moderator marks a report as resolved via `/api/v1/reports/{id}/resolve`
- The user is within the report's radius
- The user has an active WebSocket connection

**Payload**:
```json
{
  "event": "report_resolved",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Situação foi resolvida"
  }
}
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| `report_id` | UUID | The resolved report ID |
| `message` | string | Resolution message |

**What to Do**:
1. **Show notification**: "Issue resolved in your area"
2. **Update map**: Remove or grey out report marker
3. **Update report list**: Change status to resolved
4. **Refresh data**: Optionally fetch updated report details

---

## Integration Workflow

### Complete Integration Steps

#### Phase 1: Initial Setup

1. **Implement Authentication**
   - `POST /api/v1/auth/login` endpoint integration
   - JWT token storage (secure storage recommended)
   - Token refresh logic

2. **Implement WebSocket Connection**
   - WebSocket client library integration
   - Connection establishment with JWT header
   - Connection state management (connected, disconnected, reconnecting)

3. **Implement Location Services**
   - GPS permission handling
   - Location tracking (foreground/background)
   - Location update throttling (avoid excessive updates)

#### Phase 2: Message Handling

4. **Implement Message Parser**
   - JSON deserialization
   - Event type routing
   - Data validation

5. **Implement Event Handlers**
   - `update_location` sender
   - `location_updated` handler
   - `new_alert` handler
   - `report_created` handler
   - `report_verified` handler
   - `report_resolved` handler

#### Phase 3: User Experience

6. **Implement Notifications**
   - In-app notifications
   - System notifications (via FCM)
   - Notification channels/categories
   - Sound/vibration preferences

7. **Implement Map Integration**
   - Show alerts on map
   - Show reports on map
   - Distance calculation
   - Route avoidance (optional)

8. **Implement Local Storage**
   - Cache received alerts
   - Cache received reports
   - Offline viewing support

#### Phase 4: Edge Cases

9. **Implement Error Handling**
   - Connection failures
   - Authentication errors
   - Message parsing errors
   - Network changes

10. **Implement Reconnection Logic**
    - Exponential backoff
    - Connection state persistence
    - Automatic retry on network recovery

---

## Location Management

### Location Update Strategy

#### When to Send Location Updates

| Scenario | Frequency | Reason |
|----------|-----------|--------|
| **App Launch** | Once immediately | Register user location |
| **App Active** | Every 30-60s | Keep location fresh |
| **Significant Change** | On change >100m | Proximity accuracy |
| **Before Background** | Once | Last known position |
| **After Resume** | Once immediately | Update after background |

#### Location Permissions

**Required Permissions:**
- **iOS**: `NSLocationWhenInUseUsageDescription` or `NSLocationAlwaysUsageDescription`
- **Android**: `ACCESS_FINE_LOCATION` or `ACCESS_COARSE_LOCATION`

**Best Practices:**
- Request permissions contextually
- Explain why location is needed
- Provide graceful degradation if denied

#### Accuracy Requirements

| Accuracy Level | Use Case | Battery Impact |
|----------------|----------|----------------|
| **High (GPS)** | Precise notifications | High |
| **Balanced** | General proximity | Medium |
| **Low (Network)** | Basic vicinity | Low |

**Recommendation**: Use **Balanced** accuracy for periodic updates.

### Location Privacy

- **Storage**: User locations are stored in Redis temporarily
- **Retention**: No permanent location history (unless feature requires)
- **TTL**: Location data expires after 24 hours of inactivity
- **Encryption**: All data in transit via TLS (production)

---

## Notification Handling

### Notification Priority Levels

| Event | Priority | Sound | Vibration | Pop-up |
|-------|----------|-------|-----------|--------|
| `new_alert` | **Critical** | Yes (loud) | Yes (strong) | Yes |
| `report_created` | **Normal** | Yes (soft) | Yes (weak) | Optional |
| `report_verified` | **Low** | No | No | Optional |
| `report_resolved` | **Low** | No | No | Optional |

### Notification Channels (Android)

```
Channel: "alerts"
- Name: "Emergency Alerts"
- Importance: HIGH
- Sound: custom_alert_sound.mp3
- Vibration: [0, 500, 200, 500]

Channel: "reports"
- Name: "Risk Reports"
- Importance: DEFAULT
- Sound: default
- Vibration: [0, 200]

Channel: "updates"
- Name: "Report Updates"
- Importance: LOW
- Sound: none
- Vibration: none
```

### Deep Linking

Notifications should deep link to relevant screens:

| Event | Deep Link | Screen |
|-------|-----------|--------|
| `new_alert` | `riskplace://alert/{alert_id}` | Alert details |
| `report_created` | `riskplace://report/{report_id}` | Report details |
| `report_verified` | `riskplace://report/{report_id}` | User's report |
| `report_resolved` | `riskplace://report/{report_id}` | Report details |

---

## Error Scenarios

### Connection Errors

#### 1. Authentication Failed (401)

**Cause**: Invalid or expired JWT token

**Backend Response**:
```json
{
  "error": "unauthorized"
}
```

**What to Do**:
1. Clear stored JWT token
2. Prompt user to log in again
3. Retry connection with new token

#### 2. Connection Refused

**Cause**: Network issue, server down, or firewall

**What to Do**:
1. Show "Connection failed" message
2. Wait 5 seconds
3. Retry with exponential backoff
4. Max retries: 5

#### 3. Connection Timeout

**Cause**: Slow network or server overload

**What to Do**:
1. Cancel connection attempt after 10 seconds
2. Retry immediately once
3. Then follow exponential backoff

#### 4. Connection Dropped

**Cause**: Network change (WiFi → Cellular), server restart

**What to Do**:
1. Detect disconnect event
2. Wait 2 seconds
3. Attempt reconnection
4. Re-send last `update_location` after reconnect

### Message Errors

#### 1. Invalid JSON

**Cause**: Malformed message from server (rare)

**What to Do**:
1. Log error locally
2. Ignore message
3. Report to error tracking service (e.g., Sentry)

#### 2. Unknown Event Type

**Cause**: New event type not yet supported by mobile app

**What to Do**:
1. Log event type
2. Ignore gracefully
3. Update app when new version is available

#### 3. Missing Required Fields

**Cause**: Backend API change or bug

**What to Do**:
1. Log error with full payload
2. Attempt to display partial data
3. Show fallback UI

### Network Changes

#### WiFi → Cellular

**What to Do**:
1. WebSocket disconnects automatically
2. Detect network change
3. Wait 2 seconds for network stabilization
4. Reconnect WebSocket
5. Re-send location

#### Airplane Mode

**What to Do**:
1. Detect airplane mode
2. Pause reconnection attempts
3. Show "Offline" indicator
4. Resume on network recovery

---

## Best Practices

### Connection Management

1. **Single Connection**: Maintain only one WebSocket connection per user
2. **Reconnection**: Implement exponential backoff (2s, 4s, 8s, 16s, 32s)
3. **Heartbeat**: Send location update every 30-60s as heartbeat
4. **Close Gracefully**: Close WebSocket when user logs out

### Location Updates

1. **Throttle Updates**: Avoid sending updates more than once per 10 seconds
2. **Batch Coordinates**: No need to send every GPS reading
3. **Significant Changes**: Only update on meaningful movement (>50-100m)
4. **Battery Optimization**: Use appropriate accuracy level

### Message Processing

1. **Non-Blocking**: Process messages asynchronously
2. **Queue Messages**: Handle messages in order received
3. **Duplicate Detection**: Check if alert/report already received
4. **Debounce Notifications**: Avoid spamming user with notifications

### Error Handling

1. **Graceful Degradation**: App should work without WebSocket (use polling)
2. **User Feedback**: Show clear connection status
3. **Retry Logic**: Implement smart retry with backoff
4. **Logging**: Log all errors for debugging

### Security

1. **Secure Storage**: Store JWT in secure storage (Keychain, Keystore)
2. **Token Refresh**: Refresh JWT before expiration
3. **TLS Only**: Always use `wss://` in production
4. **Validate Messages**: Validate all incoming messages

### Performance

1. **Background Mode**: Close WebSocket when app is backgrounded (optional)
2. **Lazy Reconnect**: Reconnect when app returns to foreground
3. **Cache Data**: Cache alerts/reports locally
4. **Pagination**: Load old data via HTTP API, not WebSocket

---

## Testing Environment

### Development Server

**Base URL**: `ws://localhost:8000`  
**WebSocket**: `ws://localhost:8000/ws/alerts`

### Test Users

| Name | Email | Password | User ID | Location |
|------|-------|----------|---------|----------|
| Lopes Estevão | lopes@example.com | #Pwd1234 | 7aa9c0c0-14d3-4af9-a631-956c12a1f100 | -8.839987, 13.289437 |
| João Silva | joao@example.com | #Pwd1234 | f55c21ea-18e3-4fc9-99c3-d03b234bc110 | -8.915120, 13.242380 |
| Maria João | maria@example.com | #Pwd1234 | a11fd8dc-55d0-4e07-b0db-23f659ed3201 | -8.828765, 13.247865 |
| Carlos Domingos | carlos@example.com | #Pwd1234 | cc772485-bb16-4584-a8a4-3fd366478931 | -8.842560, 13.300120 |
| Ana Ferreira | ana@example.com | #Pwd1234 | 51bfcbfd-c896-4a7a-ae6a-79a0df2aab30 | -8.903290, 13.312540 |

### Testing Scenarios

#### Scenario 1: Basic Connection

1. Login as Ana (`ana@example.com`)
2. Connect to WebSocket with JWT
3. Send `update_location` with Ana's coordinates
4. Verify `location_updated` response

#### Scenario 2: Receive Alert

1. Login as Ana and Lopes (two devices/users)
2. Both connect and update location
3. Carlos creates alert at Morro Bento (via HTTP)
4. Ana receives `new_alert` (if within radius)

#### Scenario 3: Receive Report

1. Login as Maria
2. Connect and update location
3. João creates report at Zango 2 (via HTTP)
4. Maria receives `report_created` (if within radius)

#### Scenario 4: Reconnection

1. Connect as any user
2. Simulate network loss (airplane mode)
3. Restore network
4. Verify automatic reconnection
5. Verify location update is re-sent

### Testing Tools

**Recommended Tools:**
- **Postman**: WebSocket testing
- **websocat**: Command-line WebSocket client
- **Browser DevTools**: WebSocket debugging
- **Charles Proxy**: Network inspection

**Example with websocat:**
```bash
# Login to get token
TOKEN=$(curl -s -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ana@example.com","password":"#Pwd1234"}' \
  | jq -r '.token')

# Connect to WebSocket
websocat ws://localhost:8000/ws/alerts \
  -H="Authorization: Bearer $TOKEN"

# Send location update (type this after connection)
{"event":"update_location","data":{"latitude":-8.903290,"longitude":13.312540}}
```

---

## Appendix

### A. Complete Message Examples

#### Client Messages

**Update Location:**
```json
{
  "event": "update_location",
  "data": {
    "latitude": -8.903290,
    "longitude": 13.312540
  }
}
```

#### Server Messages

**Location Updated:**
```json
{
  "event": "location_updated",
  "data": {
    "status": "ok"
  }
}
```

**New Alert:**
```json
{
  "event": "new_alert",
  "data": {
    "alert_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Tiroteio reportado na área",
    "latitude": -8.839987,
    "longitude": 13.289437,
    "radius": 5000
  }
}
```

**Report Created:**
```json
{
  "event": "report_created",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Buraco grande na estrada principal",
    "latitude": -8.915120,
    "longitude": 13.242380
  }
}
```

**Report Verified:**
```json
{
  "event": "report_verified",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Seu relatório foi verificado."
  }
}
```

**Report Resolved:**
```json
{
  "event": "report_resolved",
  "data": {
    "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
    "message": "Situação foi resolvida"
  }
}
```

### B. HTTP API Endpoints

For additional data not provided via WebSocket:

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/v1/alerts/{id}` | Get full alert details |
| GET | `/api/v1/alerts` | List all alerts |
| POST | `/api/v1/alerts` | Create new alert (authorities only) |
| GET | `/api/v1/reports/{id}` | Get full report details |
| GET | `/api/v1/reports/nearby` | List reports near location |
| POST | `/api/v1/reports` | Submit new report |
| GET | `/api/v1/risks/types` | List risk types |

### C. Distance Calculation

To calculate distance from user to alert/report:

**Haversine Formula:**
```
a = sin²(Δlat/2) + cos(lat1) × cos(lat2) × sin²(Δlon/2)
c = 2 × atan2(√a, √(1−a))
d = R × c
```

Where:
- `R` = Earth's radius (6371 km)
- `lat1`, `lon1` = User coordinates
- `lat2`, `lon2` = Alert/report coordinates
- `d` = Distance in kilometers

**Libraries:**
- Flutter: `geolocator` package
- React Native: `geolib` library
- iOS: `CLLocation.distance(from:)`
- Android: `Location.distanceBetween()`

### D. Notification Radius

Different risk types have different default radii:

| Risk Type | Default Radius | Example |
|-----------|----------------|---------|
| Violence | 5000m (5km) | Shooting, assault |
| Fire | 3000m (3km) | Building fire, forest fire |
| Traffic | 2000m (2km) | Accident, road closure |
| Infrastructure | 1000m (1km) | Pothole, broken streetlight |
| Flood | 10000m (10km) | Flooding, heavy rain |

Use the HTTP API to fetch current radius values:
```bash
GET /api/v1/risks/types
```

### E. Error Codes Reference

| Code | Event | Meaning | Action |
|------|-------|---------|--------|
| 401 | Connection | Unauthorized | Re-authenticate |
| 500 | Connection | Server error | Retry later |
| N/A | Invalid message | JSON parse error | Log and ignore |
| N/A | Unknown event | Unsupported event | Log and ignore |

### F. Performance Benchmarks

Expected performance metrics:

| Metric | Target | Acceptable |
|--------|--------|------------|
| Connection time | <1s | <3s |
| Message latency | <500ms | <2s |
| Location update | <200ms | <1s |
| Notification delivery | <1s | <3s |
| Reconnection time | <2s | <5s |

### G. Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-11-12 | Initial release |

### H. Support & Contact

**Backend Team:**
- GitHub: https://github.com/risk-place-angola/backend-risk-place
- Discord: https://discord.gg/s2Nk4xYV
- Email: dev@riskplace.com

**Documentation:**
- API Docs: https://example.riskplace.com/docs/
- WebSocket Guide: [WEBSOCKET_NOTIFICATION_GUIDE.md](./WEBSOCKET_NOTIFICATION_GUIDE.md)

---

**End of Document**

*This documentation is subject to updates. Always check the latest version in the repository.*
