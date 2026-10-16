'use strict';
const {test, before, after} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const dir=fs.mkdtempSync(path.join(os.tmpdir(),'publishing-test-'));
process.env.PUBLISHING_API_TOKEN='test-token-for-local-api-testing-only-12345';
process.env.PUBLISHING_DATA_FILE=path.join(dir,'data.json');
const {handler}=require('./server.cjs');
const server=require('node:http').createServer(handler);
let base;
before(async()=>{await new Promise(resolve=>server.listen(0,'127.0.0.1',resolve));base='http://127.0.0.1:'+server.address().port;});
after(async()=>{await new Promise(resolve=>server.close(resolve));fs.rmSync(dir,{recursive:true,force:true});});
async function call(route,method='GET',body,authorized=true){
 const res=await fetch(base+route,{method,headers:{...(authorized?{Authorization:'Bearer '+process.env.PUBLISHING_API_TOKEN}:{}),'Content-Type':'application/json'},body:body?JSON.stringify(body):undefined});
 return {status:res.status,data:await res.json()};
}
test('health, authentication, validation and editorial lifecycle',async()=>{
 assert.equal((await call('/api/health', 'GET', undefined,false)).status,200);
 assert.equal((await call('/api/publications','GET',undefined,false)).status,401);
 assert.equal((await call('/api/publications','POST',{title:'',content:''})).status,422);
 const created=await call('/api/publications','POST',{title:'An article',content:'An original draft'});
 assert.equal(created.status,201);
 const id=created.data.id;
 assert.equal((await call('/api/publications/'+id+'/approve','POST')).status,409);
 assert.equal((await call('/api/publications/'+id,'PUT',{title:'Updated',content:'Updated draft'})).status,200);
 assert.equal((await call('/api/publications/'+id+'/submit','POST')).data.status,'in_review');
 assert.equal((await call('/api/publications/'+id+'/reject','POST')).data.status,'draft');
 assert.equal((await call('/api/publications/'+id+'/submit','POST')).data.status,'in_review');
 assert.equal((await call('/api/publications/'+id+'/approve','POST')).data.status,'approved');
 assert.equal((await call('/api/publications/'+id+'/schedule','POST',{scheduledFor:'2020-01-01T00:00:00Z'})).status,422);
 assert.equal((await call('/api/publications/'+id+'/schedule','POST',{scheduledFor:42})).status,422);
 assert.equal((await call('/api/publications/'+id+'/schedule','POST',{scheduledFor:null})).status,422);
 const scheduled=await call('/api/publications/'+id+'/schedule','POST',{scheduledFor:new Date(Date.now()+60000).toISOString()});
 assert.equal(scheduled.data.status,'scheduled');
 assert.equal((await call('/api/publications/'+id+'/publish','POST')).data.status,'published');
 assert.equal((await call('/api/publications/'+id+'/archive','POST')).data.status,'archived');
 assert.equal((await call('/api/publications')).data.items.length,1);
 assert.equal(JSON.parse(fs.readFileSync(process.env.PUBLISHING_DATA_FILE,'utf8')).length,1);
});

test('rejects malformed JSON without mutating publications', async()=>{
 const res=await fetch(base+'/api/publications',{method:'POST',headers:{Authorization:'Bearer '+process.env.PUBLISHING_API_TOKEN,'Content-Type':'application/json'},body:'{bad json'});
 assert.equal(res.status,400);
 assert.equal((await call('/api/publications')).data.items.length,1);
});

test('rejects non-object JSON bodies with validation errors', async()=>{
 for (const body of ['null', '[]', '42', '"text"']) {
  const response=await fetch(base+'/api/publications',{method:'POST',headers:{Authorization:'Bearer '+process.env.PUBLISHING_API_TOKEN,'Content-Type':'application/json'},body});
  assert.equal(response.status,422,body);
 }
});

test('failed persistence never exposes an uncommitted create or edit', async()=>{
 const file=process.env.PUBLISHING_DATA_FILE;
 const original=fs.readFileSync(file);
 const before=(await call('/api/publications')).data.items;
 fs.renameSync(file,file+'.backup');
 fs.mkdirSync(file); // Force the atomic replacement to fail deterministically.
 try {
  assert.equal((await call('/api/publications','POST',{title:'Must roll back',content:'Never committed'})).status,500);
  assert.deepEqual((await call('/api/publications')).data.items,before);
 } finally {
  fs.rmdirSync(file);
  fs.renameSync(file+'.backup',file);
 }
 assert.deepEqual(fs.readFileSync(file),original);
 const created=await call('/api/publications','POST',{title:'Original',content:'Persisted body'});
 fs.renameSync(file,file+'.backup');fs.mkdirSync(file);
 try {
  assert.equal((await call('/api/publications/'+created.data.id,'PUT',{title:'Uncommitted',content:'Uncommitted body'})).status,500);
  const stored=(await call('/api/publications')).data.items.find(d=>d.id===created.data.id);
  assert.deepEqual(stored,created.data);
 } finally {fs.rmdirSync(file);fs.renameSync(file+'.backup',file);}
});

test('rechecks draft state after a slow edit body finishes arriving', async()=>{
 const http=require('node:http');
 const created=await call('/api/publications','POST',{title:'Race original',content:'Original body'});
 const route='/api/publications/'+created.data.id;
 const body=JSON.stringify({title:'Late edit',content:'Must not overwrite review content'});
 let request;
 const response=new Promise((resolve,reject)=>{
  request=http.request(base+route,{method:'PUT',headers:{Authorization:'Bearer '+process.env.PUBLISHING_API_TOKEN,'Content-Type':'application/json','Content-Length':Buffer.byteLength(body)}},res=>{
   res.resume();res.on('end',()=>resolve(res.statusCode));
  });
  request.on('error',reject);
  request.write(body.slice(0,10));
 });
 // Give the server a turn to observe the partial body before the transition.
 await new Promise(resolve=>setTimeout(resolve,30));
 assert.equal((await call(route+'/submit','POST')).status,200);
 request.end(body.slice(10));
 assert.equal(await response,409);
 const stored=(await call('/api/publications')).data.items.find(d=>d.id===created.data.id);
 assert.equal(stored.status,'in_review');assert.equal(stored.title,'Race original');
});

test('oversized body is rejected and subsequent requests remain usable', async()=>{
 const res=await fetch(base+'/api/publications',{method:'POST',headers:{Authorization:'Bearer '+process.env.PUBLISHING_API_TOKEN,'Content-Type':'application/json'},body:JSON.stringify({title:'Too big',content:'x'.repeat(1024*1024)})});
 assert.equal(res.status,413);
 assert.equal((await call('/api/health')).status,200);
});
