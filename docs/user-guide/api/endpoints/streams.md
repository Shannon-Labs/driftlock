# Stream Profiles & Tuning

Manage detection sensitivity and review tuning history per stream.

## GET /v1/streams/{id}/profile

Returns the active detection profile, thresholds, and window sizes.

```
GET https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/streams/{id}/profile
```

Response fields:
- `profile`: one of `sensitive`, `balanced`, `strict`, `custom`
- `auto_tune_enabled`: whether feedback-driven auto-tuning is active
- `adaptive_window_enabled`: whether adaptive window sizing is enabled
- `current_thresholds`: includes `ncd_threshold`, `pvalue_threshold`, `baseline_size`, `window_size`

## PATCH /v1/streams/{id}/profile

Update the detection profile or enable/disable auto-tuning or adaptive windows.

Body (any subset):
```json
{
  "profile": "strict",
  "auto_tune_enabled": true,
  "adaptive_window_enabled": true
}
```

## GET /v1/streams/{id}/tuning

View current settings, recent tune history, and feedback stats (last 30 days).

```
GET https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/streams/{id}/tuning
```

Returns:
- `current_settings`: same shape as the profile response
- `tune_history`: recent threshold adjustments
- `feedback_stats`: counts and false-positive rate

## GET /v1/profiles

List available detection profiles and their defaults.

```
GET https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/profiles
```

