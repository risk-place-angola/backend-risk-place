# WebSocket Notification System Documentation

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [WebSocket Connection](#websocket-connection)
- [Event Types](#event-types)
- [Testing Guide](#testing-guide)
- [Complete Examples](#complete-examples)

---

## Overview

The Risk Place Angola backend implements a real-time notification system using WebSockets to deliver instant alerts and reports to users based on their geographic location. The system supports:

- **Real-time location updates** from connected clients
- **Proximity-based notifications** using Redis geospatial queries
- **Push notifications** via Firebase Cloud Messaging (FCM)
- **Event-driven architecture** for alerts and reports
- **Multi-user broadcasting** within defined radius zones

---

## Architecture

### Components

1. **WebSocket Hub** (`internal/adapter/websocket/websocket_hub.go`)
   - Manages all active client connections
   - Handles client registration/unregistration
   - Broadcasts messages to connected clients
   - Processes incoming location updates

2. **WebSocket Client** (`internal/adapter/websocket/websocket_client.go`)
   - Represents individual user connections
   - Handles read/write operations
   - Maintains user-specific channels

3. **Location Store** (`internal/infra/location/redis_location_store.go`)
   - Stores user locations in Redis using geospatial indexing
   - Performs radius-based queries to find nearby users
   - Key: `user_locations`

4. **Event Dispatcher** (`internal/domain/event/dispatcher.go`)
   - Coordinates event handling across the system
   - Triggers notifications for domain events

### Notification Flow

```
┌─────────────┐
│   Client    │
│  (Mobile)   │
└──────┬──────┘
       │ 1. WebSocket Connect
       │    + JWT Auth
       ▼
┌─────────────────┐
│  WSHandler      │
│  Authenticate   │
└────────┬────────┘
         │ 2. Register Client
         ▼
┌─────────────────┐         ┌──────────────────┐
│   WebSocket     │◄────────│  Event Listener  │
│      Hub        │         │  (Alert/Report)  │
└────────┬────────┘         └──────────────────┘
         │                           ▲
         │ 3. Store Location         │ 5. Dispatch Event
         ▼                           │
┌─────────────────┐         ┌──────────────────┐
│  Redis Geo      │         │   Use Case       │
│  Location Store │         │  (Alert/Report)  │
└─────────────────┘         └──────────────────┘
         │
         │ 4. Find Users in Radius
         ▼
┌─────────────────┐
│  Broadcast to   │
│  Nearby Users   │
└─────────────────┘
```

### Geolocation System

The system uses **Redis GEO commands** to efficiently store and query user locations:

- **Storage**: `GEOADD user_locations <longitude> <latitude> <user_id>`
- **Search**: `GEOSEARCH user_locations FROMLONLAT <lon> <lat> BYRADIUS <meters> M`

When an alert or report is created:
1. The system queries Redis for all users within the specified radius
2. Notifications are sent via WebSocket to connected users
3. Push notifications are sent to offline users via FCM

---

## WebSocket Connection

### Endpoint

```
ws://localhost:8000/ws/alerts
```

Or in production:
```
wss://api.riskplace.com/ws/alerts
```

### Authentication

The WebSocket connection requires JWT authentication. Include the JWT token in the request:

**Option 1: Query Parameter**
```
ws://localhost:8000/ws/alerts?token=<JWT_TOKEN>
```

**Option 2: Authorization Header** (Recommended)
```javascript
const ws = new WebSocket('ws://localhost:8000/ws/alerts');
// Send auth in first message or validate via middleware
```

The JWT token is validated by `AuthMiddleware.ValidateJWTFromRequest()` before upgrading the connection.

### Connection Lifecycle

1. **Connection Established**
   - Client connects with valid JWT
   - Server creates a `Client` instance
   - Client is registered in the Hub
   - Server logs: `websocket client connected user_id=<user_id>`

2. **Active Connection**
   - Client can send location updates
   - Client receives real-time notifications
   - Bidirectional communication maintained

3. **Disconnection**
   - Client unregisters from Hub
   - Client channel is closed
   - Connection cleanup performed

---

## Event Types

### 1. Client to Server Events

#### `update_location`

Updates the user's current location in the system.

**Request Payload:**
```json
{
    "event": "update_location",
    "data": {
        "latitude": -8.903290,
        "longitude": 13.312540
    }
}
```

**Fields:**
- `event` (string): Must be `"update_location"`
- `data.latitude` (float64): User's latitude coordinate
- `data.longitude` (float64): User's longitude coordinate

**Response:**
```json
{
    "event": "location_updated",
    "data": {
        "status": "ok"
    }
}
```

**What Happens:**
1. Location is stored in Redis: `GEOADD user_locations <longitude> <latitude> <user_id>`
2. User becomes visible for proximity-based notifications
3. Confirmation is sent back to the client

---

### 2. Server to Client Events

#### `new_alert`

Sent when a new alert is created within the user's proximity.

**Notification Payload:**
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

**Fields:**
- `alert_id` (UUID): Unique identifier for the alert
- `message` (string): Alert description
- `latitude` (float64): Alert location latitude
- `longitude` (float64): Alert location longitude  
- `radius` (float64): Alert broadcast radius in meters

**Trigger:**
- Alert created via `/api/v1/alerts` endpoint
- User is within the alert's radius
- Event `AlertCreatedEvent` is dispatched

**Push Notification:**
If the user is offline, they receive an FCM push notification:
- **Title**: 🚨 Alerta de Risco
- **Body**: `<message>`
- **Data**: `{ "alert_id": "<uuid>" }`

---

#### `report_created`

Sent when a new report is created within the user's proximity.

**Notification Payload:**
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

**Fields:**
- `report_id` (UUID): Unique identifier for the report
- `message` (string): Report description
- `latitude` (float64): Report location latitude
- `longitude` (float64): Report location longitude

**Trigger:**
- Report created via `/api/v1/reports` endpoint
- User is within the report's radius (based on risk type default radius)
- Event `ReportCreatedEvent` is dispatched

**Push Notification:**
If the user is offline:
- **Title**: 📍 Novo Relato de Risco
- **Body**: `<message>`
- **Data**: `{ "report_id": "<uuid>" }`

---

#### `report_verified`

Sent to the report creator when their report is verified by a moderator.

**Notification Payload:**
```json
{
    "event": "report_verified",
    "data": {
        "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
        "message": "Seu relatório foi verificado."
    }
}
```

**Fields:**
- `report_id` (UUID): The verified report ID
- `message` (string): Verification confirmation message

**Trigger:**
- Moderator verifies report via `/api/v1/reports/{id}/verify`
- Event `ReportVerifiedEvent` is dispatched
- Only sent to the original report creator

---

#### `report_resolved`

Sent when a report is marked as resolved, notifying nearby users.

**Notification Payload:**
```json
{
    "event": "report_resolved",
    "data": {
        "report_id": "7d3a4b10-f29c-41d4-a716-446655440001",
        "message": "Situação foi resolvida"
    }
}
```

**Fields:**
- `report_id` (UUID): The resolved report ID
- `message` (string): Resolution message

**Trigger:**
- Moderator resolves report via `/api/v1/reports/{id}/resolve`
- Event `ReportResolvedEvent` is dispatched
- Sent to all users within the report's radius

---

## Testing Guide

### Test Users (from seed data)

| Name            | Email                  | User ID                                | Location (Lat, Lon)       | Neighborhood  |
|-----------------|------------------------|----------------------------------------|---------------------------|---------------|
| Lopes Estevão   | lopes@example.com      | 7aa9c0c0-14d3-4af9-a631-956c12a1f100   | -8.839987, 13.289437      | Benfica       |
| João Silva      | joao@example.com       | f55c21ea-18e3-4fc9-99c3-d03b234bc110   | -8.915120, 13.242380      | Zango 2       |
| Maria João      | maria@example.com      | a11fd8dc-55d0-4e07-b0db-23f659ed3201   | -8.828765, 13.247865      | Gamek         |
| Carlos Domingos | carlos@example.com     | cc772485-bb16-4584-a8a4-3fd366478931   | -8.842560, 13.300120      | Morro Bento   |
| Ana Ferreira    | ana@example.com        | 51bfcbfd-c896-4a7a-ae6a-79a0df2aab30   | -8.903290, 13.312540      | Sequele       |

**Default Password for all test users:** `#Pwd1234`

### Prerequisites

1. **Start the backend server:**
   ```bash
   make run
   ```

2. **Ensure Redis is running:**
   ```bash
   docker-compose up -d redis
   ```

3. **Login and obtain JWT token:**
   ```bash
   curl -X POST http://localhost:8000/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "ana@example.com",
       "password": "#Pwd1234"
     }'
   ```

   **Response:**
   ```json
   {
     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
     "user": { ... }
   }
   ```

### Test Tools

#### Option 1: Using `websocat` (Recommended)

Install websocat:
```bash
brew install websocat
```

Connect to WebSocket:
```bash
websocat ws://localhost:8000/ws/alerts \
  -H="Authorization: Bearer <JWT_TOKEN>"
```

#### Option 2: Using JavaScript/Browser Console

```javascript
const token = "YOUR_JWT_TOKEN";
const ws = new WebSocket('ws://localhost:8000/ws/alerts');

ws.onopen = () => {
    console.log('Connected to WebSocket');
    
    // Update location
    ws.send(JSON.stringify({
        event: "update_location",
        data: {
            latitude: -8.903290,
            longitude: 13.312540
        }
    }));
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('WebSocket closed');
};
```

#### Option 3: Using Postman

1. Create a new WebSocket Request
2. URL: `ws://localhost:8000/ws/alerts`
3. Headers: Add `Authorization: Bearer <JWT_TOKEN>`
4. Connect and send messages

---

## Complete Examples

### Example 1: Ana Updates Her Location

**Scenario:** Ana (at Sequele) connects and updates her location.

**Step 1: Connect WebSocket**
```bash
websocat ws://localhost:8000/ws/alerts \
  -H="Authorization: Bearer eyJhbGci..."
```

**Step 2: Send Location Update**
```json
{
    "event": "update_location",
    "data": {
        "latitude": -8.903290,
        "longitude": 13.312540
    }
}
```

**Step 3: Receive Confirmation**
```json
{
    "event": "location_updated",
    "data": {
        "status": "ok"
    }
}
```

**Backend Log:**
```
INFO websocket client connected user_id=51bfcbfd-c896-4a7a-ae6a-79a0df2aab30
INFO location updated user_id=51bfcbfd-c896-4a7a-ae6a-79a0df2aab30
```

---

### Example 2: Carlos Creates an Alert Near Lopes

**Scenario:** Carlos creates an alert at Morro Bento. Lopes (at Benfica, ~1.5km away) should receive it.

**Step 1: Lopes Connects and Updates Location**

As Lopes (login with `lopes@example.com`):
```json
{
    "event": "update_location",
    "data": {
        "latitude": -8.839987,
        "longitude": 13.289437
    }
}
```

**Step 2: Carlos Creates Alert**

As Carlos (via HTTP):
```bash
curl -X POST http://localhost:8000/api/v1/alerts \
  -H "Authorization: Bearer <CARLOS_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "risk_type_id": "uuid-of-tiroteio",
    "risk_topic_id": "uuid-of-violencia",
    "message": "Tiroteio reportado na área do Morro Bento",
    "latitude": -8.842560,
    "longitude": 13.300120,
    "radius": 5000,
    "severity": "high"
  }'
```

**Step 3: Lopes Receives Alert (WebSocket)**

Lopes' WebSocket receives:
```json
{
    "event": "new_alert",
    "data": {
        "alert_id": "9c8e7f10-a29b-41d4-b716-446655440123",
        "message": "Tiroteio reportado na área do Morro Bento",
        "latitude": -8.842560,
        "longitude": 13.300120,
        "radius": 5000
    }
}
```

**Backend Logs:**
```
INFO broadcasting alert alert_id=9c8e7f10-... user_count=1
INFO created alert notification alert_id=9c8e7f10-... user_id=7aa9c0c0-...
```

---

### Example 3: Maria Creates Report, João Receives Notification

**Scenario:** Maria reports a pothole at Gamek. João (at Zango 2, ~15km away) is too far but Ana (at Sequele, ~9km away) might receive it depending on the risk type radius.

**Step 1: Multiple Users Connect**

Ana connects:
```json
{
    "event": "update_location",
    "data": {
        "latitude": -8.903290,
        "longitude": 13.312540
    }
}
```

João connects:
```json
{
    "event": "update_location",
    "data": {
        "latitude": -8.915120,
        "longitude": 13.242380
    }
}
```

**Step 2: Maria Creates Report**

As Maria:
```bash
curl -X POST http://localhost:8000/api/v1/reports \
  -H "Authorization: Bearer <MARIA_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "risk_type_id": "uuid-of-buraco",
    "risk_topic_id": "uuid-of-infraestrutura",
    "description": "Buraco grande na estrada principal do Gamek",
    "province": "Luanda",
    "municipality": "Talatona",
    "neighborhood": "Gamek",
    "latitude": -8.828765,
    "longitude": 13.247865
  }'
```

**Step 3: Nearby Users Receive Notification**

Users within the risk type's default radius receive:
```json
{
    "event": "report_created",
    "data": {
        "report_id": "b4f8c210-d19c-41d4-c716-446655440789",
        "message": "Buraco grande na estrada principal do Gamek",
        "latitude": -8.828765,
        "longitude": 13.247865
    }
}
```

**Backend Logs:**
```
INFO created report report_id=b4f8c210-...
INFO found users in radius count=2
INFO created report notification report_id=b4f8c210-... user_id=<user1>
INFO created report notification report_id=b4f8c210-... user_id=<user2>
```

---

### Example 4: Report Verification and Resolution Flow

**Scenario:** João creates a report, a moderator verifies it, then resolves it.

**Step 1: João Creates Report**

João creates a report and receives confirmation via HTTP response.

**Step 2: Moderator Verifies Report**

```bash
curl -X POST http://localhost:8000/api/v1/reports/{report_id}/verify \
  -H "Authorization: Bearer <MODERATOR_JWT>"
```

**Step 3: João Receives Verification (WebSocket)**

João's WebSocket receives:
```json
{
    "event": "report_verified",
    "data": {
        "report_id": "b4f8c210-d19c-41d4-c716-446655440789",
        "message": "Seu relatório foi verificado."
    }
}
```

**Step 4: Moderator Resolves Report**

```bash
curl -X POST http://localhost:8000/api/v1/reports/{report_id}/resolve \
  -H "Authorization: Bearer <MODERATOR_JWT>"
```

**Step 5: Nearby Users Receive Resolution**

All users within radius receive:
```json
{
    "event": "report_resolved",
    "data": {
        "report_id": "b4f8c210-d19c-41d4-c716-446655440789",
        "message": "Situação foi resolvida"
    }
}
```

---

### Example 5: Testing Multiple Concurrent Connections

**Scenario:** Test system with all 5 users connected simultaneously.

**Terminal 1 - Ana:**
```bash
websocat ws://localhost:8000/ws/alerts -H="Authorization: Bearer <ANA_TOKEN>"
# Send: {"event":"update_location","data":{"latitude":-8.903290,"longitude":13.312540}}
```

**Terminal 2 - João:**
```bash
websocat ws://localhost:8000/ws/alerts -H="Authorization: Bearer <JOAO_TOKEN>"
# Send: {"event":"update_location","data":{"latitude":-8.915120,"longitude":13.242380}}
```

**Terminal 3 - Maria:**
```bash
websocat ws://localhost:8000/ws/alerts -H="Authorization: Bearer <MARIA_TOKEN>"
# Send: {"event":"update_location","data":{"latitude":-8.828765,"longitude":13.247865}}
```

**Terminal 4 - Carlos:**
```bash
websocat ws://localhost:8000/ws/alerts -H="Authorization: Bearer <CARLOS_TOKEN>"
# Send: {"event":"update_location","data":{"latitude":-8.842560,"longitude":13.300120}}
```

**Terminal 5 - Lopes:**
```bash
websocat ws://localhost:8000/ws/alerts -H="Authorization: Bearer <LOPES_TOKEN>"
# Send: {"event":"update_location","data":{"latitude":-8.839987,"longitude":13.289437}}
```

**Create Alert with Large Radius:**
```bash
curl -X POST http://localhost:8000/api/v1/alerts \
  -H "Authorization: Bearer <ANY_USER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "risk_type_id": "uuid",
    "risk_topic_id": "uuid",
    "message": "Teste de broadcast para todos",
    "latitude": -8.870,
    "longitude": 13.270,
    "radius": 20000,
    "severity": "high"
  }'
```

**Expected:** All 5 connected users should receive the alert notification simultaneously.

---

## Error Handling

### Connection Errors

**Unauthorized (401):**
```json
{
    "error": "unauthorized"
}
```
**Cause:** Invalid or missing JWT token

**Upgrade Failed (500):**
```json
{
    "error": "Failed to upgrade to WebSocket"
}
```
**Cause:** WebSocket upgrade error

### Message Errors

**Invalid Message Format:**
- No response sent
- Backend logs: `invalid message: <error>`

**Unknown Event Type:**
- No response sent  
- Backend logs: `unknown event type: <event_name>`

**Location Update Failed:**
- No response sent
- Backend logs: `failed to update location: <error>`

---

## Performance Considerations

### Redis Geospatial Queries

- **Average query time:** < 10ms for 10,000 users
- **Index key:** `user_locations`
- **Coordinate precision:** 6 decimal places (~0.11m)

### WebSocket Scalability

- **Max connections per instance:** ~10,000 (configurable)
- **Message buffer size:** 256 bytes per client
- **Broadcasting:** O(n) where n = users in radius

### Optimization Tips

1. **Limit radius queries:** Use reasonable radius values (< 50km)
2. **Throttle location updates:** Update location every 30-60 seconds
3. **Connection pooling:** Reuse connections, avoid frequent reconnects
4. **Push notifications:** Rely on FCM for offline users

---

## Security

### Authentication
- JWT required for WebSocket connections
- Token validated before upgrade
- User ID extracted from JWT claims

### Authorization
- Users only receive notifications for their location
- No cross-user data exposure
- Moderator-only endpoints protected

### Data Privacy
- Locations stored temporarily in Redis
- No permanent location history (unless required)
- TTL can be set on location data

---

## Troubleshooting

### WebSocket Not Connecting

1. **Check JWT token validity:**
   ```bash
   curl http://localhost:8000/api/v1/users/me \
     -H "Authorization: Bearer <TOKEN>"
   ```

2. **Verify WebSocket upgrade:**
   - Check browser console for errors
   - Ensure correct protocol (ws:// or wss://)

3. **Check backend logs:**
   ```bash
   tail -f logs/app.log | grep websocket
   ```

### Not Receiving Notifications

1. **Verify location was updated:**
   ```bash
   redis-cli GEOPOS user_locations <user_id>
   ```

2. **Check radius calculation:**
   ```bash
   redis-cli GEOSEARCH user_locations FROMLONLAT <lon> <lat> BYRADIUS <meters> M
   ```

3. **Confirm event dispatch:**
   - Check backend logs for `broadcasting alert` or `found users in radius`

4. **Test connection:**
   - Send `update_location` and verify `location_updated` response

### High Latency

1. **Check Redis performance:**
   ```bash
   redis-cli --latency
   ```

2. **Monitor WebSocket connections:**
   ```bash
   netstat -an | grep :8080 | grep ESTABLISHED | wc -l
   ```

3. **Review backend metrics:**
   - Connection count
   - Message throughput
   - GEO query performance

---

## API Reference

### WebSocket Endpoint

**URL:** `/ws/alerts`  
**Protocol:** WebSocket  
**Authentication:** JWT (Bearer token)

### Client Events

| Event              | Description                    | Payload Fields                      |
|--------------------|--------------------------------|-------------------------------------|
| `update_location`  | Update user's current location | `latitude` (float), `longitude` (float) |

### Server Events

| Event              | Description                      | Payload Fields                                             |
|--------------------|----------------------------------|------------------------------------------------------------|
| `location_updated` | Location update confirmation     | `status` (string)                                          |
| `new_alert`        | New alert in user's proximity    | `alert_id`, `message`, `latitude`, `longitude`, `radius`   |
| `report_created`   | New report in user's proximity   | `report_id`, `message`, `latitude`, `longitude`            |
| `report_verified`  | User's report was verified       | `report_id`, `message`                                     |
| `report_resolved`  | Nearby report was resolved       | `report_id`, `message`                                     |

---

## Additional Resources

- **Swagger API Docs:** http://localhost:8000/docs/
- **WebSocket RFC:** https://tools.ietf.org/html/rfc6455
- **Redis GEO Commands:** https://redis.io/commands/geosearch
- **Gorilla WebSocket:** https://github.com/gorilla/websocket

---

## Support

For issues or questions:
- **GitHub Issues:** https://github.com/risk-place-angola/backend-risk-place/issues
- **Email:** support@riskplace.com

---

**Last Updated:** November 11, 2025  
**Version:** 1.0.0
