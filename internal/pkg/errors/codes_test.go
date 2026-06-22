package errors

import (
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := ErrInvalidPhone
	if err.Error() != "手机号格式不正确" {
		t.Errorf("unexpected message: %s", err.Error())
	}
}

func TestAppError_WithStatus(t *testing.T) {
	tests := []struct {
		err      *AppError
		expected int
	}{
		{ErrInvalidPhone, http.StatusBadRequest},
		{ErrUnauthorized, http.StatusUnauthorized},
		{ErrJobNotFound, http.StatusNotFound},
		{ErrFileTooLarge, http.StatusBadRequest},
		{ErrInternal, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		status := tt.err.WithStatus()
		if status != tt.expected {
			t.Errorf("%s: expected status %d, got %d", tt.err.Error(), tt.expected, status)
		}
	}
}

func TestErrorCodes_Unique(t *testing.T) {
	codes := []*AppError{
		ErrInvalidPhone, ErrInvalidCode, ErrCodeExpired, ErrCodeSendFreq, ErrCodeSendFailed, ErrInvalidRole,
		ErrUnauthorized, ErrTokenExpired, ErrForbidden, ErrUserNotFound, ErrNameRequired, ErrCertRequired,
		ErrResumeNotFound, ErrResumeTitleEmpty, ErrInvalidJSON,
		ErrJobNotFound, ErrJobTitleEmpty, ErrJobNotOwned, ErrInvalidStatus,
		ErrAppNotFound, ErrAlreadyApplied, ErrAppNotOwned,
		ErrInvNotFound, ErrInvNotAuthorized, ErrAddressInvalid, ErrAppNotFoundForInv,
		ErrFileTooLarge, ErrInvalidFileType, ErrFileRequired, ErrUploadFailed,
		ErrInternal, ErrInvalidParam, ErrNotFound,
	}
	seen := make(map[int]bool)
	for _, e := range codes {
		if seen[e.Code] {
			t.Errorf("duplicate error code: %d (%s)", e.Code, e.Message)
		}
		seen[e.Code] = true
	}
}

func TestErrorCodeRanges(t *testing.T) {
	tests := []struct {
		code     int
		minRange int
		maxRange int
	}{
		{10001, 10000, 19999},
		{20001, 20000, 29999},
		{30001, 30000, 39999},
		{40001, 40000, 49999},
		{50001, 50000, 59999},
		{60001, 60000, 69999},
		{70001, 70000, 79999},
		{90001, 90000, 99999},
	}
	for _, tt := range tests {
		if tt.code < tt.minRange || tt.code > tt.maxRange {
			t.Errorf("code %d outside range %d-%d", tt.code, tt.minRange, tt.maxRange)
		}
	}
}
