from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, Optional

import requests


@dataclass
class DriftlockError(Exception):
  message: str
  status: Optional[int] = None
  details: Any = None

  def __post_init__(self) -> None:
    super().__init__(self.message)

  def __str__(self) -> str:
    status = f" (status={self.status})" if self.status is not None else ""
    return f"{self.message}{status}"


class DriftlockClient:
  """Lightweight wrapper around Driftlock's /v1 API."""

  def __init__(
      self,
      api_key: str,
      base_url: str = "http://localhost:8080",
      timeout: float = 10.0,
      session: Optional[requests.Session] = None) -> None:
    if not api_key:
      raise DriftlockError("api_key is required")
    self.api_key = api_key
    self.base_url = base_url.rstrip("/")
    self.timeout = timeout
    self.session = session or requests.Session()

  def health(self) -> Dict[str, Any]:
    return self._request("GET", "/healthz")

  def detect(self, payload: Dict[str, Any]) -> Dict[str, Any]:
    return self._request("POST", "/v1/detect", json=payload)

  def list_anomalies(self, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    return self._request("GET", "/v1/anomalies", params=params or {})

  def get_anomaly(self, anomaly_id: str) -> Dict[str, Any]:
    if not anomaly_id:
      raise DriftlockError("anomaly_id is required")
    return self._request("GET", f"/v1/anomalies/{anomaly_id}")

  def _request(
      self,
      method: str,
      path: str,
      json: Optional[Dict[str, Any]] = None,
      params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    url = f"{self.base_url}{path}"
    headers = {"X-Api-Key": self.api_key}
    if json is not None:
      headers["Content-Type"] = "application/json"

    response = self.session.request(
        method=method,
        url=url,
        json=json,
        params=params,
        headers=headers,
        timeout=self.timeout,
    )

    content_type = response.headers.get("content-type", "")
    is_json = "application/json" in content_type

    try:
      payload = response.json() if is_json else response.text
    except Exception:  # pylint: disable=broad-except
      payload = response.text

    if not response.ok:
      raise DriftlockError(
          f"Request failed with status {response.status_code}",
          status=response.status_code,
          details=payload)

    if isinstance(payload, str):
      raise DriftlockError("Expected JSON response", details=payload)

    return payload
