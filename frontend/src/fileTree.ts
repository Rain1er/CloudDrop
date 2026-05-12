export type FileNode = {
  id: string;
  name: string;
  path: string;
  type: 'dir' | 'file';
  size: string;
  perms: string;
  modified: string;
  children?: FileNode[];
};

export type FileListParseResult = {
  roots: string;
  currentPath: string;
  nodes: FileNode[];
};

export function defaultFilePath(shellType: string) {
  return String(shellType || '').toLowerCase() === 'php' ? '/app/' : '/';
}

export function isWindowsPath(path: string) {
  return /^[a-zA-Z]:[\\/]/.test(path.trim()) || path.includes('\\');
}

function pathSeparator(path: string) {
  return isWindowsPath(path) ? '\\' : '/';
}

export function stripTrailingPathSeparators(path: string) {
  const value = String(path || '').trim();
  if (!value) return '/';
  if (value === '/') return '/';
  const drive = value.match(/^([a-zA-Z]:)[\\/]*$/);
  if (drive) return `${drive[1]}\\`;
  return value.replace(/[\\/]+$/g, '');
}

export function ensureDirectoryPath(path: string) {
  const stripped = stripTrailingPathSeparators(path);
  if (stripped === '/' || /^[a-zA-Z]:\\$/.test(stripped)) return stripped;
  return `${stripped}${pathSeparator(stripped)}`;
}

export function joinPath(base: string, name: string) {
  const root = stripTrailingPathSeparators(base || '/');
  const child = String(name || '').replace(/^[\\/]+/g, '');
  if (root === '/') return `/${child}`;
  if (/^[a-zA-Z]:\\$/.test(root)) return `${root}${child}`;
  return `${root}${pathSeparator(root)}${child}`;
}

export function parentPath(path: string) {
  const stripped = stripTrailingPathSeparators(path);
  if (stripped === '/' || /^[a-zA-Z]:\\$/.test(stripped)) return stripped;
  const windows = isWindowsPath(stripped);
  const normalized = windows ? stripped.replace(/\//g, '\\') : stripped;
  const sep = windows ? '\\' : '/';
  const idx = normalized.lastIndexOf(sep);
  if (idx < 0) return windows ? stripped : '/';
  if (windows && idx === 2) return `${normalized.slice(0, 2)}\\`;
  if (!windows && idx === 0) return '/';
  return normalized.slice(0, idx);
}

function looksLikePath(value: string) {
  const text = value.trim();
  return text === '/' || text.startsWith('/') || /^[a-zA-Z]:[\\/]/.test(text);
}

function inferDirectoryFromBareName(name: string) {
  const clean = name.trim();
  if (!clean) return false;
  if (clean === '.' || clean === '..') return true;
  return !/\.[^\\/.\s]+$/.test(clean);
}

export function parseFileList(raw: string, requestedPath: string): FileListParseResult {
  const lines = String(raw || '')
    .split(/\r?\n/)
    .map((line) => line.trimEnd())
    .filter((line) => line.trim() !== '');
  const hasHeader = lines.length >= 2 && looksLikePath(lines[1]);
  const roots = hasHeader ? lines[0].trim() : '';
  const currentPath = ensureDirectoryPath(hasHeader ? lines[1].trim() : requestedPath);
  const entryLines = hasHeader ? lines.slice(2) : lines;

  const nodes = entryLines
    .map((line): FileNode | null => {
      const columns = line.split('\t');
      const rawName = (columns[0] || '').trim();
      if (!rawName || rawName === '.' || rawName === '..') return null;
      const explicitDir = /^dic:/i.test(rawName);
      const name = rawName.replace(/^dic:/i, '');
      const structured = columns.length > 1;
      const type: FileNode['type'] = explicitDir || (!structured && inferDirectoryFromBareName(name)) ? 'dir' : 'file';
      const size = type === 'dir' && columns[1] === '-' ? '' : (columns[1] || '');
      const nodePath = joinPath(currentPath, name);
      return {
        id: nodePath,
        name,
        path: nodePath,
        type,
        size,
        perms: columns[2] || '',
        modified: columns[3] || '',
        children: type === 'dir' ? ([] as FileNode[]) : undefined,
      };
    })
    .filter((node): node is FileNode => Boolean(node))
    .sort((a, b) => (a.type === b.type ? a.name.localeCompare(b.name) : a.type === 'dir' ? -1 : 1));

  return { roots, currentPath, nodes };
}
