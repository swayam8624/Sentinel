# Policy Packs for Sentinel + CipherMesh

## Overview

This document defines the initial policy packs for Sentinel + CipherMesh, covering default policies for PCI, PHI, PII, and other common data classes. These policies will be used as starting points for organizations to customize based on their specific requirements.

## Policy Structure

Each policy pack consists of:

1. **Data Classification Rules**: Define what constitutes each data class
2. **Detection Configuration**: Specify how to detect each data type
3. **Redaction Actions**: Define what actions to take when data is detected
4. **Detokenization Permissions**: Specify who can access original data
5. **Security Policies**: Define threat detection and response rules
6. **Audit Requirements**: Specify what events to log and monitor

## PCI Policy Pack

### Data Classification Rules

```yaml
pci:
  description: "Payment Card Industry Data"
  patterns:
    - name: "Primary Account Number (PAN)"
      regex: "(?:\\d[ -]*?){13,19}\\d"
      context: ["card", "credit", "debit", "payment", "transaction"]
    - name: "Card Verification Value (CVV)"
      regex: "\\d{3,4}"
      context: ["cvv", "cvc", "cid", "verification"]
    - name: "Expiration Date"
      regex: "(0[1-9]|1[0-2])/?([0-9]{4}|[0-9]{2})"
      context: ["expire", "expiration", "valid"]
    - name: "Cardholder Name"
      ner: true
      context: ["cardholder", "name on card"]
```

### Detection Configuration

```yaml
detection:
  confidence_threshold: 0.85
  context_window: 50 # characters before and after match
  languages: ["en"]
  enable_ocr: false
  code_secrets: false
```

### Redaction Actions

```yaml
actions:
  pan:
    action: "fpe"
    format: "XXXX-XXXX-XXXX-NNNN"
  cvv:
    action: "tokenize"
    ttl_hours: 24
  expiration_date:
    action: "mask"
    mask_char: "*"
  cardholder_name:
    action: "tokenize"
    ttl_hours: 168 # 1 week
```

### Detokenization Permissions

```yaml
detokenize:
  roles_allowed: ["pci_auditor", "security_admin"]
  justification_required: true
  session_timeout: 3600 # 1 hour
  audit_level: "detailed"
```

### Security Policies

```yaml
security:
  threat_detection:
    enabled: true
    sensitivity: "high"
  tool_guard:
    financial_tools_allowed: false
    payment_apis_allowed: false
  response_actions:
    - on_violation: "encrypt"
    - on_jailbreak: "block"
    - on_exfiltration_attempt: "encrypt_and_log"
```

### Audit Requirements

```yaml
audit:
  log_sensitive_access: true
  log_detokenization: true
  log_policy_violations: true
  retention_period: "730d" # 2 years
```

## PHI Policy Pack

### Data Classification Rules

```yaml
phi:
  description: "Protected Health Information"
  patterns:
    - name: "Medical Record Number"
      regex: "[A-Z]{2,3}\\d{6,10}"
      context: ["medical record", "patient id", "mrn"]
    - name: "Health Plan Beneficiary Number"
      regex: "\\d{8,12}"
      context: ["health plan", "beneficiary", "insurance"]
    - name: "Biometric Identifier"
      ner: true
      context: ["fingerprint", "retina", "iris", "voiceprint", "dna"]
    - name: "Full-face Photograph"
      ocr: true
      context: ["photo", "photograph", "picture", "image"]
```

### Detection Configuration

```yaml
detection:
  confidence_threshold: 0.90
  context_window: 100
  languages: ["en"]
  enable_ocr: true
  code_secrets: false
```

### Redaction Actions

```yaml
actions:
  medical_record_number:
    action: "tokenize"
    ttl_hours: 168 # 1 week
  health_plan_number:
    action: "tokenize"
    ttl_hours: 168
  biometric_identifier:
    action: "encrypt"
  full_face_photo:
    action: "block"
    reason: "Prohibited PHI type"
```

### Detokenization Permissions

```yaml
detokenize:
  roles_allowed: ["healthcare_provider", "medical_researcher", "hipaa_auditor"]
  justification_required: true
  session_timeout: 1800 # 30 minutes
  audit_level: "phi_comprehensive"
```

### Security Policies

```yaml
security:
  threat_detection:
    enabled: true
    sensitivity: "high"
  tool_guard:
    medical_devices_allowed: false
    health_records_apis_allowed: false
  response_actions:
    - on_violation: "encrypt"
    - on_jailbreak: "block"
    - on_exfiltration_attempt: "encrypt_and_alert"
```

### Audit Requirements

```yaml
audit:
  log_sensitive_access: true
  log_detokenization: true
  log_policy_violations: true
  retention_period: "3650d" # 10 years for HIPAA
```

## PII Policy Pack

### Data Classification Rules

```yaml
pii:
  description: "Personally Identifiable Information"
  patterns:
    - name: "Social Security Number"
      regex: "(?!000|666|9\\d{2})\\d{3}-(?!00)\\d{2}-(?!0000)\\d{4}"
      context: ["ssn", "social security", "tax id"]
    - name: "Driver's License Number"
      regex: "[A-Z]\\d{3,12}"
      context: ["driver", "license", "dl number"]
    - name: "Passport Number"
      regex: "[A-Z]\\d{8}"
      context: ["passport", "travel document"]
    - name: "Email Address"
      regex: "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
      context: ["email", "e-mail", "contact"]
    - name: "Phone Number"
      regex: "(\\+\\d{1,3}\\s?)?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}"
      context: ["phone", "telephone", "mobile", "cell"]
```

### Detection Configuration

```yaml
detection:
  confidence_threshold: 0.80
  context_window: 30
  languages: ["en", "es", "fr", "de"]
  enable_ocr: false
  code_secrets: false
```

### Redaction Actions

```yaml
actions:
  ssn:
    action: "fpe"
    format: "XXX-XX-NNNN"
  drivers_license:
    action: "tokenize"
    ttl_hours: 720 # 30 days
  passport_number:
    action: "tokenize"
    ttl_hours: 720
  email_address:
    action: "mask"
    mask_char: "*"
    preserve_domain: true
  phone_number:
    action: "fpe"
    format: "(XXX) XXX-NNNN"
```

### Detokenization Permissions

```yaml
detokenize:
  roles_allowed: ["analyst", "customer_support", "security_admin"]
  justification_required: true
  session_timeout: 7200 # 2 hours
  audit_level: "standard"
```

### Security Policies

```yaml
security:
  threat_detection:
    enabled: true
    sensitivity: "medium"
  tool_guard:
    pii_tools_allowed: true
    restricted_tools: ["social_media_post", "public_api"]
  response_actions:
    - on_violation: "encrypt"
    - on_jailbreak: "reframe"
    - on_exfiltration_attempt: "encrypt_and_log"
```

### Audit Requirements

```yaml
audit:
  log_sensitive_access: true
  log_detokenization: true
  log_policy_violations: true
  retention_period: "365d" # 1 year
```

## Credentials Policy Pack

### Data Classification Rules

```yaml
credentials:
  description: "Authentication Credentials"
  patterns:
    - name: "Password"
      entropy: true
      min_length: 8
      context: ["password", "passwd", "pwd", "secret"]
    - name: "API Key"
      regex: "^(sk|pk)_[a-zA-Z0-9]{32,}$"
      context: ["api key", "apikey", "secret key"]
    - name: "SSH Key"
      regex: "^(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ecdsa-sha2-nistp384|ecdsa-sha2-nistp521|ssh-ed25519) "
      context: ["ssh", "private key", "public key"]
    - name: "Database Credentials"
      ner: true
      context: ["database", "db password", "connection string"]
```

### Detection Configuration

```yaml
detection:
  confidence_threshold: 0.95
  context_window: 20
  languages: ["en"]
  enable_ocr: false
  code_secrets: true
```

### Redaction Actions

```yaml
actions:
  password:
    action: "tokenize"
    ttl_hours: 1 # 1 hour
  api_key:
    action: "tokenize"
    ttl_hours: 24
  ssh_key:
    action: "encrypt"
  database_credentials:
    action: "tokenize"
    ttl_hours: 168 # 1 week
```

### Detokenization Permissions

```yaml
detokenize:
  roles_allowed: ["security_admin", "devops_engineer"]
  justification_required: true
  session_timeout: 900 # 15 minutes
  audit_level: "detailed"
```

### Security Policies

```yaml
security:
  threat_detection:
    enabled: true
    sensitivity: "critical"
  tool_guard:
    credential_tools_allowed: false
  response_actions:
    - on_violation: "encrypt"
    - on_jailbreak: "block"
    - on_exfiltration_attempt: "encrypt_alert_and_notify"
```

### Audit Requirements

```yaml
audit:
  log_sensitive_access: true
  log_detokenization: true
  log_policy_violations: true
  retention_period: "1825d" # 5 years
```

## Default Policy Pack

### Data Classification Rules

```yaml
default:
  description: "Default policies for general data protection"
  patterns:
    - name: "Generic Sensitive Data"
      entropy: true
      min_length: 10
      context: ["sensitive", "confidential", "private"]
```

### Detection Configuration

```yaml
detection:
  confidence_threshold: 0.70
  context_window: 25
  languages: ["en"]
  enable_ocr: false
  code_secrets: true
```

### Redaction Actions

```yaml
actions:
  generic_sensitive:
    action: "mask"
    mask_char: "*"
```

### Detokenization Permissions

```yaml
detokenize:
  roles_allowed: ["admin"]
  justification_required: true
  session_timeout: 3600 # 1 hour
  audit_level: "standard"
```

### Security Policies

```yaml
security:
  threat_detection:
    enabled: true
    sensitivity: "low"
  tool_guard:
    default_tools_allowed: true
  response_actions:
    - on_violation: "reframe"
    - on_jailbreak: "reframe"
    - on_exfiltration_attempt: "log"
```

### Audit Requirements

```yaml
audit:
  log_sensitive_access: true
  log_detokenization: true
  log_policy_violations: true
  retention_period: "365d" # 1 year
```

## Policy Pack Implementation

### Policy Versioning

All policy packs will follow semantic versioning:

- MAJOR version when you make incompatible policy changes
- MINOR version when you add functionality in a backward compatible manner
- PATCH version when you make backward compatible bug fixes

### Policy Deployment

Policy packs can be deployed in three modes:

1. **Audit Mode**: Test policies without enforcement
2. **Shadow Mode**: Run policies alongside existing ones for comparison
3. **Enforce Mode**: Full policy enforcement

### Policy Customization

Organizations can customize policy packs by:

1. Modifying detection thresholds
2. Adding/removing data patterns
3. Changing redaction actions
4. Adjusting permissions
5. Customizing security responses

## Testing and Validation

Each policy pack will be tested against:

1. Standard test datasets for each data class
2. Common bypass techniques
3. Performance impact measurements
4. False positive/negative analysis
5. Compliance requirement validation
