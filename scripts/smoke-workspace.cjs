'use strict';
// Real HTTP contract against the built Compose stack; no provider delivery implied.
const assert = require('node:assert/strict');
const base = process.env.WORKSPACE_BASE_URL || 'http://127.0.0.1:8080';
const token = process.env.PUBLISHING_API_TOKEN;
if (!token) throw new Error('PUBLISHING_API_TOKEN is required');
async function request(route, method = 'GET', body, authenticated = true) {
  return fetch(base + route, {
    method,
    headers: { ...(authenticated ? { Authorization: 'Bearer ' + token } : {}), 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
    signal: AbortSignal.timeout(5000),
  });
}
async function main() {
  let healthy = false;
  for (let attempt = 0; attempt < 30; attempt++) {
    try {
      const response = await request('/api/health', 'GET', undefined, false);
      if (response.ok && (await response.json()).status === 'ok') { healthy = true; break; }
    } catch { /* allow the freshly built stack to start */ }
    await new Promise(resolve => setTimeout(resolve, 2000));
  }
  assert.ok(healthy, 'API must become healthy');
  const web = await request('/', 'GET', undefined, false);
  assert.equal(web.status, 200);
  assert.match(await web.text(), /<div id="root"><\/div>/);
  assert.equal((await request('/api/publications', 'GET', undefined, false)).status, 401);
  const created = await request('/api/publications', 'POST', { title: 'Container smoke', content: 'Internal editorial lifecycle smoke test.' });
  assert.equal(created.status, 201);
  const id = (await created.json()).id;
  for (const [action, expected, body] of [
    ['submit', 'in_review'], ['approve', 'approved'],
    ['schedule', 'scheduled', { scheduledFor: new Date(Date.now() + 3600000).toISOString() }],
    ['publish', 'published'], ['archive', 'archived'],
  ]) {
    const response = await request('/api/publications/' + id + '/' + action, 'POST', body);
    assert.equal(response.status, 200, action);
    assert.equal((await response.json()).status, expected, action);
  }
  console.log('Web, authentication, and internal editorial lifecycle smoke passed');
}
main().catch(error => { console.error(error.message); process.exitCode = 1; });
