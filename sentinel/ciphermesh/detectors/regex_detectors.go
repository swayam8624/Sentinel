package detectors

import (
	"fmt"
)

// CommonRegexDetectors returns a set of common regex-based detectors
func CommonRegexDetectors() ([]Detector, error) {
	var detectors []Detector

	// Social Security Number detector
	ssnDetector, err := NewRegexDetector(
		"ssn_detector",
		"pii",
		"ssn",
		`(?<!\d)(?!000|666|9\d{2})\d{3}-(?!00)\d{2}-(?!0000)\d{4}(?!\d)`,
		0.95,
		50,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSN detector: %w", err)
	}
	detectors = append(detectors, ssnDetector)

	// Credit Card detector
	ccDetector, err := NewRegexDetector(
		"credit_card_detector",
		"pci",
		"credit_card",
		`(?<!\d)(?:\d[ -]*?){13,19}\d(?!\d)`,
		0.90,
		50,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create credit card detector: %w", err)
	}
	detectors = append(detectors, ccDetector)

	// Email detector
	emailDetector, err := NewRegexDetector(
		"email_detector",
		"pii",
		"email",
		`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
		0.85,
		30,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email detector: %w", err)
	}
	detectors = append(detectors, emailDetector)

	// Phone number detector
	phoneDetector, err := NewRegexDetector(
		"phone_detector",
		"pii",
		"phone",
		`(\+\d{1,3}\s?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}`,
		0.80,
		30,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create phone detector: %w", err)
	}
	detectors = append(detectors, phoneDetector)

	// Driver's License detector (simplified)
	dlDetector, err := NewRegexDetector(
		"drivers_license_detector",
		"pii",
		"drivers_license",
		`[A-Z]\d{3,12}`,
		0.75,
		40,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver's license detector: %w", err)
	}
	detectors = append(detectors, dlDetector)

	return detectors, nil
}

// USFinancialDetectors returns detectors for US financial data
func USFinancialDetectors() ([]Detector, error) {
	var detectors []Detector

	// Bank Account Number detector
	bankDetector, err := NewRegexDetector(
		"bank_account_detector",
		"pii",
		"bank_account",
		`(?<!\d)\d{10,17}(?!\d)`,
		0.85,
		40,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create bank account detector: %w", err)
	}
	detectors = append(detectors, bankDetector)

	// Routing Number detector
	routingDetector, err := NewRegexDetector(
		"routing_number_detector",
		"pii",
		"routing_number",
		`(?<!\d)\d{9}(?!\d)`,
		0.90,
		40,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create routing number detector: %w", err)
	}
	detectors = append(detectors, routingDetector)

	return detectors, nil
}

// MedicalDetectors returns detectors for medical/healthcare data
func MedicalDetectors() ([]Detector, error) {
	var detectors []Detector

	// Medical Record Number detector
	mrnDetector, err := NewRegexDetector(
		"mrn_detector",
		"phi",
		"medical_record_number",
		`[A-Z]{2,3}\d{6,10}`,
		0.85,
		40,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MRN detector: %w", err)
	}
	detectors = append(detectors, mrnDetector)

	// Health Plan Beneficiary Number detector
	hpbDetector, err := NewRegexDetector(
		"hpbn_detector",
		"phi",
		"health_plan_beneficiary",
		`\d{8,12}`,
		0.80,
		40,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HPBN detector: %w", err)
	}
	detectors = append(detectors, hpbDetector)

	return detectors, nil
}
