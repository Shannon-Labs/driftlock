//! API Key validation

use uuid::Uuid;

/// Parse API key format: dlk_{uuid}.{secret}
pub fn parse_api_key(key: &str) -> Option<(Uuid, String)> {
    let key = key.strip_prefix("dlk_")?;
    let parts: Vec<&str> = key.splitn(2, '.').collect();
    if parts.len() != 2 {
        return None;
    }
    let id = Uuid::parse_str(parts[0]).ok()?;
    Some((id, parts[1].to_string()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_api_key() {
        let key = "dlk_12345678-1234-1234-1234-123456789012.abc123secret";
        let result = parse_api_key(key);
        assert!(result.is_some());
        let (id, secret) = result.unwrap();
        assert_eq!(id.to_string(), "12345678-1234-1234-1234-123456789012");
        assert_eq!(secret, "abc123secret");
    }
}
