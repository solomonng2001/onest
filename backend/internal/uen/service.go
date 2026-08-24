package uen

import (
	"regexp"
	"strings"
)

const (
	FormatBusiness     = "business"
	FormatLocalCompany = "local_company"
	FormatOtherEntity  = "other_entity"

	validFormatMessage = "UEN format is valid"
)

var (
	// ACRA-registered business: nnnnnnnnX
	businessPattern = regexp.MustCompile(`^[0-9]{8}[A-Z]$`)

	// ACRA-registered local company: yyyynnnnnX
	localCompanyPattern = regexp.MustCompile(`^[0-9]{9}[A-Z]$`)

	// Other entity: TyyPQnnnnX. Historical values can begin with S or R
	otherEntityPattern = regexp.MustCompile(`^[TSR][0-9]{2}[A-Z][A-Z0-9][0-9]{4}[A-Z]$`)
)

var validEntityTypeIndicators = map[string]struct{}{
	// Accounting and Corporate Regulatory Authority
	"LP": {}, // Limited Partnership
	"LL": {}, // Limited Liability Partnership
	"FC": {}, // Foreign Company
	"PF": {}, // Public Accounting Firm

	// Enterprise Singapore
	"RF": {}, // Representative Office

	// Islamic Religious Council of Singapore
	"MQ": {}, // Mosque
	"MM": {}, // Madrasah

	// Ministry of Communications and Information
	"NB": {}, // News Bureau

	// Ministry of Culture, Community and Youth
	"CC": {}, // Charity or Institution of a Public Character
	"CS": {}, // Co-operative Society
	"MB": {}, // Mutual Benefit Organisation

	// Ministry of Defence
	"FM": {}, // Foreign Military Unit

	// Ministry of Education
	"GS": {}, // Government or Government-aided School

	// Ministry of Foreign Affairs
	"DP": {}, // High Commission or Embassy
	"CP": {}, // Consulate
	"NR": {}, // International Organisation

	// Ministry of Health
	"CM": {}, // Medical Clinic
	"CD": {}, // Dental Clinic
	"MD": {}, // Medical and Dental Clinic
	"HS": {}, // Hospital
	"VH": {}, // Voluntary Welfare Home
	"CH": {}, // Commercial Home
	"MH": {}, // Maternity Home
	"CL": {}, // Clinical Laboratory
	"XL": {}, // X-ray Laboratory
	"CX": {}, // Clinical and X-ray Laboratory
	"HC": {}, // Healthcare Service Providers

	// Ministry of Law
	"RP": {}, // Foreign Law Practice Representative Office

	// Ministry of Manpower
	"TU": {}, // Trade Union

	// Ministry of National Development
	"TC": {}, // Town Council

	// Monetary Authority of Singapore
	"FB": {}, // Bank Representative Office
	"FN": {}, // Insurance Representative Office

	// People's Association
	"PA": {}, // PA Service
	"PB": {}, // Grassroots Unit

	// Registry of Societies
	"SS": {}, // Society

	// Singapore Land Authority
	"MC": {}, // Management Corporation
	"SM": {}, // Subsidiary Management Corporation

	// Government agencies and public bodies
	"GA": {}, // Organ of State, Ministry or Department
	"GB": {}, // Statutory Board or public body
}

// Service contains the UEN validation rules.
type Service struct{}

func NewService() *Service {
	return &Service{}
}

func isValidOtherEntity(uen string) bool {
	if !otherEntityPattern.MatchString(uen) {
		return false
	}

	pq := uen[3:5]

	_, exists := validEntityTypeIndicators[pq]
	return exists
}

func (s *Service) Validate(value string) ValidateResponse {
	normalized := strings.ToUpper(strings.TrimSpace(value))

	if normalized == "" {
		return ValidateResponse{
			UEN:     normalized,
			Valid:   false,
			Message: "UEN is required",
		}
	}

	if businessPattern.MatchString(normalized) {
		return ValidateResponse{
			UEN:     normalized,
			Valid:   true,
			Format:  FormatBusiness,
			Message: validFormatMessage,
		}
	}

	if localCompanyPattern.MatchString(normalized) {
		return ValidateResponse{
			UEN:     normalized,
			Valid:   true,
			Format:  FormatLocalCompany,
			Message: validFormatMessage,
		}
	}

	if isValidOtherEntity(normalized) {
		return ValidateResponse{
			UEN:     normalized,
			Valid:   true,
			Format:  FormatOtherEntity,
			Message: validFormatMessage,
		}
	}

	return ValidateResponse{
		UEN:     normalized,
		Valid:   false,
		Message: "UEN does not match a recognised format",
	}
}
