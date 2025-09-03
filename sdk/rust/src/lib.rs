use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum SentinelError {
    #[error("HTTP request failed: {0}")]
    Http(#[from] reqwest::Error),
    #[error("JSON serialization failed: {0}")]
    Json(#[from] serde_json::Error),
    #[error("Authentication failed")]
    Authentication,
    #[error("Threat detected: {0}")]
    ThreatDetected(String),
    #[error("Rate limit exceeded")]
    RateLimit,
    #[error("Invalid configuration: {0}")]
    InvalidConfig(String),
}

#[derive(Debug, Clone)]
pub struct SentinelClient {
    client: Client,
    base_url: String,
    api_key: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ThreatAnalysisRequest {
    pub prompt: String,
    pub user_id: Option<String>,
    pub session_id: Option<String>,
    pub context: Option<HashMap<String, serde_json::Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ThreatAnalysisResponse {
    pub threat_score: f64,
    pub is_safe: bool,
    pub threat_type: Option<String>,
    pub confidence: f64,
    pub explanation: String,
    pub recommendations: Vec<String>,
    pub request_id: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PolicyValidationRequest {
    pub content: String,
    pub policy_type: String,
    pub metadata: Option<HashMap<String, serde_json::Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PolicyValidationResponse {
    pub is_compliant: bool,
    pub violations: Vec<PolicyViolation>,
    pub score: f64,
    pub recommendations: Vec<String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PolicyViolation {
    pub rule_id: String,
    pub severity: String,
    pub description: String,
    pub suggestion: String,
}

impl SentinelClient {
    pub fn new(base_url: String, api_key: String) -> Result<Self, SentinelError> {
        if api_key.is_empty() {
            return Err(SentinelError::InvalidConfig("API key cannot be empty".to_string()));
        }

        let client = Client::new();
        
        Ok(SentinelClient {
            client,
            base_url,
            api_key,
        })
    }

    pub async fn analyze_threat(&self, request: ThreatAnalysisRequest) -> Result<ThreatAnalysisResponse, SentinelError> {
        let url = format!("{}/api/v1/analyze", self.base_url);
        
        let response = self.client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .header("Content-Type", "application/json")
            .json(&request)
            .send()
            .await?;

        if response.status() == 401 {
            return Err(SentinelError::Authentication);
        }

        if response.status() == 429 {
            return Err(SentinelError::RateLimit);
        }

        let analysis: ThreatAnalysisResponse = response.json().await?;
        
        if !analysis.is_safe {
            return Err(SentinelError::ThreatDetected(analysis.explanation.clone()));
        }

        Ok(analysis)
    }

    pub async fn validate_policy(&self, request: PolicyValidationRequest) -> Result<PolicyValidationResponse, SentinelError> {
        let url = format!("{}/api/v1/policy/validate", self.base_url);
        
        let response = self.client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.api_key))
            .header("Content-Type", "application/json")
            .json(&request)
            .send()
            .await?;

        if response.status() == 401 {
            return Err(SentinelError::Authentication);
        }

        Ok(response.json().await?)
    }

    pub async fn health_check(&self) -> Result<bool, SentinelError> {
        let url = format!("{}/health", self.base_url);
        
        let response = self.client
            .get(&url)
            .send()
            .await?;

        Ok(response.status().is_success())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_client_creation() {
        let client = SentinelClient::new(
            "https://api.sentinel.example.com".to_string(),
            "test-api-key".to_string(),
        );
        assert!(client.is_ok());
    }

    #[tokio::test]
    async fn test_empty_api_key() {
        let client = SentinelClient::new(
            "https://api.sentinel.example.com".to_string(),
            "".to_string(),
        );
        assert!(client.is_err());
    }
}
