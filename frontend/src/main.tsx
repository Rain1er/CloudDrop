import React, { FormEvent, useDeferredValue, useEffect, useMemo, useRef, useState } from 'react';
import { createRoot } from 'react-dom/client';
import axios from 'axios';
import { Tree } from 'react-arborist';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';
import './styles.css';
import { defaultFilePath, ensureDirectoryPath, parentPath, parseFileList, stripTrailingPathSeparators, type FileNode } from './fileTree';

type ShellType = 'php' | 'java' | 'c#' | 'asp';
type Page = 'home' | 'detail';
type ModuleKey = 'system' | 'terminal' | 'code' | 'files' | 'database';
type WebShell = { ID: number; CreatedAt?: string; UpdatedAt?: string; name: string; url: string; password?: string; type: ShellType; encode: string; status?: string; note?: string };
type ApiEnvelope<T = string> = { info?: T; results?: WebShell[]; data?: WebShell; error?: string; message?: string; status?: string };
type ShellForm = { name: string; url: string; type: ShellType; encode: string; note: string };
type DbForm = { driver: string; host: string; port: string; user: string; pass: string; database: string; sql: string; option: string; encoding: string };
type MenuState = { x: number; y: number; shell: WebShell } | null;

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://127.0.0.1:8989';
const api = axios.create({ baseURL: `${API_BASE}/api/v1`, timeout: 45000 });
const emptyForm: ShellForm = { name: '', url: '', type: 'php', encode: 'utf-8', note: '' };
const defaultDbForm: DbForm = { driver: 'mysql', host: '127.0.0.1', port: '3306', user: 'root', pass: '', database: '', sql: 'SHOW DATABASES', option: '[]', encoding: 'utf8mb4' };

function Icon({ name }: { name: 'cloud' | 'terminal' | 'folder' | 'database' | 'shield' | 'plus' | 'refresh' | 'code' | 'arrow' | 'file' | 'play' | 'edit' | 'trash' }) {
  const common = { width: 20, height: 20, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', strokeWidth: 1.8, strokeLinecap: 'round' as const, strokeLinejoin: 'round' as const };
  const paths: Record<typeof name, React.ReactNode> = {
    cloud: <path d="M17.5 19H7a4 4 0 0 1-.8-7.92 5.6 5.6 0 0 1 10.84-1.7A4.82 4.82 0 0 1 17.5 19Z" />,
    terminal: <><path d="m7 8 4 4-4 4" /><path d="M12 16h5" /><rect x="3" y="4" width="18" height="16" rx="3" /></>,
    folder: <path d="M3 7.5A2.5 2.5 0 0 1 5.5 5H9l2 2h7.5A2.5 2.5 0 0 1 21 9.5v7A2.5 2.5 0 0 1 18.5 19h-13A2.5 2.5 0 0 1 3 16.5Z" />,
    database: <><ellipse cx="12" cy="5" rx="7" ry="3" /><path d="M5 5v12c0 1.66 3.13 3 7 3s7-1.34 7-3V5" /><path d="M5 11c0 1.66 3.13 3 7 3s7-1.34 7-3" /></>,
    shield: <><path d="M12 3 20 6v6c0 5-3.4 7.8-8 9-4.6-1.2-8-4-8-9V6Z" /><path d="m9.5 12 1.8 1.8 3.7-4" /></>,
    plus: <><path d="M12 5v14" /><path d="M5 12h14" /></>,
    refresh: <><path d="M20 12a8 8 0 0 1-14.9 4" /><path d="M4 12A8 8 0 0 1 18.9 8" /><path d="M3 16h4v4" /><path d="M21 8h-4V4" /></>,
    code: <><path d="m9 18-6-6 6-6" /><path d="m15 6 6 6-6 6" /></>,
    arrow: <><path d="M19 12H5" /><path d="m12 19-7-7 7-7" /></>,
    file: <><path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8Z" /><path d="M14 3v5h5" /></>,
    play: <path d="m8 5 11 7-11 7Z" />,
    edit: <><path d="M12 20h9" /><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z" /></>,
    trash: <><path d="M3 6h18" /><path d="M8 6V4h8v2" /><path d="M19 6l-1 14H6L5 6" /></>,
  };
  return <svg {...common} aria-hidden="true">{paths[name]}</svg>;
}

function normalizeShell(shell: WebShell): WebShell { const type = String(shell.type || '').toLowerCase(); return { ...shell, ID: Number(shell.ID), type: (['php', 'java', 'c#', 'asp'].includes(type) ? type : ((type.startsWith('c#') || type === 'net' || type === 'csharp') ? 'c#' : 'c#')) as ShellType }; }
const displayShellType = (type: string) => (String(type || '').toLowerCase().startsWith('c#') ? 'c#' : type);
function codeExampleForShell(type: string) {
  const normalized = String(type || '').toLowerCase();
  if (normalized === 'php') return 'echo "hello world";';
  if (normalized === 'java') return 'System.out.print("hello world");';
  if (normalized === 'asp') return 'GlobalResult = "hello world"';
  return 'public class RunCode {\n    public override string ToString() {\n        return "hello world";\n    }\n}';
}
function languageForShell(type: string) {
  const normalized = String(type || '').toLowerCase();
  if (normalized === 'php') return 'php';
  if (normalized === 'java') return 'java';
  if (normalized === 'asp') return 'vbscript';
  return 'csharp';
}
const statusTone = (status?: string) => status === 'alive' ? 'alive' : status === 'dead' || status === 'unsupported' ? 'dead' : 'unknown';
function readError(err: unknown) { if (axios.isAxiosError(err)) { const data = err.response?.data as ApiEnvelope | undefined; return data?.message || data?.error || err.message; } return err instanceof Error ? err.message : '未知错误'; }
function postForm<T = string>(path: string, data: Record<string, string>) { return api.post<ApiEnvelope<T>>(path, new URLSearchParams(data), { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }); }
function useShells() {
  const [shells, setShells] = useState<WebShell[]>([]); const [loading, setLoading] = useState(false); const [error, setError] = useState('');
  const refresh = async () => { setLoading(true); setError(''); try { const res = await api.get<ApiEnvelope<WebShell[]>>('/webshells/list'); const next = (res.data.results ?? []).map(normalizeShell); setShells(next); return next.length; } catch (err) { setError(readError(err)); return 0; } finally { setLoading(false); } };
  useEffect(() => { refresh(); }, []);
  return { shells, loading, error, refresh };
}

function App() {
  const { shells, loading, error, refresh } = useShells();
  const [page, setPage] = useState<Page>('home');
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [query, setQuery] = useState('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editShell, setEditShell] = useState<WebShell | null>(null);
  const [menu, setMenu] = useState<MenuState>(null);
  const [batchBusy, setBatchBusy] = useState(false);
  const [notice, setNotice] = useState('');
  const deferredQuery = useDeferredValue(query);

  useEffect(() => { if (shells.length === 0) { setSelectedId(null); setPage('home'); return; } if (selectedId && !shells.some((item) => item.ID === selectedId)) { setSelectedId(null); setPage('home'); } }, [shells, selectedId]);
  useEffect(() => { const close = () => setMenu(null); window.addEventListener('click', close); return () => window.removeEventListener('click', close); }, []);
  const filtered = useMemo(() => { const q = deferredQuery.trim().toLowerCase(); if (!q) return shells; return shells.filter((item) => [item.name, item.url, item.type, item.note, item.status].filter(Boolean).join(' ').toLowerCase().includes(q)); }, [shells, deferredQuery]);
  const selected = shells.find((item) => item.ID === selectedId) ?? null;
  const stats = useMemo(() => ({ total: shells.length, alive: shells.filter((item) => item.status === 'alive').length, dead: shells.filter((item) => item.status === 'dead').length }), [shells]);
  const refreshWithNotice = async () => { const count = await refresh(); setNotice(`已刷新 ${count} 条连接`); };
  const batchTest = async () => { setBatchBusy(true); try { const res = await api.post<ApiEnvelope<Record<string, boolean>>>('/webshells/batch-test'); const result = res.data.info ?? {}; const alive = Object.values(result).filter(Boolean).length; setNotice(`存活验证完成：${alive}/${Object.keys(result).length}`); await refresh(); } catch (err) { setNotice(readError(err)); } finally { setBatchBusy(false); } };
  const deleteShell = async (shell: WebShell) => { if (!window.confirm(`确认删除 ${shell.name} ?`)) return; await api.delete(`/webshells/${shell.ID}`); setNotice(`已删除 ${shell.name}`); await refresh(); };
  const openDetail = (shell: WebShell) => { setSelectedId(shell.ID); setPage('detail'); };
  if (page === 'detail' && selected) return <DetailPage shell={selected} onBack={() => setPage('home')} refreshList={refresh} />;

  return <main className="home-shell"><section className="panel home-panel"><div className="panel-head connection-head"><div><span className="eyebrow">Connections</span><h2>连接列表</h2></div><div className="topbar-actions"><label className="search-label"><input aria-label="搜索连接" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="搜索名称、URL、类型、状态" /></label><button className="button ghost" type="button" onClick={refreshWithNotice} disabled={loading}><Icon name="refresh" />{loading ? '刷新中' : '刷新'}</button><button className="button secondary" type="button" onClick={batchTest} disabled={batchBusy}><Icon name="shield" />{batchBusy ? '验证中' : '验证存活'}</button><button className="button primary" type="button" onClick={() => setCreateOpen(true)}><Icon name="plus" />新增连接</button></div></div>{error && <div className="alert danger" role="alert">{error}</div>}{notice && <div className="alert" role="status">{notice}</div>}<div className="table-wrap home-table-wrap"><table><thead><tr><th>名称</th><th>类型</th><th>URL</th><th>状态</th><th>创建时间</th><th>更新时间</th><th>Note</th></tr></thead><tbody>{filtered.map((item) => <tr key={item.ID} onContextMenu={(event) => { event.preventDefault(); setSelectedId(item.ID); setMenu({ x: Math.min(event.clientX, window.innerWidth - 180), y: Math.min(event.clientY, window.innerHeight - 120), shell: item }); }} onDoubleClick={() => openDetail(item)} onClick={() => setSelectedId(item.ID)} className={item.ID === selectedId ? 'selected' : ''} tabIndex={0} onKeyDown={(event) => event.key === 'Enter' && openDetail(item)}><td><strong>{item.name}</strong></td><td><span className="pill">{displayShellType(item.type)}</span></td><td className="mono">{item.url}</td><td><StatusBadge status={item.status} /></td><td className="mono muted">{formatTime(item.CreatedAt)}</td><td className="mono muted">{formatTime(item.UpdatedAt)}</td><td className="note-cell">{formatNote(item.note)}</td></tr>)}{filtered.length === 0 && <tr><td colSpan={7} className="empty">暂无连接，点击右上角新增。</td></tr>}</tbody></table></div></section>{menu && <div className="context-menu" style={{ left: menu.x, top: menu.y }} onClick={(event) => event.stopPropagation()}><button type="button" onClick={() => { setEditShell(menu.shell); setMenu(null); }}><Icon name="edit" />编辑</button><button type="button" className="danger" onClick={() => { const shell = menu.shell; setMenu(null); deleteShell(shell); }}><Icon name="trash" />删除</button></div>}{createOpen && <ShellModal title="新增连接" onClose={() => setCreateOpen(false)} onSaved={() => { setCreateOpen(false); refresh(); setNotice('连接记录已写入数据库'); }} />}{editShell && <ShellModal title="编辑连接" shell={editShell} onClose={() => setEditShell(null)} onSaved={async () => { setEditShell(null); await refresh(); setNotice('连接已更新'); }} />}</main>;
}

function DetailPage({ shell, onBack, refreshList }: { shell: WebShell; onBack: () => void; refreshList: () => Promise<number> }) {
  const [module, setModule] = useState<ModuleKey>('system');
  return <main className="app-shell"><aside className="sidebar" aria-label="CloudDrop detail navigation"><div className="brand"><span className="brand-mark"><Icon name="cloud" /></span><div><strong>CloudDrop</strong><small>Interaction Workspace</small></div></div><button className="nav-item back-button" type="button" onClick={onBack}><Icon name="arrow" />返回连接列表</button><nav className="nav-list" aria-label="功能模块"><NavButton active={module === 'system'} icon="shield" label="系统信息" onClick={() => setModule('system')} /><NavButton active={module === 'terminal'} icon="terminal" label="命令执行" onClick={() => setModule('terminal')} /><NavButton active={module === 'code'} icon="code" label="代码执行" onClick={() => setModule('code')} /><NavButton active={module === 'files'} icon="folder" label="文件管理" onClick={() => setModule('files')} /><NavButton active={module === 'database'} icon="database" label="数据库" onClick={() => setModule('database')} /></nav><div className="sidebar-card"><strong>{shell.name}</strong><p>{shell.url}</p></div></aside><section className="workspace"><ShellWorkspace shell={shell} module={module} refreshList={refreshList} /></section></main>;
}

function ShellWorkspace({ shell, module, refreshList }: { shell: WebShell; module: ModuleKey; refreshList: () => Promise<number> }) {
  const [busy, setBusy] = useState(''); const [systemInfo, setSystemInfo] = useState(''); const [systemError, setSystemError] = useState(''); const [code, setCode] = useState(codeExampleForShell(shell.type)); const [codeOutput, setCodeOutput] = useState('');
  const run = async (label: string, task: () => Promise<string>) => { setBusy(label); try { return await task(); } finally { setBusy(''); } };
  const loadSystemInfo = async () => { setSystemError(''); try { const info = await run('系统信息', async () => (await api.get<ApiEnvelope>(`/webshells/BaseInfo/${shell.ID}`)).data.info ?? ''); setSystemInfo(info); await refreshList(); } catch (err) { setSystemError(readError(err)); } };
  useEffect(() => { setSystemInfo(''); setCodeOutput(''); setCode(codeExampleForShell(shell.type)); loadSystemInfo(); }, [shell.ID, shell.type]);
  const executeCode = async () => { try { const info = await run('代码执行', async () => (await postForm(`/webshells/ExecCode/${shell.ID}`, { code })).data.info ?? ''); setCodeOutput(info || '执行完成，无返回内容。'); } catch (err) { setCodeOutput(readError(err)); } };
  return <div className="workspace-stack">{module === 'system' && <SystemPanel info={systemInfo} error={systemError} busy={busy === '系统信息'} onRefresh={loadSystemInfo} />}{module === 'terminal' && <WebTerminalPanel shell={shell} />}{module === 'code' && <CodePanel shellType={shell.type} code={code} setCode={setCode} output={codeOutput} busy={busy === '代码执行'} onRun={executeCode} />}{module === 'files' && <FileManager shell={shell} />}{module === 'database' && <DatabasePanel shell={shell} />}</div>;
}

function SystemPanel({ info, error, busy, onRefresh }: { info: string; error: string; busy: boolean; onRefresh: () => void }) { const parsed = parseBaseInfo(info); return <section className="panel ide-panel"><div className="panel-head"><div><span className="eyebrow">System</span><h2>系统信息</h2></div><button className="button ghost" type="button" onClick={onRefresh} disabled={busy}><Icon name="refresh" />{busy ? '获取中' : '重新获取'}</button></div>{error && <div className="alert danger">{error}</div>}<div className="system-grid"><div className="system-summary">{['WebPath', 'ClassPath', 'DriveList', 'OS', 'os.name', 'os.arch'].map((key) => <article key={key}><span>{key}</span><strong className="mono">{parsed[key] || '—'}</strong></article>)}</div><pre className="code-view system-output">{info || '正在获取目标系统信息...'}</pre></div></section>; }

function WebTerminalPanel({ shell }: { shell: WebShell }) {
  const hostRef = useRef<HTMLDivElement | null>(null);
  const termRef = useRef<Terminal | null>(null);
  const inputRef = useRef('');
  const busyRef = useRef(false);
  const prompt = '> ';

  useEffect(() => {
    if (!hostRef.current) return;
    const term = new Terminal({
      cursorBlink: true,
      convertEol: true,
      fontFamily: 'Fira Code, ui-monospace, SFMono-Regular, Menlo, monospace',
      fontSize: 14,
      lineHeight: 1.45,
      scrollback: 3000,
      theme: {
        background: '#0f172a',
        foreground: '#dbeafe',
        cursor: '#7dd3fc',
        selectionBackground: '#334155',
        black: '#0f172a',
        red: '#f87171',
        green: '#86efac',
        yellow: '#fde68a',
        blue: '#7dd3fc',
        magenta: '#c4b5fd',
        cyan: '#67e8f9',
        white: '#e2e8f0',
        brightBlack: '#64748b',
        brightRed: '#fca5a5',
        brightGreen: '#bbf7d0',
        brightYellow: '#fef3c7',
        brightBlue: '#bae6fd',
        brightMagenta: '#ddd6fe',
        brightCyan: '#a5f3fc',
        brightWhite: '#f8fafc',
      },
    });
    const fit = new FitAddon();
    term.loadAddon(fit);
    term.open(hostRef.current);
    fit.fit();
    termRef.current = term;
    term.writeln('CloudDrop terminal ready');
    term.write(`\x1b[32m${prompt}\x1b[0m`);

    const writePrompt = () => term.write(`\r\n\x1b[32m${prompt}\x1b[0m`);
    const runCommand = async () => {
      const command = inputRef.current.trim();
      inputRef.current = '';
      if (!command || busyRef.current) {
        writePrompt();
        return;
      }
      busyRef.current = true;
      term.write('\r\n');
      try {
        const res = await postForm(`/webshells/ExecCommand/${shell.ID}`, { command });
        const output = String(res.data.info ?? '').replace(/\n$/, '');
        term.writeln(output || '执行完成，无返回内容。');
      } catch (err) {
        term.writeln(`\x1b[31m${readError(err)}\x1b[0m`);
      } finally {
        busyRef.current = false;
        writePrompt();
      }
    };

    const disposable = term.onData((data) => {
      if (data === '\r') { runCommand(); return; }
      if (data === '\u007F') {
        if (inputRef.current.length > 0) {
          inputRef.current = inputRef.current.slice(0, -1);
          term.write('\b \b');
        }
        return;
      }
      if (data === '\u0003') {
        inputRef.current = '';
        term.write('^C');
        writePrompt();
        return;
      }
      if (data >= ' ') {
        inputRef.current += data;
        term.write(data);
      }
    });

    const resize = () => fit.fit();
    window.addEventListener('resize', resize);
    return () => {
      disposable.dispose();
      window.removeEventListener('resize', resize);
      term.dispose();
      termRef.current = null;
    };
  }, [shell.ID]);

  return <section className="panel xterm-shell-panel"><div className="xterm-console"><div className="terminal-toolbar"><div className="traffic"><span /><span /><span /></div><strong>Terminal</strong><span>React + Axios + xterm.js</span></div><div className="xterm-host monokai-xterm-host" ref={hostRef} /></div></section>;
}

function CodePanel({ shellType, code, setCode, output, busy, onRun }: { shellType: ShellType; code: string; setCode: (value: string) => void; output: string; busy: boolean; onRun: () => void }) { const example = codeExampleForShell(shellType); const language = languageForShell(shellType); return <section className="panel ide-panel"><div className="panel-head"><div><span className="eyebrow">Code Runner</span><h2>代码执行</h2></div><button className="button primary" type="button" onClick={onRun} disabled={busy}><Icon name="code" />{busy ? '执行中' : '执行代码'}</button></div><div className="code-example"><div><span className="field-title">{displayShellType(shellType)} hello world 示例</span><SyntaxHighlighter language={language} style={oneLight} customStyle={{ ...syntaxStyle, minHeight: 80 }}>{example}</SyntaxHighlighter></div><button className="button ghost" type="button" onClick={() => setCode(example)}>使用示例</button></div><div className="code-runner"><label><span>Payload</span><textarea className="code-editor" value={code} onChange={(event) => setCode(event.target.value)} rows={14} /></label><div><span className="field-title">返回结果</span><SyntaxHighlighter language={language} style={oneLight} customStyle={syntaxStyle}>{output || '执行结果会显示在这里。'}</SyntaxHighlighter></div></div></section>; }

function FileManager({ shell }: { shell: WebShell }) {
  const [path, setPath] = useState(defaultFilePath(shell.type)); const [nodes, setNodes] = useState<FileNode[]>([]); const [selected, setSelected] = useState<FileNode | null>(null); const [content, setContent] = useState(''); const [busy, setBusy] = useState(''); const [error, setError] = useState(''); const [treeHeight, setTreeHeight] = useState(Math.max(520, window.innerHeight - 180));
  const loadPath = async (nextPath = path) => { const requestedPath = ensureDirectoryPath(nextPath || defaultFilePath(shell.type)); setBusy('list'); setError(''); try { const res = await postForm(`/webshells/FileList/${shell.ID}`, { path: requestedPath }); const parsed = parseFileList(String(res.data.info ?? ''), requestedPath); setNodes(parsed.nodes); setPath(parsed.currentPath); setSelected(null); setContent(''); } catch (err) { setError(readError(err)); } finally { setBusy(''); } };
  const openNode = async (node: FileNode) => { setSelected(node); if (node.type === 'dir') { await loadPath(node.path); return; } setBusy('show'); setError(''); try { const res = await postForm(`/webshells/FileShow/${shell.ID}`, { path: node.path }); setContent(String(res.data.info ?? '')); } catch (err) { setContent(readError(err)); } finally { setBusy(''); } };
  const runArchive = async (mode: 'zip' | 'unzip') => { const defaultSrc = selected?.path || path; const srcPath = window.prompt(mode === 'zip' ? '输入要压缩的源路径' : '输入要解压的压缩包路径', defaultSrc); if (!srcPath) return; const defaultTo = mode === 'zip' ? `${stripTrailingPathSeparators(srcPath)}.zip` : path; const toPath = window.prompt(mode === 'zip' ? '输入压缩包输出路径' : '输入解压目标目录', defaultTo); if (!toPath) return; setBusy(mode); setError(''); try { const endpoint = mode === 'zip' ? 'FileZip' : 'FileUnZip'; const res = await postForm(`/webshells/${endpoint}/${shell.ID}`, { srcPath, toPath }); setContent(String(res.data.info ?? '执行完成')); await loadPath(path); } catch (err) { setError(readError(err)); } finally { setBusy(''); } };
  useEffect(() => { const resize = () => setTreeHeight(Math.max(520, window.innerHeight - 180)); resize(); window.addEventListener('resize', resize); return () => window.removeEventListener('resize', resize); }, []);
  useEffect(() => { loadPath(defaultFilePath(shell.type)); }, [shell.ID, shell.type]);
  return <section className="panel file-ide"><div className="file-toolbar"><button className="button ghost" type="button" onClick={() => loadPath(parentPath(path))}>上级</button><button className="button secondary" type="button" onClick={() => runArchive('zip')} disabled={Boolean(busy)}>压缩</button><button className="button secondary" type="button" onClick={() => runArchive('unzip')} disabled={Boolean(busy)}>解压</button><input aria-label="文件路径" value={path} onChange={(event) => setPath(event.target.value)} onKeyDown={(event) => event.key === 'Enter' && loadPath(path)} /><button className="button primary" type="button" onClick={() => loadPath(path)} disabled={busy === 'list'}><Icon name="refresh" />查询</button></div>{error && <div className="alert danger">{error}</div>}<div className="file-layout"><aside className="file-tree"><div className="tree-head"><strong>react-arborist</strong><span>{nodes.length} items</span></div><Tree data={nodes} openByDefault width="100%" height={treeHeight} rowHeight={34} indent={18} onActivate={(node: any) => openNode(node.data)}>{({ node, style }: any) => <div style={style} className={`tree-node arborist-row ${selected?.path === node.data.path ? 'active' : ''}`}><Icon name={node.data.type === 'dir' ? 'folder' : 'file'} /><span>{node.data.name}</span><small>{node.data.type === 'dir' ? 'dir' : node.data.size}</small></div>}</Tree>{nodes.length === 0 && <p className="muted tree-empty">{busy ? '加载中...' : '目录为空'}</p>}</aside><main className="file-preview"><div className="preview-head"><strong>{selected?.name || '选择文件查看'}</strong><span className="mono">{selected?.path || path}</span></div><SyntaxHighlighter language={guessLanguage(selected?.name)} style={oneLight} showLineNumbers customStyle={syntaxStyle}>{content || '点击左侧文件查看内容。'}</SyntaxHighlighter></main></div></section>;
}

function DatabasePanel({ shell }: { shell: WebShell }) {
  const [form, setForm] = useState<DbForm>(defaultDbForm); const [output, setOutput] = useState(''); const [busy, setBusy] = useState(false);
  const runSql = async (override?: Partial<DbForm>) => { setBusy(true); const payload = { ...form, ...override }; setForm(payload); try { const res = await postForm(`/webshells/ExecSql/${shell.ID}`, payload); setOutput(String(res.data.info ?? '')); } catch (err) { setOutput(JSON.stringify({ status: 'error', message: readError(err) }, null, 2)); } finally { setBusy(false); } };
  const parsed = parseMaybeJson(output);
  return <section className="panel db-panel"><div className="panel-head"><div><span className="eyebrow">Database</span><h2>数据库管理</h2></div><div className="topbar-actions"><button className="button ghost" type="button" onClick={() => runSql({ database: '', sql: 'SHOW DATABASES' })} disabled={busy}>列库</button><button className="button secondary" type="button" onClick={() => runSql({ sql: 'SHOW TABLES' })} disabled={busy}>列表</button><button className="button primary" type="button" onClick={() => runSql()} disabled={busy}><Icon name="database" />执行</button></div></div><div className="db-layout"><form className="db-config" onSubmit={(event) => { event.preventDefault(); runSql(); }}><label><span>驱动</span><select value={form.driver} onChange={(event) => setForm({ ...form, driver: event.target.value })}><option value="mysql">mysql</option><option value="pgsql">pgsql</option><option value="sqlite">sqlite</option><option value="sqlsrv">sqlsrv</option></select></label><label><span>Host</span><input value={form.host} onChange={(event) => setForm({ ...form, host: event.target.value })} /></label><label><span>Port</span><input value={form.port} onChange={(event) => setForm({ ...form, port: event.target.value })} /></label><label><span>User</span><input value={form.user} onChange={(event) => setForm({ ...form, user: event.target.value })} /></label><label><span>Password</span><input value={form.pass} onChange={(event) => setForm({ ...form, pass: event.target.value })} /></label><label><span>Database</span><input value={form.database} onChange={(event) => setForm({ ...form, database: event.target.value })} /></label><label className="span-2"><span>SQL</span><textarea className="code-editor" value={form.sql} onChange={(event) => setForm({ ...form, sql: event.target.value })} rows={8} /></label><button className="button primary span-2" type="submit" disabled={busy}>执行 SQL</button></form><div className="db-result"><DbResult value={parsed} raw={output} /></div></div></section>;
}

function DbResult({ value, raw }: { value: unknown; raw: string }) { if (!raw) return <pre className="code-view">数据库返回会显示在这里。</pre>; if (typeof value === 'object' && value && 'databases' in value && Array.isArray((value as { databases: unknown[] }).databases)) return <div className="db-cards">{(value as { databases: string[] }).databases.map((name) => <article key={name}><Icon name="database" /><strong>{name}</strong></article>)}</div>; const rows = typeof value === 'object' && value && 'result' in value ? (value as { result?: { data?: Record<string, unknown>[] } }).result?.data : undefined; if (Array.isArray(rows) && rows.length > 0) return <div className="table-wrap"><table><thead><tr>{Object.keys(rows[0]).map((key) => <th key={key}>{key}</th>)}</tr></thead><tbody>{rows.map((row, index) => <tr key={index}>{Object.values(row).map((cell, cellIndex) => <td key={cellIndex} className="mono">{String(cell)}</td>)}</tr>)}</tbody></table></div>; return <SyntaxHighlighter language="json" style={oneLight} customStyle={syntaxStyle}>{typeof value === 'string' ? value : JSON.stringify(value, null, 2)}</SyntaxHighlighter>; }
function Metric({ label, value, tone = 'blue' }: { label: string; value: number; tone?: string }) { return <article className={`metric ${tone}`}><span>{label}</span><strong>{value}</strong></article>; }
function StatusBadge({ status }: { status?: string }) { const tone = statusTone(status); const text = status === 'alive' ? '存活' : status === 'dead' ? '异常' : status === 'unsupported' ? '不支持' : '未验证'; return <span className={`status-badge ${tone}`}><span className="status-dot" />{text}</span>; }
function NavButton({ active, icon, label, onClick }: { active: boolean; icon: 'shield' | 'terminal' | 'code' | 'folder' | 'database'; label: string; onClick: () => void }) { return <button className={`nav-item ${active ? 'active' : ''}`} type="button" onClick={onClick}><Icon name={icon} />{label}</button>; }
function formatTime(value?: string) { if (!value || value.startsWith('0001-')) return '—'; return new Date(value).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }); }
function formatNote(value?: string) { return value && value !== 'null' ? value : '—'; }
function parseBaseInfo(info: string) { return info.split(/\r?\n/).reduce<Record<string, string>>((acc, line) => { const idx = line.indexOf(':'); if (idx > 0) acc[line.slice(0, idx)] = line.slice(idx + 1); return acc; }, {}); }
function parseMaybeJson(raw: string): unknown { try { return JSON.parse(raw); } catch { return raw; } }
function guessLanguage(name?: string) { if (!name) return 'text'; const lower = name.toLowerCase(); if (lower.endsWith('.php')) return 'php'; if (lower.endsWith('.jsp') || lower.endsWith('.java')) return 'java'; if (lower.endsWith('.asp') || lower.endsWith('.aspx') || lower.endsWith('.ashx') || lower.endsWith('.asmx') || lower.endsWith('.vbs')) return 'vbscript'; if (lower.endsWith('.cs')) return 'csharp'; if (lower.endsWith('.js')) return 'javascript'; if (lower.endsWith('.json')) return 'json'; if (lower.endsWith('.html')) return 'html'; if (lower.endsWith('.css')) return 'css'; return 'text'; }
const syntaxStyle: React.CSSProperties = { minHeight: 240, margin: 0, borderRadius: 18, border: '1px solid #dbeafe', background: '#f8fafc', fontSize: 13 };

function ShellModal({ title, shell, onClose, onSaved }: { title: string; shell?: WebShell; onClose: () => void; onSaved: () => void }) { const [form, setForm] = useState<ShellForm>(shell ? { name: shell.name, url: shell.url, type: shell.type, encode: shell.encode, note: formatNote(shell.note) === '—' ? '' : shell.note ?? '' } : emptyForm); const [error, setError] = useState(''); const [busy, setBusy] = useState(false); const submit = async (event: FormEvent) => { event.preventDefault(); setBusy(true); setError(''); try { if (shell) await api.put(`/webshells/${shell.ID}`, { ...form, password: '' }); else await api.post('/webshells/create', { ...form, password: '' }); onSaved(); } catch (err) { setError(readError(err)); } finally { setBusy(false); } }; return <div className="modal-backdrop" role="presentation" onMouseDown={onClose}><section className="modal" role="dialog" aria-modal="true" aria-labelledby="shell-title" onMouseDown={(event) => event.stopPropagation()}><div className="panel-head compact"><div><span className="eyebrow">Target</span><h2 id="shell-title">{title}</h2></div><button className="icon-button" type="button" onClick={onClose} aria-label="关闭">×</button></div>{error && <div className="alert danger">{error}</div>}<form className="create-form" onSubmit={submit}><label><span>名称</span><input required value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} /></label><label><span>URL</span><input required value={form.url} onChange={(event) => setForm({ ...form, url: event.target.value })} placeholder="http://127.0.0.1/api.php" /></label><p className="form-hint">当前协议会在后端自动生成通信种子，新增连接无需填写首次连接密码。</p><label><span>类型</span><select value={form.type} onChange={(event) => setForm({ ...form, type: event.target.value as ShellType })}><option value="php">php</option><option value="java">java</option><option value="c#">c#</option><option value="asp">asp</option></select></label><label><span>编码</span><input required value={form.encode} onChange={(event) => setForm({ ...form, encode: event.target.value })} /></label><label><span>备注</span><textarea value={form.note} onChange={(event) => setForm({ ...form, note: event.target.value })} rows={3} /></label><div className="modal-actions"><button className="button ghost" type="button" onClick={onClose}>取消</button><button className="button primary" type="submit" disabled={busy}>{busy ? '保存中' : '保存连接'}</button></div></form></section></div>; }
class ErrorBoundary extends React.Component<{ children: React.ReactNode }, { hasError: boolean }> { state = { hasError: false }; static getDerivedStateFromError() { return { hasError: true }; } render() { return this.state.hasError ? <main className="home-shell"><div className="alert danger">前端渲染异常，请刷新页面。</div></main> : this.props.children; } }
createRoot(document.getElementById('root')!).render(<ErrorBoundary><App /></ErrorBoundary>);
