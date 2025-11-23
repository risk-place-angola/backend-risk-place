# HTTP Status Codes Reference

## Authentication Endpoints

### POST /auth/signup
- **201 Created**: User successfully registered
- **400 Bad Request**: Invalid request body or failed to create account
- **500 Internal Server Error**: Unexpected server error

### POST /auth/login
- **200 OK**: Login successful, returns access and refresh tokens
- **400 Bad Request**: Invalid credentials (wrong email/phone or password)
- **403 Forbidden**: Account not verified (verification code sent)
- **500 Internal Server Error**: Unexpected server error

### POST /auth/refresh
- **200 OK**: Tokens refreshed successfully
- **400 Bad Request**: Invalid request body or missing refresh token
- **401 Unauthorized**: Expired or invalid refresh token
- **403 Forbidden**: Account not verified
- **500 Internal Server Error**: Unexpected server error

### POST /auth/logout
- **200 OK**: Logout successful
- **401 Unauthorized**: Missing or invalid JWT token
- **500 Internal Server Error**: Unexpected server error

### POST /auth/password/forgot
- **200 OK**: Password reset code sent
- **400 Bad Request**: Invalid request body or user not found
- **500 Internal Server Error**: Unexpected server error

### POST /auth/password/reset
- **200 OK**: Password reset successfully
- **400 Bad Request**: Invalid request, user not found, or code not verified
- **500 Internal Server Error**: Unexpected server error

### POST /auth/confirm
- **204 No Content**: Account verified successfully
- **400 Bad Request**: Invalid request, user not found, invalid or expired code
- **405 Method Not Allowed**: Invalid HTTP method
- **500 Internal Server Error**: Unexpected server error

### POST /auth/resend-code
- **204 No Content**: Verification code resent successfully
- **400 Bad Request**: Invalid request body
- **500 Internal Server Error**: Failed to resend verification code

## User Endpoints

### GET /users/me
- **200 OK**: Returns current user profile
- **401 Unauthorized**: Missing or invalid JWT token
- **404 Not Found**: User not found
- **500 Internal Server Error**: Unexpected server error

### PUT /users/profile
- **200 OK**: Profile updated successfully
- **400 Bad Request**: Invalid request body
- **401 Unauthorized**: Missing or invalid JWT token
- **404 Not Found**: User not found
- **500 Internal Server Error**: Unexpected server error

## Common Error Response Format

```json
{
  "error": "error message description"
}
```

## Important Notes

1. **400 vs 401 vs 403**:
   - **400 Bad Request**: Invalid data, wrong credentials, validation errors
   - **401 Unauthorized**: Missing or invalid authentication token
   - **403 Forbidden**: Valid authentication but action not allowed (e.g., unverified account)

2. **Login with Email or Phone**:
   - Use field `identifier` to send either email or phone number
   - Backend will search both fields automatically

3. **Account Verification**:
   - After signup, verification code is sent via SMS (primary) or email (fallback)
   - Login attempts with unverified accounts return 403 and automatically resend verification code
   - Use `/auth/confirm` to verify account with the code
   - Use `/auth/resend-code` to request a new verification code

4. **Password Reset Flow**:
   - Request reset code with `/auth/password/forgot`
   - Verify code with `/auth/confirm`
   - Reset password with `/auth/password/reset` (requires verified code)

5. **Token Refresh**:
   - Access tokens expire after 1 hour
   - Use refresh token to get new access and refresh tokens
   - Refresh tokens are rotated on each refresh request
