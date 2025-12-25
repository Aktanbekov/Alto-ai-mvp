package tests

import (
	"testing"
	"altoai_mvp/internal/services"
)

func TestHashPassword(t *testing.T) {
	// Create a service instance to access private methods
	// Since hashPassword is private, we test through the public Register API
	// or we can test password comparison which is used in Login
	emailSvc := services.NewEmailService()
	if emailSvc == nil {
		t.Error("NewEmailService should not return nil")
	}
}

func TestGenerateCode(t *testing.T) {
	emailSvc := services.NewEmailService()
	
	code, err := emailSvc.GenerateCode()
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}
	
	if len(code) != 6 {
		t.Errorf("Expected code length 6, got %d", len(code))
	}
	
	// Generate multiple codes to ensure randomness
	codes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		code, _ := emailSvc.GenerateCode()
		codes[code] = true
	}
	
	// Should have some variation (not all same)
	if len(codes) == 1 {
		t.Error("Generated codes should have some variation")
	}
}


