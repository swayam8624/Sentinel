# Sentinel + CipherMesh

## Software Requirements Specification (SRS) **and** End-to-End Workflow Plan

> Goal: Let enterprises use any LLM provider without exposing raw data, while adding a self-healing security layer that detects, corrects, or cryptographically contains adversarial prompts—drop-in, provider-agnostic, production-grade.

---

## 0. Document Control

- **Product name:** Sentinel (self-healing LLM firewall) + CipherMesh (PII/tokenization & crypto layer)
- **Status:** Draft v1.0 (Foundational SRS)
- **Audience:** Product, Security, ML Platform, SWE, DevSecOps, Compliance, Infra
- **Change policy:** Semantic versioning for SRS; major bumps when public APIs/policies change.

---

## 1. Introduction

### 1.1 Purpose

Define functional and non-functional requirements for a gateway/SDK that:

- Shields upstream LLM providers from raw sensitive data (pre-prompt tokenization / FPE).
- Detects jailbreaks & injections; self-reflects & rewrites; or **encrypts/locks** unsafe outputs.
- Works with **any** model/vendor via adapters/proxy/sidecar.
- Provides **tamper-evident** audit, policy versioning, and BYOK/HSM key management.

### 1.2 Scope

- In-line, real-time data redaction/tokenization + reversible detokenization.
- Self-healing security pipeline (detect → reflect → rewrite → encrypt/log).
- Tool/function-call guarding and streaming cutover.
- Admin console for policies/keys/logs; observability; multi-tenancy.

### 1.3 Definitions & Acronyms

- **FPE**: Format-Preserving Encryption (FF3-1).
- **HKDF**: RFC5869 key derivation.
- **AEAD**: Authenticated Encryption w/ Associated Data (AES-256-GCM).
- **BYOK**: Bring Your Own Key (KMS/HSM).
- **OPA**: Open Policy Agent policies.
- **FAISS/Chroma**: Vector index stores.
- **PII/PHI/PCI**: Sensitive data classes.

### 1.4 References (standards & guidance)

- AES-GCM (NIST SP 800-38D), FF3-1 (NIST SP 800-38G), HKDF (RFC 5869), OWASP ASVS/Top-10, ISO/IEC 27001, GDPR/CCPA/HIPAA (where applicable).

---

## 2. Overall Description

### 2.1 Product perspective

A drop-in **gateway/proxy** (and SDK) that sits between applications and LLM providers. It can also run as a **sidecar** for microservices.

### 2.2 Product functions (high level)

- **CipherMesh**: detect PII/secrets → tokenize/FPE → provider → detokenize (policy-gated).
- **Sentinel**: context trace, violation detection, reflection, safe rewriting, tool guard, **encryption on violation**.
- **Policy Engine**: OPA-style rules, versioned & tenant-scoped.
- **Key Management**: BYOK/HSM; envelope encryption; rotation.
- **Audit/Observability**: Tamper-evident logs; metrics; alerts.

### 2.3 User classes

- **Integrators** (backend/ML engineers).
- **Security/Admin** (policies/keys/alerts).
- **Compliance** (audit/exports).
- **Developers** (SDK usage).

### 2.4 Constraints

- Added latency budgets on hot path.
- Provider quirks (streaming/function calls/rate limits).
- Multilingual, multimodal inputs.
- Regulatory demands (data residency, DPA).

### 2.5 Assumptions/Dependencies

- Tenant BYOK via cloud KMS (preferred).
- Reasonable compute for NER/embeddings/OCR where enabled.
- Network egress to providers (proxy mode).

---

## 3. System Features (Functional Requirements)

> Use IDs for traceability: **FR-xxx** (functional), **NFR-xxx** (non-functional).

### 3.1 Adapters & Integration

- **FR-101**: Provide adapters for OpenAI/Anthropic/Mistral/HF/Ollama at launch.
- **FR-102**: Reverse-proxy mode exposing provider-compatible endpoints.
- **FR-103**: SDK middleware (Python/Node) wrapping any LLM client.
- **FR-104**: Sidecar/gateway (gRPC/HTTP) for polyglot stacks.
- **FR-105**: Support streaming; mid-stream inspection & cutover.

### 3.2 Context Tracking & Policy

- **FR-110**: Maintain conversation/system fingerprint & short-term memory window.
- **FR-111**: OPA-style policy evaluation for data classes and actions.
- **FR-112**: Policy versioning; per-tenant pinning; canary/shadow evaluation.

### 3.3 CipherMesh (Detection & Redaction)

- **FR-120**: Detect PII/PHI/PCI via hybrid rules: regex + multilingual NER + entropy-based secret scanners.
- **FR-121**: Normalize inputs (Unicode NFKC, zero-width char stripping, homoglyph mapping).
- **FR-122**: Optional OCR for image/PDF attachments (configurable).
- **FR-123**: Actions by class: tokenize, FPE (FF3-1), mask, hash, drop, allow.
- **FR-124**: **Format-preserving** placeholders for model utility; semantic placeholders when needed.
- **FR-125**: Redact prompts, function-call args, tool I/O, RAG chunks, and logs.
- **FR-126**: Maintain reversible token map in encrypted vault; TTL by policy.
- **FR-127**: Streaming chunk redaction; function-call JSON field-level redaction.
- **FR-128**: Policy-gated detokenization on egress (role/RBAC aware).

### 3.4 Sentinel (Self-Healing Security)

- **FR-140**: Semantic violation detector (embeddings + signatures + rules).
- **FR-141**: Signature store (FAISS/Chroma) with online updates.
- **FR-142**: Self-alignment reflection pass (constitutional prompt).
- **FR-143**: Prompt reconstructor (multi-candidate rewrite + ranking).
- **FR-144**: Router: allow / reframe / encrypt / block decisions.
- **FR-145**: ToolGuard: disable downstream tools on violation.
- **FR-146**: Mid-stream cut (≤150 ms) on violation detection.

### 3.5 Cryptography & Vault

- **FR-160**: Envelope encryption with tenant BYOK (KMS/HSM).
- **FR-161**: Key derivation per message using **HKDF-SHA-512** with per-message random salt; include token_count as associated info (not entropy).
- **FR-162**: AEAD via AES-256-GCM; unique nonce per (key,message); nonce management enforced.
- **FR-163**: FF3-1 FPE with per-field tweaks; rotation supported.
- **FR-164**: Split-knowledge for vault master keys (optional), strict RBAC.
- **FR-165**: Zero plaintext storage in logs; PII scrubbing in traces.

### 3.6 Audit, Alerts, Governance

- **FR-180**: Tamper-evident logs: hash chain per event; daily Merkle root.
- **FR-181**: Emit metrics & events (OTel) with privacy-safe fields.
- **FR-182**: Admin alerts (webhook/email/Slack) on encrypt/block decisions.
- **FR-183**: SIEM export; scoped search; retention policies.
- **FR-184**: DSR/DSAR support: locate & purge tenant data/maps by keys.

### 3.7 Admin Console & APIs

- **FR-200**: Web console: policies, tenants, adapters, keys, alerts, dashboards.
- **FR-201**: Admin APIs: policy CRUD/versioning, key mgmt (BYOK hooks), logs export.
- **FR-202**: Read-only viewer role; per-tenant isolation.

### 3.8 Multi-Tenancy & Residency

- **FR-220**: Hard isolation of token vault & vectors per tenant.
- **FR-221**: Data residency controls (pin vault & compute region).
- **FR-222**: Rate limits & quotas per tenant.

---

## 4. External Interface Requirements

### 4.1 Gateway (proxy) endpoints (examples)

- `POST /v1/chat/completions` (OpenAI-compatible)
  Headers: `X-Tenant`, `X-Policy-Version` (optional)
  Request body proxied after redaction; response post-processed (detokenize/guard).

- `POST /sentinel/admin/policies` (admin)

- `GET /sentinel/admin/logs?from=...&to=...` (admin)

### 4.2 SDK Interface (provider-agnostic)

```python
resp = sentinel.chat(
  model="gpt-4o",
  messages=[{"role":"user","content":"..."}],
  stream=False,
  metadata={"data_classes":["pci"], "tools_allowed":True}
)
```

### 4.3 Config-as-Policy (YAML excerpt)

```yaml
sentinel:
  mode: enforce # audit | enforce | silent
  thresholds:
    violation_similarity: 0.78
    reflect_confidence: 0.65
  rewriter:
    enabled: true
    confirm_user: true
  encryption:
    enabled: true
    base_secret_env: SENTINEL_SECRET
  tools:
    lockdown_on_violation: true

ciphermesh:
  detectors:
    languages: [en, es, fr, de, ja, hi]
    enable_ocr: false
    code_secrets: true
  actions:
    pii: fpe
    credentials: tokenize
    generic: mask
  detokenize:
    roles_allowed: ["analyst", "ops"]
```

---

## 5. System Models (Diagrams)

### 5.1 Component Diagram

```
[Client] -> [Sentinel Gateway]
  ├─ CipherMesh (NER/regex/secret-scan/OCR)
  ├─ Sentinel (detector/reflector/rewriter/router/toolguard)
  ├─ Adapters (OpenAI/Anthropic/...)
  ├─ Policy Engine (OPA)
  ├─ Crypto/KMS Vault (FF3-1, AEAD, HKDF)
  └─ Audit/Observability (OTel, Merkle logs)
-> [LLM Provider(s)]
```

### 5.2 Sequence (unsafe prompt)

```
Client → Gateway: prompt
Gateway → CipherMesh: redact
Gateway → Sentinel.Detector: score=0.91 (violation)
Sentinel.Router: path=encrypt
Sentinel.Encryptor: AES-GCM(key=HKDF(..., info=token_count))
Audit.Log: event(hash), ciphertext, metadata
Gateway → Client: {status:"encrypted", token_count, log_id}
```

### 5.3 Sequence (safe with rewriting)

```
Client → Gateway: prompt
CipherMesh: redact
Detector: suspicious
Reflector: misaligned
Rewriter: safe_prompt_candidate
Client ↔ Gateway: confirm
Gateway → Provider: redacted safe prompt
Provider → Gateway: response
Detokenize (policy OK)
Gateway → Client: final
```

---

## 6. Non-Functional Requirements

### Performance & Availability

- **NFR-001**: Added latency p95 (non-stream) ≤ **700 ms**; p50 ≤ 300 ms.
- **NFR-002**: Streaming redaction overhead per chunk ≤ **200 ms**.
- **NFR-003**: Uptime (gateway) **99.9%**; graceful degradation if OCR off.

### Security & Privacy

- **NFR-010**: Provider-bound **raw PII leakage ≈ 0%** (measured).
- **NFR-011**: Plaintext unsafe output leakage when flagged = **0%**.
- **NFR-012**: BYOK/HSM enforced; no KEK at rest; strict RBAC.
- **NFR-013**: Unique nonce per message; HKDF with random salt per message.

### Reliability & DR

- **NFR-020**: Vault RPO ≤ 5 min; RTO ≤ 30 min.
- **NFR-021**: Hash-chained logs; daily Merkle anchor.

### Scalability & Maintainability

- **NFR-030**: Horizontal autoscale to ≥ **200 req/s** per pod baseline.
- **NFR-031**: SemVer for adapters; policy backward compatibility.

### Internationalization

- **NFR-040**: Multilingual NER for listed languages; Unicode canonicalization.

---

## 7. Data Design

### 7.1 Core Schemas (conceptual)

- **tenants**(id, name, region, kms_key_arn, policy_version_pin, …)
- **policies**(id, version, ruleset_json, status, created_at)
- **redaction_maps**(tenant_id, conversation_id, token_id, enc_value_aead, fpe_tweak, created_at, ttl)
- **events**(tenant_id, event_id, ts, decision, scores, token_count, ciphertext_len, prev_hash, hash)
- **signatures**(tenant_id, vector_id, embedding, label, created_at)

### 7.2 Token Vault

- AEAD-encrypted values, indexed by token_id; per-tenant envelope keys.
- TTL & purge tasks; access reason codes; full audit.

---

## 8. Threat Model (STRIDE overview)

- **Spoofing**: mTLS to providers/internal; JWT for console.
- **Tampering**: AEAD everywhere; Merkle logs; WORM storage (optional).
- **Repudiation**: Signed logs; operator identity tagging.
- **Information disclosure**: Redaction pre-provider; policy-gated detokenization.
- **DoS**: Rate limits; circuit breakers; model/provider fallback.
- **Escalation**: RBAC; least privilege; break-glass tracking.

---

## 9. Verification & Validation

### Security & Privacy Metrics

- Provider-bound PII leakage: **≈0%** across languages/modalities.
- Final output sensitive leakage: **<0.5%**.
- Jailbreak block uplift vs baseline: **≥90%** (AdvBench/JBB/custom).
- Mid-stream cutover: **≤150 ms**.

### Utility & Quality

- Task success retention vs unmasked:

  - **≥95%** generic tasks,
  - **≥85%** PII-heavy (human eval + semantic similarity).

- RAG answer fidelity with masked sources: **≥90%** semantic match.

### Performance

- p95 ≤ 700 ms; throughput ≥ 200 req/s per pod in enforce mode.

### Governance

- 100% events logged with hash chain; daily anchored root.

---

## 10. Test Plan (summary)

- **Unit**: detectors, crypto, policy engine, adapters.
- **Integration**: end-to-end redaction→provider→detokenize; function-call & tools; streaming.
- **Security**: red-team prompts; obfuscation (unicode, base64, homoglyphs); code secret scans; OCR bypass attempts.
- **Perf/Load**: concurrency, backpressure, autoscale.
- **Chaos/DR**: KMS outages, nonce collision tests (must never pass), vault failover.

---

## 11. Deployment & Operations

- **Packaging**: Docker images for gateway/sidecar; Helm charts; SDKs (pip/npm).
- **Secrets**: KMS-managed; no plaintext keys at rest.
- **Observability**: Metrics (latency, leak attempts, block rate), traces, structured logs (PII-safe).
- **Runbooks**: policy rollout (audit→enforce), key rotation, vault recovery, incident response.
- **Versioning**: SemVer; policy pinning per tenant; shadow eval before enforce.

---

## 12. Roadmap (feature gates)

- v1.0: Proxy + SDK, redaction (text), detectors (NER/regex/secrets), FF3-1, vault, Sentinel core (detect/reflect/rewrite/encrypt), ToolGuard, policy v1, audit logs, metrics, alerts.
- v1.1: Streaming redaction/cutover, function-call deep masking, RAG source scanning.
- v1.2: OCR (images/PDF), multilingual expansion, admin console v1, SIEM exporters.
- v1.3: Signature learner pipeline, confidential compute option, DP toggles.

---

# End-to-End Workflow Structure (All Phases & Tasks)

> Designed so you can ship MVP, enforce without downtime, and iterate with versioned policies.

## Phase A — Inception & Threat Modeling

- **A1** Define use cases, data classes, residency.
- **A2** STRIDE threat model; abuse cases (roleplay, obfuscation, tool exfil).
- **A3** SLOs & metrics (latency, leakage, block rate); acceptance gates.
- **A4** Draft policy packs (PCI/PHI/PII defaults).

**Deliverables:** Threat model doc, acceptance metric sheet, initial policies.

---

## Phase B — Architecture & Foundations

- **B1** Choose integration modes (proxy/SDK/sidecar) & provider adapters.
- **B2** Crypto design: HKDF-SHA-512 per-message, AES-GCM nonce policy, FF3-1 FPE domain/tweaks.
- **B3** Data model for vault, events, policies, signatures.
- **B4** OPA policy engine scaffolding & versioning semantics.

**Deliverables:** ADRs, schemas, policy engine MVP.

---

## Phase C — CipherMesh Core

- **C1** Detectors: regex, multilingual NER, secret scanners, canonicalization.
- **C2** Redaction actions & reversible tokenization; vault AEAD.
- **C3** FPE (FF3-1) library & per-field tweaks; test domains (credit card, phone, IDs).
- **C4** Detokenization gate w/ RBAC; audit entries.
- **C5** Streaming chunk redaction scaffolding.

**Deliverables:** Redaction/Detokenization library + tests; vault service.

---

## Phase D — Sentinel Core

- **D1** Violation detector (embeddings + rules + signature index).
- **D2** Reflection pass (constitutional prompts).
- **D3** Rewriter (multi-candidate, ranking, user confirm flow).
- **D4** Router & ToolGuard integration; function-call guarding.
- **D5** Mid-stream cutover logic.

**Deliverables:** Sentinel pipeline module with decision API.

---

## Phase E — Adapters & Proxy

- **E1** OpenAI/Anthropic/Mistral/HF/Ollama adapters.
- **E2** Reverse proxy w/ provider-compatible endpoints.
- **E3** SDKs (Python/Node) with minimal surface.
- **E4** Streaming support end-to-end; backpressure & timeouts.

**Deliverables:** Running gateway; SDKs; adapter conformance tests.

---

## Phase F — Security & Crypto Hardening

- **F1** HKDF salts and nonce uniqueness enforcement.
- **F2** KMS/HSM envelope integration; BYOK flows; rotation runbook.
- **F3** Vault split-knowledge (optional) & access reason codes.
- **F4** Merkle hash chain logs; daily root anchoring (e.g., external timestamping).

**Deliverables:** Crypto compliance checklist; penetration & misuse tests.

---

## Phase G — Observability & Admin

- **G1** OTel metrics/traces; PII-safe logging; SIEM exporters.
- **G2** Admin console (v1): policies, tenants, alerts, dashboards.
- **G3** Policy versioning UX: audit → canary → enforce; per-tenant pinning.
- **G4** DSAR tools: locate/purge mappings & logs.

**Deliverables:** Dashboards; admin API/UI; ops runbooks.

---

## Phase H — MVP Rollout (Audit Mode)

- **H1** Deploy proxy in **audit** mode; gather leakage/block metrics.
- **H2** Calibrate thresholds; tune detectors & policies.
- **H3** Load & chaos tests; SLO validation.

**Gate:** Provider-bound PII leakage ≈ 0%; p95 latency within budget.

---

## Phase I — Enforce Mode & Tooling

- **I1** Switch to **enforce** for select tenants (canary).
- **I2** Enable ToolGuard; function-call deep masking; streaming cutover.
- **I3** RAG source scanning; partial detokenization policies.

**Gate:** Jailbreak block uplift ≥ 90%; false positives < 2%.

---

## Phase J — Enterprise Readiness

- **J1** Residency controls; multi-region deployment.
- **J2** SOC2-style audit artifacts; DPA templates.
- **J3** Policy version migrations with graceful fallback; adapter SemVer.

**Gate:** Multi-tenant isolation verified; zero-downtime rotations & policy swaps.

---

## Phase K — Advanced (optional)

- **K1** Signature learner retraining pipeline.
- **K2** OCR enablement for images/PDF; multimodal detectors.
- **K3** Confidential compute (Nitro/SEV-SNP) for vault operations.
- **K4** Differential privacy flags for analytics exports.

---

## Acceptance Criteria (Go/No-Go for Production)

- **Security**: 0% raw PII to providers; 0% plaintext unsafe leakage on violations.
- **Robustness**: ≥90% jailbreak block uplift vs baseline; mid-stream cutover ≤150 ms.
- **Utility**: ≥95% task success retention generic; ≥85% PII-heavy.
- **Performance**: p95 ≤700 ms non-stream; streaming overhead ≤200 ms/chunk.
- **Ops**: 99.9% uptime; full audit hashing; policy versioning live; BYOK operational.

---

### Final Note

This SRS + workflow is built to be **auditable, enforceable, and evolvable**. Your edge is the **integration** of (1) **high-recall redaction** that preserves utility, (2) **self-healing security cognition**, and (3) **cryptographic containment**—all delivered as a **provider-agnostic plugin/gateway** with policy versioning and hard SLOs.
