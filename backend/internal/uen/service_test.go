package uen

import "testing"

func TestServiceValidate(t *testing.T) {
	service := NewService()

	tests := []struct {
		name  string
		input string
		want  ValidateResponse
	}{
		{
			name:  "empty UEN",
			input: "",
			want: ValidateResponse{
				UEN:     "",
				Valid:   false,
				Message: "UEN is required",
			},
		},
		{
			name:  "whitespace-only UEN",
			input: "   \t\n",
			want: ValidateResponse{
				UEN:     "",
				Valid:   false,
				Message: "UEN is required",
			},
		},
		{
			name:  "valid business UEN",
			input: "12345678A",
			want: ValidateResponse{
				UEN:     "12345678A",
				Valid:   true,
				Format:  FormatBusiness,
				Message: validFormatMessage,
			},
		},
		{
			name:  "business UEN is trimmed and uppercased",
			input: "  12345678a  ",
			want: ValidateResponse{
				UEN:     "12345678A",
				Valid:   true,
				Format:  FormatBusiness,
				Message: validFormatMessage,
			},
		},
		{
			name:  "valid local company UEN",
			input: "202612345B",
			want: ValidateResponse{
				UEN:     "202612345B",
				Valid:   true,
				Format:  FormatLocalCompany,
				Message: validFormatMessage,
			},
		},
		{
			name:  "valid current other-entity UEN",
			input: "T26LL1234C",
			want: ValidateResponse{
				UEN:     "T26LL1234C",
				Valid:   true,
				Format:  FormatOtherEntity,
				Message: validFormatMessage,
			},
		},
		{
			name:  "valid historical S-prefix other-entity UEN",
			input: "S99SS1234D",
			want: ValidateResponse{
				UEN:     "S99SS1234D",
				Valid:   true,
				Format:  FormatOtherEntity,
				Message: validFormatMessage,
			},
		},
		{
			name:  "valid historical R-prefix other-entity UEN",
			input: "R00GA0001E",
			want: ValidateResponse{
				UEN:     "R00GA0001E",
				Valid:   true,
				Format:  FormatOtherEntity,
				Message: validFormatMessage,
			},
		},
		{
			name:  "business UEN has too few digits",
			input: "1234567A",
			want: ValidateResponse{
				UEN:     "1234567A",
				Valid:   false,
				Message: "UEN does not match a recognised format",
			},
		},
		{
			name:  "local company UEN ends in a digit",
			input: "2026123456",
			want: ValidateResponse{
				UEN:     "2026123456",
				Valid:   false,
				Message: "UEN does not match a recognised format",
			},
		},
		{
			name:  "other-entity UEN has unsupported prefix",
			input: "X26LL1234A",
			want: ValidateResponse{
				UEN:     "X26LL1234A",
				Valid:   false,
				Message: "UEN does not match a recognised format",
			},
		},
		{
			name:  "other-entity UEN has unknown entity indicator",
			input: "T26ZZ1234A",
			want: ValidateResponse{
				UEN:     "T26ZZ1234A",
				Valid:   false,
				Message: "UEN does not match a recognised format",
			},
		},
		{
			name:  "UEN contains internal whitespace",
			input: "1234 5678A",
			want: ValidateResponse{
				UEN:     "1234 5678A",
				Valid:   false,
				Message: "UEN does not match a recognised format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.Validate(tt.input)

			if got != tt.want {
				t.Fatalf("Validate(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidOtherEntityRecognisesEveryConfiguredIndicator(t *testing.T) {
	for indicator := range validEntityTypeIndicators {
		uen := "T26" + indicator + "1234A"

		if !isValidOtherEntity(uen) {
			t.Errorf("isValidOtherEntity(%q) = false, want true", uen)
		}
	}
}
