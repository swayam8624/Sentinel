package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sentinel-platform/sentinel/sentinel/crypto/merkle"
	"go.opentelemetry.io/otel/attribute"
)

// ObservabilityInterface defines the interface for observability operations
type ObservabilityInterface interface {
	RecordMetric(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue)
	LogEvent(level string, message string, fields map[string]interface{})
}

// AuditFramework handles security auditing and compliance reporting
type AuditFramework struct {
	merkleTree    *merkle.MerkleTree
	observability ObservabilityInterface
	events        []AuditEvent
}

// AuditEvent represents a security audit event
type AuditEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	Source      string                 `json:"source"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Severity    string                 `json:"severity"` // low, medium, high, critical
}

// ComplianceReport represents a compliance report
type ComplianceReport struct {
	ReportID         string        `json:"report_id"`
	GeneratedAt      time.Time     `json:"generated_at"`
	PeriodStart      time.Time     `json:"period_start"`
	PeriodEnd        time.Time     `json:"period_end"`
	Findings         []Finding     `json:"findings"`
	ComplianceStatus string        `json:"compliance_status"` // compliant, non_compliant, pending
	Summary          ReportSummary `json:"summary"`
}

// Finding represents a security finding
type Finding struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"` // low, medium, high, critical
	Category    string                 `json:"category"` // policy, crypto, access, etc.
	Status      string                 `json:"status"`   // open, resolved, mitigated
	Details     map[string]interface{} `json:"details"`
	Remediation string                 `json:"remediation"`
}

// ReportSummary contains summary statistics for a compliance report
type ReportSummary struct {
	TotalEvents      int            `json:"total_events"`
	EventsByType     map[string]int `json:"events_by_type"`
	EventsBySeverity map[string]int `json:"events_by_severity"`
	CriticalFindings int            `json:"critical_findings"`
	HighFindings     int            `json:"high_findings"`
	MediumFindings   int            `json:"medium_findings"`
	LowFindings      int            `json:"low_findings"`
}

// NewAuditFramework creates a new audit framework
func NewAuditFramework(merkleTree *merkle.MerkleTree, obs ObservabilityInterface) *AuditFramework {
	return &AuditFramework{
		merkleTree:    merkleTree,
		observability: obs,
		events:        make([]AuditEvent, 0),
	}
}

// LogEvent logs a security audit event
func (af *AuditFramework) LogEvent(ctx context.Context, eventType, source, description string, details map[string]interface{}, severity string) error {
	event := AuditEvent{
		ID:          generateEventID(),
		Timestamp:   time.Now(),
		EventType:   eventType,
		Source:      source,
		Description: description,
		Details:     details,
		Severity:    severity,
	}

	// Add to events list
	af.events = append(af.events, event)

	// Add to Merkle tree for tamper evidence
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if af.merkleTree != nil {
		af.merkleTree.AddLeaf(eventData)
	}

	// Log the event
	if af.observability != nil {
		af.observability.LogEvent("AUDIT", description, details)
		af.observability.RecordMetric(ctx, "request.count", 1,
			attribute.String("event_type", eventType),
			attribute.String("severity", severity))
	}

	return nil
}

// GenerateComplianceReport generates a compliance report for a time period
func (af *AuditFramework) GenerateComplianceReport(ctx context.Context, start, end time.Time) (*ComplianceReport, error) {
	// Filter events for the time period
	var periodEvents []AuditEvent
	for _, event := range af.events {
		if event.Timestamp.After(start) && event.Timestamp.Before(end) {
			periodEvents = append(periodEvents, event)
		}
	}

	// Convert events to findings
	var findings []Finding
	for _, event := range periodEvents {
		finding := Finding{
			ID:          event.ID,
			Timestamp:   event.Timestamp,
			Title:       event.EventType,
			Description: event.Description,
			Severity:    event.Severity,
			Category:    categorizeEvent(event.EventType),
			Status:      "open",
			Details:     event.Details,
			Remediation: getRemediationForEvent(event.EventType),
		}
		findings = append(findings, finding)
	}

	// Generate summary statistics
	summary := af.generateSummary(periodEvents, findings)

	// Determine compliance status
	complianceStatus := "compliant"
	if summary.CriticalFindings > 0 || summary.HighFindings > 5 {
		complianceStatus = "non_compliant"
	} else if summary.MediumFindings > 10 {
		complianceStatus = "pending"
	}

	report := &ComplianceReport{
		ReportID:         generateReportID(),
		GeneratedAt:      time.Now(),
		PeriodStart:      start,
		PeriodEnd:        end,
		Findings:         findings,
		ComplianceStatus: complianceStatus,
		Summary:          summary,
	}

	return report, nil
}

// generateSummary generates summary statistics for a compliance report
func (af *AuditFramework) generateSummary(events []AuditEvent, findings []Finding) ReportSummary {
	summary := ReportSummary{
		TotalEvents:      len(events),
		EventsByType:     make(map[string]int),
		EventsBySeverity: make(map[string]int),
	}

	// Count events by type
	for _, event := range events {
		summary.EventsByType[event.EventType]++
		summary.EventsBySeverity[event.Severity]++
	}

	// Count findings by severity
	for _, finding := range findings {
		switch finding.Severity {
		case "critical":
			summary.CriticalFindings++
		case "high":
			summary.HighFindings++
		case "medium":
			summary.MediumFindings++
		case "low":
			summary.LowFindings++
		}
	}

	return summary
}

// GetRecentEvents returns recent audit events
func (af *AuditFramework) GetRecentEvents(count int) []AuditEvent {
	if count > len(af.events) {
		count = len(af.events)
	}

	start := len(af.events) - count
	return af.events[start:]
}

// GetEventsByType returns events filtered by type
func (af *AuditFramework) GetEventsByType(eventType string) []AuditEvent {
	var filtered []AuditEvent
	for _, event := range af.events {
		if event.EventType == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GetEventsBySeverity returns events filtered by severity
func (af *AuditFramework) GetEventsBySeverity(severity string) []AuditEvent {
	var filtered []AuditEvent
	for _, event := range af.events {
		if event.Severity == severity {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%s", time.Now().Format("20060102150405"))
}

// generateReportID generates a unique report ID
func generateReportID() string {
	return fmt.Sprintf("rep_%s", time.Now().Format("20060102150405"))
}

// categorizeEvent categorizes an event type
func categorizeEvent(eventType string) string {
	switch eventType {
	case "policy_violation", "jailbreak_attempt", "injection_attack":
		return "policy"
	case "key_rotation", "encryption_failure", "decryption_failure":
		return "crypto"
	case "unauthorized_access", "access_denied":
		return "access"
	case "system_error", "component_failure":
		return "system"
	default:
		return "other"
	}
}

// getRemediationForEvent provides remediation guidance for an event type
func getRemediationForEvent(eventType string) string {
	switch eventType {
	case "policy_violation":
		return "Review and update security policies"
	case "jailbreak_attempt":
		return "Strengthen prompt validation and monitoring"
	case "injection_attack":
		return "Implement input sanitization and validation"
	case "key_rotation":
		return "Verify key rotation completed successfully"
	case "encryption_failure":
		return "Check KMS connectivity and key availability"
	case "decryption_failure":
		return "Verify key integrity and access permissions"
	case "unauthorized_access":
		return "Review access controls and user permissions"
	case "access_denied":
		return "Validate user authentication and authorization"
	default:
		return "Investigate and address the issue"
	}
}
