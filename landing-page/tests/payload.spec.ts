import { test, expect } from '@playwright/test';
import { parseEventsInput, PayloadError } from '../src/utils/payload';

test('parses array of events', () => {
  const raw = JSON.stringify([{ body: { a: 1 } }, { body: { b: 2 } }]);
  const result = parseEventsInput(raw);
  expect(result.events.length).toBe(2);
});

test('parses object with events property', () => {
  const raw = JSON.stringify({ events: [{ body: { a: 1 } }] });
  const result = parseEventsInput(raw);
  expect(result.events[0].body.a).toBe(1);
});

test('parses ndjson', () => {
  const raw = '{"body":{"a":1}}\n{"body":{"b":2}}';
  const result = parseEventsInput(raw);
  expect(result.events[1].body.b).toBe(2);
});

test('throws on too many events', () => {
  const events = Array.from({ length: 10_001 }, () => ({ body: {} }));
  const raw = JSON.stringify(events);
  expect(() => parseEventsInput(raw)).toThrow(PayloadError);
});
