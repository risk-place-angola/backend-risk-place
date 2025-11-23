# Verification Code Flow

## Overview

After signup, users receive a 6-digit verification code via SMS (primary) or email (fallback). Code expires in 10 minutes with rate limiting to prevent brute force attacks.

## Security Measures

### Rate Limiting
- **Max Attempts**: 5 incorrect code attempts
- **Lockout Duration**: 15 minutes after max attempts
- **Resend Cooldown**: 60 seconds between resend requests
- **Code Expiration**: 10 minutes

### Attack Prevention
- Code space: 1,000,000 combinations (000000-999999)
- Auto-deleted after successful verification
- Attempt counter tracks failed verifications
- Lockout prevents brute force attacks

## Mobile Integration

### 1. Signup Flow

```
POST /api/v1/auth/signup
{
  "name": "João Silva",
  "email": "joao@example.com",
  "phone": "+244923456789",
  "password": "SecurePass123",
  "device_fcm_token": "fcm_token_here",
  "device_language": "pt"
}

Response 201:
{
  "id": "uuid-here"
}
```

**Action**: Redirect to verification screen immediately.

### 2. Verify Code

```
POST /api/v1/auth/confirm
{
  "email": "joao@example.com",
  "code": "123456"
}

Response 204: Success (no content)

Response 400:
{
  "error": "invalid verification code"
}
{
  "error": "verification code expired"
}
{
  "error": "Too many incorrect attempts. Wait 15 minutes"
}
```

**Mobile Handling**:
- **204**: Code verified → redirect to login or auto-login
- **400 (invalid code)**: Show error, allow retry (track attempts client-side)
- **400 (expired code)**: Show "Code expired, request new one" → enable resend button
- **400 (locked)**: Show "Too many attempts, wait 15 minutes" → disable input for 15min

### 3. Resend Code

```
POST /api/v1/auth/resend-code
{
  "email": "joao@example.com"
}

Response 204: Code resent successfully

Response 400:
{
  "error": "Wait 60 seconds before resending"
}
{
  "error": "Code already sent, please wait"
}
{
  "error": "Too many incorrect attempts. Wait 15 minutes"
}
```

**Mobile Handling**:
- **204**: Show "Code sent" toast, start 60s countdown on resend button
- **400 (cooldown)**: Disable resend button, show countdown timer
- **400 (locked)**: Disable all verification actions for 15 minutes

### 4. Login with Unverified Account

```
POST /api/v1/auth/login
{
  "identifier": "joao@example.com",
  "password": "SecurePass123"
}

Response 403:
{
  "error": "account not verified"
}
```

**Action**: Redirect to verification screen. New code sent automatically.

## UI/UX Recommendations

### Verification Screen

```
┌─────────────────────────────────┐
│  Verify Your Account            │
│                                  │
│  Code sent to +244 923 ***789   │
│                                  │
│  ┌─┬─┬─┬─┬─┬─┐                 │
│  │1│2│3│4│5│6│  [Verify]       │
│  └─┴─┴─┴─┴─┴─┘                 │
│                                  │
│  Didn't receive code?           │
│  [Resend in 48s]                │
│                                  │
│  Code expires in 9:32           │
└─────────────────────────────────┘
```

### State Management

**Local State**:
- `attemptCount`: Track verification attempts (0-5)
- `cooldownTimer`: Countdown for resend button (60s)
- `expirationTimer`: Show code expiration countdown (10min)
- `lockoutTimer`: If locked, show time remaining (15min)

**Error Display**:
```typescript
switch (response.error) {
  case "invalid verification code":
    incrementAttemptCount()
    showError(`Incorrect code. ${5 - attemptCount} attempts left`)
    break
    
  case "verification code expired":
    showError("Code expired")
    enableResendButton()
    break
    
  case "Too many incorrect attempts. Wait 15 minutes":
    disableAllInputs()
    startLockoutTimer(15 * 60 * 1000)
    showError("Account locked for 15 minutes")
    break
    
  case "Wait 60 seconds before resending":
    startCooldownTimer(60 * 1000)
    break
}
```

### Best Practices

1. **Auto-focus**: Focus on first code input on screen load
2. **Auto-advance**: Move to next input after digit entry
3. **Paste Support**: Allow pasting full 6-digit code
4. **Clear Button**: Allow clearing all inputs to start over
5. **Attempt Counter**: Show remaining attempts (visual feedback)
6. **Timers**: Display cooldown/expiration timers clearly
7. **Auto-resend**: Offer to resend when code expires
8. **SMS Detection**: Auto-fill code from SMS (Android/iOS)

### Error Messages by Language

**Portuguese**:
- "Código inválido. X tentativas restantes"
- "Código expirado. Solicite um novo código"
- "Muitas tentativas. Aguarde 15 minutos"
- "Aguarde 60 segundos antes de reenviar"

**English**:
- "Invalid code. X attempts left"
- "Code expired. Request a new code"
- "Too many attempts. Wait 15 minutes"
- "Wait 60 seconds before resending"

## Testing Scenarios

### Happy Path
1. User signs up → receives code within 30s
2. User enters correct code → verified immediately
3. User redirected to app home screen

### Error Cases
1. **Expired Code**: Wait 11 minutes → code invalid
2. **Wrong Code 5x**: Account locked for 15 minutes
3. **Resend Spam**: Try resend <60s → error shown
4. **SMS Delay**: Wait 90s → still valid for 8.5min

### Edge Cases
1. **Network Failure**: Retry with exponential backoff
2. **App Kill**: Persist state, resume verification
3. **Code Already Used**: Show "Code already verified"
4. **Lockout During Entry**: Disable inputs immediately

## API Response Codes Summary

| Code | Scenario | Mobile Action |
|------|----------|---------------|
| 204 | Success | Proceed to next screen |
| 400 | Invalid code | Show error, allow retry |
| 400 | Expired code | Enable resend button |
| 400 | Locked | Disable for 15min |
| 400 | Cooldown | Show countdown timer |
| 500 | Server error | Retry with backoff |
