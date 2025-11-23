# Email Fallback Response

## Overview

When SMS delivery fails (Twilio issues, invalid phone, etc), the system automatically sends verification codes via email and returns a specific response indicating the fallback method was used.

## Response Format

### Signup with Email Fallback

```
POST /api/v1/auth/signup

Response 201:
{
  "success": true,
  "message": "verification code sent via email",
  "data": {
    "id": "uuid-here",
    "email": "user@example.com"
  }
}
```

### Resend Code with Email Fallback

```
POST /api/v1/auth/resend-code

Response 200:
{
  "success": true,
  "message": "verification code sent via email",
  "data": {
    "identifier": "user@example.com"
  }
}
```

## Mobile Integration

### Detection

```typescript
if (response.success && response.message === "verification code sent via email") {
  showNotification("Code sent to your email", "info")
  updateUI({
    method: "email",
    destination: response.data.email || response.data.identifier
  })
}
```

### UI Recommendations

```
┌─────────────────────────────────┐
│  Check Your Email               │
│                                  │
│  ℹ️  Code sent to your email   │
│  user@example.com               │
│                                  │
│  ┌─┬─┬─┬─┬─┬─┐                 │
│  │ │ │ │ │ │ │  [Verify]       │
│  └─┴─┴─┴─┴─┴─┘                 │
│                                  │
│  Didn't receive code?           │
│  [Resend]                       │
└─────────────────────────────────┘
```

## Status Codes

| Code | Scenario | Response |
|------|----------|----------|
| 201 | Signup - SMS sent | `{"id": "..."}` |
| 201 | Signup - Email fallback | `{"success": true, "message": "verification code sent via email", "data": {...}}` |
| 204 | Resend - SMS sent | No content |
| 200 | Resend - Email fallback | `{"success": true, "message": "verification code sent via email", "data": {...}}` |
