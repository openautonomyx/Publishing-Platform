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
 const scheduled=await call('/api/publications/'+id+'/schedule','POST',{scheduledFor:new Date(Date.now()+60000).toISOString()});
 assert.equal(scheduled.data.status,'scheduled');
 assert.equal((await call('/api/publications/'+id+'/publish','POST')).data.status,'published');
 assert.equal((await call('/api/publications/'+id+'/archive','POST')).data.status,'archived');
 assert.equal((await call('/api/publications')).data.items.length,1);
 assert.equal(JSON.parse(fs.readFileSync(process.env.PUBLISHING_DATA_FILE,'utf8')).length,1);
});
