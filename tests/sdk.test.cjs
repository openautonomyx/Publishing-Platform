'use strict';
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const sdk = require('../dist');

test('packed entrypoint loads and exposes the documented client and validators', () => {
  assert.equal(sdk.default, sdk.OpenAutonomyX);
  assert.throws(() => sdk.createOpenAutonomyX({apiKey: ''}), /API key is required/);
  assert.equal(sdk.validateEmail('author@example.test'), true);
  assert.equal(sdk.validateEmail('invalid'), false);
  assert.throws(() => sdk.validateContentRequest({title:'',content:'long enough content'}), sdk.ValidationError);
  assert.doesNotThrow(() => sdk.validateLanguages(['en','zip']));
  assert.throws(() => sdk.validateLanguages(['invalid']), sdk.ValidationError);
});

test('legacy SDK sends auth and content contracts to an explicitly configured HTTP backend', async t => {
  const seen = [];
  const server = http.createServer(async (req,res) => {
    let body = '';
    for await (const chunk of req) body += chunk;
    seen.push({method:req.method,path:req.url,auth:req.headers.authorization,body:body ? JSON.parse(body) : undefined});
    res.setHeader('Content-Type','application/json');
    if (req.url === '/api/v1/publications/missing') {res.writeHead(404);res.end(JSON.stringify({message:'Missing publication'}));return;}
    res.end(JSON.stringify({id:'publication-1',title:'Contract title',content:'Contract content'}));
  });
  await new Promise(resolve => server.listen(0,'127.0.0.1',resolve));
  t.after(() => new Promise(resolve => server.close(resolve)));
  const client = new sdk.OpenAutonomyX({apiKey:'local-contract-test',apiUrl:'http://127.0.0.1:'+server.address().port+'/api/v1'});
  const input = {title:'Contract title',content:'Contract content'};
  assert.equal((await client.content.create(input)).id,'publication-1');
  assert.deepEqual(seen[0],{method:'POST',path:'/api/v1/publications',auth:'Bearer local-contract-test',body:input});
  await assert.rejects(client.content.get('missing'),error => error.response.status === 404);
});
