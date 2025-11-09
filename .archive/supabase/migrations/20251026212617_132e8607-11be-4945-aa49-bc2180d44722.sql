-- Fix Security Definer View warnings by explicitly setting SECURITY INVOKER
-- and ensuring views don't bypass RLS

-- Recreate views with explicit SECURITY INVOKER
DROP VIEW IF EXISTS public.v_current_period_usage;
CREATE VIEW public.v_current_period_usage
WITH (security_invoker=true)
AS
SELECT
  s.organization_id,
  s.plan,
  s.status,
  s.current_period_start,
  s.current_period_end,
  COALESCE(u.total_calls, 0) as total_calls,
  COALESCE(u.included_calls_used, 0) as included_calls_used,
  COALESCE(u.overage_calls, 0) as overage_calls,
  s.included_calls,
  ROUND((COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) * 100, 2) as percent_used,
  COALESCE(u.estimated_charges_cents, 0) as estimated_charges_cents,
  EXTRACT(DAY FROM (s.current_period_end - NOW())) as days_remaining,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.70 as needs_70_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.90 as needs_90_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 1.00 as needs_100_alert
FROM public.subscriptions s
LEFT JOIN public.usage_counters u ON u.organization_id = s.organization_id
  AND u.period_start = s.current_period_start
WHERE s.status IN ('active', 'trialing', 'incomplete');

DROP VIEW IF EXISTS public.v_org_anomaly_summary;
CREATE VIEW public.v_org_anomaly_summary
WITH (security_invoker=true)
AS
SELECT
  ae.organization_id,
  s.current_period_start,
  s.current_period_end,
  COUNT(*) as total_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'critical') as critical_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'high') as high_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'medium') as medium_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'low') as low_anomalies,
  MAX(ae.detection_timestamp) as last_detection,
  MIN(ae.detection_timestamp) as first_detection_in_period
FROM public.anomaly_events ae
JOIN public.subscriptions s ON ae.organization_id = s.organization_id
WHERE ae.detection_timestamp >= s.current_period_start
  AND ae.detection_timestamp <= s.current_period_end
GROUP BY ae.organization_id, s.current_period_start, s.current_period_end;