export interface ParseResult {
  events: any[];
}

export class PayloadError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'PayloadError';
  }
}

// Parse JSON, { events: [] }, or ndjson strings into an events array.
export const parseEventsInput = (raw: string, maxEvents = 10_000): ParseResult => {
  const trimmed = raw.trim();
  if (!trimmed) {
    throw new PayloadError('No data provided. Paste JSON or upload a file.');
  }

  const tryJson = (): any => {
    try {
      return JSON.parse(trimmed);
    } catch {
      return null;
    }
  };

  // Try full JSON parse first
  const parsed = tryJson();
  let events: any[] | null = null;

  if (Array.isArray(parsed)) {
    events = parsed;
  } else if (parsed && Array.isArray(parsed.events)) {
    events = parsed.events;
  }

  // Fallback to ndjson (one JSON object per line)
  if (!events) {
    const lines = trimmed
      .split('\n')
      .map((l) => l.trim())
      .filter(Boolean);

    try {
      events = lines.map((line) => JSON.parse(line));
    } catch (e) {
      throw new PayloadError('Invalid JSON or NDJSON format.');
    }
  }

  if (!events || events.length === 0) {
    throw new PayloadError('No events found. Provide at least one event.');
  }

  if (events.length > maxEvents) {
    throw new PayloadError(`Too many events. Limit is ${maxEvents}.`);
  }

  return { events };
};

export interface DetectConfigOverride {
  baseline_size?: number;
  window_size?: number;
  hop_size?: number;
  ncd_threshold?: number;
  p_value_threshold?: number;
  compressor?: string;
}

export const sensitivityPresets: Record<
  'sensitive' | 'balanced' | 'strict',
  DetectConfigOverride
> = {
  sensitive: {
    baseline_size: 200,
    window_size: 30,
    hop_size: 10,
    ncd_threshold: 0.2,
    p_value_threshold: 0.1,
    compressor: 'zstd',
  },
  balanced: {
    baseline_size: 400,
    window_size: 50,
    hop_size: 10,
    ncd_threshold: 0.3,
    p_value_threshold: 0.05,
    compressor: 'zstd',
  },
  strict: {
    baseline_size: 800,
    window_size: 100,
    hop_size: 20,
    ncd_threshold: 0.45,
    p_value_threshold: 0.01,
    compressor: 'zstd',
  },
};

export const buildDetectPayload = (
  events: any[],
  streamId: string,
  sensitivity: keyof typeof sensitivityPresets,
  addIdempotency = true
) => {
  const payload: Record<string, unknown> = {
    stream_id: streamId || 'default',
    events: addIdempotency
      ? events.map((evt, idx) => ({
          ...evt,
          idempotency_key:
            evt.idempotency_key ||
            `evt_${Date.now()}_${Math.random().toString(16).slice(2)}_${idx}`,
        }))
      : events,
  };

  if (sensitivity !== 'balanced') {
    payload.config_override = sensitivityPresets[sensitivity];
  }

  return payload;
};
