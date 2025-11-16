/**
 * Cloudflare Pages Function handling contact submissions from the landing page.
 *
 * Behaviour:
 * - Accepts POSTed JSON payloads from the CTA form.
 * - Forwards the payload to CRM_WEBHOOK_URL when configured.
 * - Otherwise logs the submission so we have an audit trail without failing the UX.
 */

interface Env {
  CRM_WEBHOOK_URL?: string;
  CONTACT_LOG_KV?: KVNamespace; // optional KV binding for persistence (configure later)
}

interface ContactPayload {
  name: string;
  email: string;
  company?: string;
  message?: string;
  [key: string]: unknown;
}

const corsHeaders: HeadersInit = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'POST, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With',
  'Access-Control-Max-Age': '86400',
  'Content-Type': 'application/json',
};

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const badRequest = (message: string): Response =>
  new Response(JSON.stringify({ ok: false, error: message }), {
    status: 400,
    headers: corsHeaders,
  });

export const onRequestOptions = () =>
  new Response(null, { status: 204, headers: corsHeaders });

export const onRequestPost = async ({ request, env }: { request: Request; env: Env }) => {
  let payload: ContactPayload;

  try {
    payload = (await request.json()) as ContactPayload;
  } catch {
    return badRequest('Invalid JSON body');
  }

  if (!payload.name || payload.name.trim().length < 2) {
    return badRequest('Name is required');
  }

  if (!payload.email || !emailRegex.test(payload.email)) {
    return badRequest('Valid email is required');
  }

  const submission = {
    ...payload,
    name: payload.name.trim(),
    email: payload.email.trim(),
    receivedAt: new Date().toISOString(),
    userAgent: request.headers.get('user-agent') ?? 'unknown',
  };

  try {
    if (env.CRM_WEBHOOK_URL) {
      const upstream = await fetch(env.CRM_WEBHOOK_URL, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify(submission),
      });

      if (!upstream.ok) {
        console.error('CRM webhook failed', upstream.status, await upstream.text());
        return new Response(JSON.stringify({ ok: false, error: 'Upstream_error' }), {
          status: 502,
          headers: corsHeaders,
        });
      }

      return new Response(JSON.stringify({ ok: true, mode: 'webhook' }), {
        status: 200,
        headers: corsHeaders,
      });
    }

    if (env.CONTACT_LOG_KV) {
      const key = `contact:${submission.receivedAt}:${crypto.randomUUID()}`;
      await env.CONTACT_LOG_KV.put(key, JSON.stringify(submission));
    } else {
      console.log('Contact submission (no webhook configured)', submission);
    }

    return new Response(JSON.stringify({ ok: true, mode: 'logged' }), {
      status: 200,
      headers: corsHeaders,
    });
  } catch (error) {
    console.error('Contact handler error', error);
    return new Response(JSON.stringify({ ok: false, error: 'handler_error' }), {
      status: 500,
      headers: corsHeaders,
    });
  }
};