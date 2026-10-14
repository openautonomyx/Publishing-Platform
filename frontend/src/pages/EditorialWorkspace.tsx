import React, {useCallback, useEffect, useState} from 'react';
import './Pages.css';

type Status = 'draft'|'in_review'|'approved'|'scheduled'|'published'|'archived';
interface Publication {id:string;title:string;content:string;status:Status;updatedAt:string;scheduledFor:string|null;publishedAt:string|null}
interface Props {onContextChange:(context:string,data:any)=>void}
const API='/api/publications';
async function request(endpoint:string, token:string, method='GET', body?:object) {
  const response=await fetch(endpoint,{method,headers:{Authorization:'Bearer '+token,'Content-Type':'application/json'},body:body?JSON.stringify(body):undefined});
  const data=await response.json();
  if(!response.ok) throw new Error(data.error||'Request failed');
  return data;
}
export default function ContentPage({onContextChange}:Props) {
  const [token,setToken]=useState('');
  const [items,setItems]=useState<Publication[]>([]);
  const [selected,setSelected]=useState<string|null>(null);
  const [title,setTitle]=useState('');
  const [content,setContent]=useState('');
  const [schedule,setSchedule]=useState('');
  const [error,setError]=useState('');
  const [busy,setBusy]=useState(false);
  const refresh=useCallback(async()=>{
    if(!token) return;
    try {const data=await request(API,token);setItems(data.items);setError('');}
    catch(e) {setError((e as Error).message);}
  },[token]);
  useEffect(()=>{onContextChange('content',{});},[]);
  useEffect(()=>{refresh();},[refresh]);
  const current=items.find(x=>x.id===selected);
  const choose=(item?:Publication)=>{
    setSelected(item?.id||null);setTitle(item?.title||'');setContent(item?.content||'');setError('');
  };
  async function execute(action:'save'|'submit'|'approve'|'reject'|'schedule'|'publish'|'archive') {
    setBusy(true);setError('');
    try {
      if(action==='save') {
        const doc=await request(current?API+'/'+current.id:API,token,current?'PUT':'POST',{title,content});
        setSelected(doc.id);
      } else {
        if(!current) throw new Error('Save the publication first');
        await request(API+'/'+current.id+'/'+action,token,'POST',action==='schedule'?{scheduledFor:new Date(schedule).toISOString()}:undefined);
      }
      await refresh();
    } catch(e) {setError((e as Error).message);} finally {setBusy(false);}
  }
  const allowed = (action:string)=> {
    const status=current?.status;
    return (action==='submit'&&status==='draft')||(action==='approve'&&status==='in_review')||(action==='reject'&&status==='in_review')||(action==='schedule'&&status==='approved')||(action==='publish'&&(status==='approved'||status==='scheduled'))||(action==='archive'&&(status==='draft'||status==='approved'||status==='published'));
  };
  return <div className="page-container">
    <div className="page-header"><h1>Editorial workspace</h1><p>Draft, review, approve, schedule and publish to the internal catalog.</p></div>
    <div className="blog-post" style={{marginBottom:20}}>
      <label htmlFor="workspace-key">Workspace API token (kept only in this tab)</label>
      <input id="workspace-key" type="password" value={token} onChange={e=>setToken(e.target.value)} placeholder="Enter PUBLISHING_API_TOKEN" style={{display:'block',width:'100%',marginTop:8,padding:12}} autoComplete="off"/>
      <button className="btn-secondary" style={{marginTop:10}} onClick={refresh}>Connect / Refresh</button>
    </div>
    {error&&<p role="alert" style={{color:'#fca5a5'}}>{error}</p>}
    <div className="content-section">
      <div className="blog-post">
        <h3>Publications ({items.length})</h3>
        <button className="btn-secondary" onClick={()=>choose()}>+ New draft</button>
        <div style={{display:'grid',gap:8,marginTop:16}}>
          {items.map(item=><button key={item.id} className="btn-secondary" onClick={()=>choose(item)} style={{textAlign:'left',border:item.id===selected?'2px solid #60a5fa':undefined}}>{item.title} — {item.status}</button>)}
        </div>
      </div>
      <div className="blog-post">
        <h3>{current?'Edit publication':'Create publication'}</h3>
        <input aria-label="Title" placeholder="Title" value={title} onChange={e=>setTitle(e.target.value)} disabled={!!current&&current.status!=='draft'} style={{width:'100%',padding:12,marginBottom:12}}/>
        <textarea aria-label="Content" placeholder="Write your article..." value={content} onChange={e=>setContent(e.target.value)} disabled={!!current&&current.status!=='draft'} rows={12} style={{width:'100%',padding:12,marginBottom:12}}/>
        {(!current||current.status==='draft')&&<button className="btn-action" disabled={busy||!token} onClick={()=>execute('save')}>Save draft</button>}
        {current&&<p>Status: <strong>{current.status}</strong>{current.scheduledFor?' · Due '+new Date(current.scheduledFor).toLocaleString():''}</p>}
        {(['submit','approve','reject','publish','archive'] as const).filter(allowed).map(action=><button key={action} className="btn-secondary" disabled={busy||!token} style={{margin:4}} onClick={()=>execute(action)}>{action}</button>)}
        {allowed('schedule')&&<div style={{marginTop:12}}><label htmlFor="publish-at">Schedule publication</label><input id="publish-at" type="datetime-local" value={schedule} onChange={e=>setSchedule(e.target.value)} style={{padding:10,margin:8}}/><button className="btn-secondary" disabled={busy||!schedule} onClick={()=>execute('schedule')}>Schedule</button></div>}
        <p style={{fontSize:12,marginTop:16}}>Publication marks content as published internally. External social channels require separate authenticated adapters and are not simulated.</p>
      </div>
    </div>
  </div>;
}
