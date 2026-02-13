import React, { useState } from 'react';
import { useFileHistoryReport } from '../api/reports';

export const FileHistoryPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [file, setFile] = useState('');
  const [limit, setLimit] = useState(10);
  const [mode, setMode] = useState<'commits' | 'prs'>('prs');
  const [submitted, setSubmitted] = useState(false);

  const { data, isLoading, error } = useFileHistoryReport(
    submitted && repo && file
      ? { repo, file, limit, mode }
      : { repo: '', file: '' } // desativa query se não submetido
  );

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
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

      {isLoading && <p>Carregando...</p>}
      {error && <p style={{ color: 'red' }}>{(error as Error).message}</p>}

      {data && (
        <pre style={{ background: '#f5f5f5', padding: '1rem' }}>
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  );
};
