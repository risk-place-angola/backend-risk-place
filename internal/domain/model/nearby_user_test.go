package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAssignAvatar_NoIndexOutOfRange(t *testing.T) {
	// Test with multiple UUIDs including ones that might generate negative hashes
	testCases := []struct {
		name string
		uuid uuid.UUID
	}{
		{
			name: "UUID that caused panic",
			uuid: uuid.MustParse("8bc22b4a-ad8a-4365-950c-39cc66e06769"),
		},
		{
			name: "Random UUID 1",
			uuid: uuid.New(),
		},
		{
			name: "Random UUID 2",
			uuid: uuid.New(),
		},
		{
			name: "Random UUID 3",
			uuid: uuid.New(),
		},
		{
			name: "Nil UUID",
			uuid: uuid.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			avatarID, color := AssignAvatar(tc.uuid)

			// Validate avatar ID is within expected range
			assert.GreaterOrEqual(t, avatarID, 1, "Avatar ID should be at least 1")
			assert.LessOrEqual(t, avatarID, totalAvatars, "Avatar ID should not exceed total avatars")

			// Validate color is not empty
			assert.NotEmpty(t, color, "Color should not be empty")

			// Validate color starts with #
			assert.Equal(t, "#", color[0:1], "Color should start with #")
		})
	}
}

func TestGenerateAnonymousID_Consistency(t *testing.T) {
	userID := uuid.New()

	// Generate ID multiple times - should be consistent for same UUID
	id1 := GenerateAnonymousID(userID)
	id2 := GenerateAnonymousID(userID)

	assert.Equal(t, id1, id2, "Anonymous ID should be deterministic for same UUID")
	assert.Contains(t, id1, "neter_", "Anonymous ID should contain neter_ prefix")
	assert.Equal(t, len("neter_")+anonymousIDLength, len(id1), "Anonymous ID should have correct length")
}

func TestGenerateAnonymousID_NoIndexOutOfRange(t *testing.T) {
	// Test with multiple UUIDs including ones that might generate negative hashes
	testCases := []uuid.UUID{
		uuid.MustParse("8bc22b4a-ad8a-4365-950c-39cc66e06769"),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.Nil,
	}

	for _, userID := range testCases {
		t.Run(userID.String(), func(t *testing.T) {
			// Should not panic
			anonymousID := GenerateAnonymousID(userID)

			assert.NotEmpty(t, anonymousID, "Anonymous ID should not be empty")
			assert.Contains(t, anonymousID, "neter_", "Anonymous ID should contain neter_ prefix")
		})
	}
}

func TestNewUserLocation(t *testing.T) {
	userID := uuid.New()
	deviceID := "device-123"
	lat := 40.7128
	lon := -74.0060
	speed := 5.0
	heading := 90.0

	// Test non-anonymous user
	location := NewUserLocation(userID, deviceID, lat, lon, speed, heading, false)

	assert.NotNil(t, location)
	assert.Equal(t, userID, location.UserID)
	assert.Equal(t, deviceID, location.DeviceID)
	assert.Equal(t, lat, location.Latitude)
	assert.Equal(t, lon, location.Longitude)
	assert.Equal(t, speed, location.Speed)
	assert.Equal(t, heading, location.Heading)
	assert.False(t, location.IsAnonymous)
	assert.NotZero(t, location.AvatarID)
	assert.NotEmpty(t, location.Color)

	// Test anonymous user
	anonymousLocation := NewUserLocation(userID, deviceID, lat, lon, speed, heading, true)
	assert.True(t, anonymousLocation.IsAnonymous)
}
