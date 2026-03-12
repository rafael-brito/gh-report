import React, { useState } from 'react';
import { useFileHistoryReport, getFileHistoryDownloadUrl } from '../api/reports';

export const FileHistoryPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [file, setFile] = useState('');
  const [limit, setLimit] = useState(10);
  const [mode, setMode] = useState<'commits' | 'prs'>('prs');
  const [submitted, setSubmitted] = useState(false);

  const queryParams = submitted && repo && file
    ? { repo, file, limit, mode }
    : { repo: '', file: '' };

  const { data, isLoading, error } = useFileHistoryReport(queryParams as any);

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
  };

  const canDownload = submitted && !!repo && !!file;

  const handleOpenMarkdown = () => {
    if (!canDownload) return;
    const url = getFileHistoryDownloadUrl({ repo, file, limit, mode }, 'markdown');
    window.open(url, '_blank');
  };

  const handleOpenCSV = () => {
    if (!canDownload) return;
    const url = getFileHistoryDownloadUrl({ repo, file, limit, mode }, 'csv');
    window.open(url, '_blank');
  };

  return (
    <div style={{ padding: '1rem' }}>
      <h1>Relatório de Histórico de Arquivo</h1>

      <form onSubmit={onSubmit} style={{ marginBottom: '1rem' }}>
        <div>
          <label>Repo (owner/repo): </label>
          <input
            value={repo}
            onChange={e => setRepo(e.target.value)}
            placeholder="org/projeto-x"
            style={{ width: '300px' }}
          />
        </div>
        <div>
          <label>File path: </label>
          <input
            value={file}
            onChange={e => setFile(e.target.value)}
            placeholder="src/main/.../MinhaClasse.java"
            style={{ width: '300px' }}
          />
        </div>
        <div>
          <label>Limit: </label>
          <input
            type="number"
            value={limit}
            onChange={e => setLimit(Number(e.target.value) || 10)}
            style={{ width: '80px' }}
          />
        </div>
        <div>
          <label>Mode: </label>
          <select
            value={mode}
            onChange={e => setMode(e.target.value as 'commits' | 'prs')}
          >
            <option value="prs">PRs</option>
            <option value="commits">Commits</option>
          </select>
        </div>
        <button type="submit" style={{ marginTop: '0.5rem' }}>
          Gerar
        </button>
      </form>

      {canDownload && (
        <div style={{ marginBottom: '1rem', display: 'flex', gap: '0.5rem' }}>
          <button type="button" onClick={handleOpenMarkdown}>
            Abrir em Markdown
          </button>
          <button type="button" onClick={handleOpenCSV}>
            Baixar CSV
          </button>
        </div>
      )}

      {isLoading && <p>Carregando...</p>}
      {error && <p style={{ color: 'red' }}>{(error as Error).message}</p>}

      {data && (
        <pre style={{ background: '#f5f5f5', padding: '1rem', maxHeight: '400px', overflow: 'auto' }}>
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  );
};
