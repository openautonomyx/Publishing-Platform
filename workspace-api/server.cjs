'use strict';
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');

const PORT = Number(process.env.PORT || 3000);
const DATA_FILE = path.resolve(process.env.PUBLISHING_DATA_FILE || './data/publications.json');
const TOKEN = process.env.PUBLISHING_API_TOKEN;
if (!TOKEN || TOKEN.length < 32) {
  console.error('PUBLISHING_API_TOKEN must be set to a random value of at least 32 characters');
  process.exit(1);
}
const ACTIONS = {submit:['draft'],approve:['in_review'],reject:['in_review'],schedule:['approved'],publish:['approved','scheduled'],archive:['draft','approved','published']};
const timing = () => new Date().toISOString();
function load() {
  try {
    const value = JSON.parse(fs.readFileSync(DATA_FILE, 'utf8'));
    if (!Array.isArray(value)) throw new Error('Data file must contain an array');
    return value;
  } catch (e) {
    if (e.code === 'ENOENT') return [];
    throw e;
  }
}
let documents = load();
function persist() {
  fs.mkdirSync(path.dirname(DATA_FILE), {recursive:true,mode:0o700});
  const temp = DATA_FILE + '.' + process.pid + '.tmp';
  fs.writeFileSync(temp, JSON.stringify(documents, null, 2), {mode:0o600});
  fs.renameSync(temp, DATA_FILE);
}
function respond(res,status,data) {
  res.writeHead(status, {'Content-Type':'application/json; charset=utf-8','Cache-Control':'no-store','X-Content-Type-Options':'nosniff'});
  res.end(JSON.stringify(data));
}
function authenticate(req) {
  const candidate = req.headers.authorization || '';
  const expected = 'Bearer ' + TOKEN;
  const a = Buffer.from(candidate), b = Buffer.from(expected);
  return a.length === b.length && crypto.timingSafeEqual(a,b);
}
function json(req) {
  return new Promise((resolve,reject) => {
    let body = '', size = 0;
    req.on('data', chunk => {
      size += chunk.length;
      if (size > 1024 * 1024) {reject(Object.assign(new Error('Payload too large'),{status:413}));req.pause();return;}
      body += chunk;
    });
    req.on('end', () => {try {resolve(JSON.parse(body || '{}'));}catch {reject(Object.assign(new Error('Invalid JSON'),{status:400}));}});
    req.on('error',reject);
  });
}
function audit(doc,action) {
  doc.history.push({action,at:timing()});
  doc.updatedAt = timing();
}
function validText(value,max=200) {return typeof value === 'string' && value.trim().length>0 && value.length<=max;}
function advanceScheduled() {
  let changed=false;
  for(const doc of documents) {
    if(doc.status === 'scheduled' && doc.scheduledFor && Date.parse(doc.scheduledFor)<=Date.now()) {
      doc.status='published';
      doc.publishedAt=timing();
      audit(doc,'publish_internal');
      changed=true;
    }
  }
  if(changed) persist();
}
const handler = async (req,res) => {
  try {
    const pathname = new URL(req.url,'http://localhost').pathname;
    if(pathname==='/api/health' && req.method==='GET') return respond(res,200,{status:'ok',service:'publishing-api'});
    if(!pathname.startsWith('/api/')) return respond(res,404,{error:'Not found'});
    if(!authenticate(req)) return respond(res,401,{error:'Unauthorized'});
    advanceScheduled();
    if(pathname==='/api/publications' && req.method==='GET') return respond(res,200,{items:documents});
    if(pathname==='/api/publications' && req.method==='POST') {
      const input=await json(req);
      if(!validText(input.title) || !validText(input.content,100000)) return respond(res,422,{error:'Title and content are required (max 200/100000 characters)'});
      const now=timing();
      const doc={id:crypto.randomUUID(),title:input.title.trim(),content:input.content,status:'draft',createdAt:now,updatedAt:now,scheduledFor:null,publishedAt:null,history:[{action:'create',at:now}]};
      documents.unshift(doc);persist();return respond(res,201,doc);
    }
    const match=pathname.match(/^\/api\/publications\/([0-9a-f-]{36})(?:\/(submit|approve|reject|schedule|publish|archive))?$/);
    if(!match) return respond(res,404,{error:'Not found'});
    const doc=documents.find(d=>d.id===match[1]);
    if(!doc) return respond(res,404,{error:'Publication not found'});
    if(req.method==='PUT' && !match[2]) {
      if(doc.status!=='draft') return respond(res,409,{error:'Only drafts may be edited'});
      const input=await json(req);
      if(!validText(input.title) || !validText(input.content,100000)) return respond(res,422,{error:'Title and content are required'});
      doc.title=input.title.trim();doc.content=input.content;audit(doc,'edit');persist();return respond(res,200,doc);
    }
    if(req.method!=='POST'||!match[2]) return respond(res,405,{error:'Method not allowed'});
    const action=match[2];
    if(!ACTIONS[action].includes(doc.status)) return respond(res,409,{error:'Invalid transition from '+doc.status+' via '+action});
    const input=action==='schedule' ? await json(req) : {};
    if(action==='schedule') {
      const when=typeof input.scheduledFor === 'string' ? Date.parse(input.scheduledFor) : NaN;
      if(!Number.isFinite(when)||when<=Date.now()) return respond(res,422,{error:'scheduledFor must be a future ISO-8601 datetime'});
      doc.scheduledFor=new Date(when).toISOString();
    }
    if(action==='publish') doc.publishedAt=timing();
    if(action==='reject') doc.scheduledFor=null;
    const states={submit:'in_review',approve:'approved',reject:'draft',schedule:'scheduled',publish:'published',archive:'archived'};
    doc.status=states[action];audit(doc,action);persist();return respond(res,200,doc);
  } catch(error) {console.error(error);if (!res.headersSent && !res.destroyed) respond(res,error.status||500,{error:error.status?error.message:'Internal server error'});}
};
if(require.main===module) {
  const server=http.createServer(handler);
  server.listen(PORT,()=>console.log('Publishing API listening on '+PORT));
  const timer=setInterval(advanceScheduled,15000);
  timer.unref();
}
module.exports={handler};
