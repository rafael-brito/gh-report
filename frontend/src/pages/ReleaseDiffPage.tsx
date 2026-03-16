import React, { useState } from 'react';

export const ReleaseDiffPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canDownload = !!repo && !!from && !!to;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    setData(null);

    try {
      const params = new URLSearchParams({
        repo,
        from,
        to,
        format: 'json',
      });

      const res = await fetch(`/api/reports/release-diff?` + params.toString());
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Erro: ${res.status} ${res.statusText}`);
      }
      const json = await res.json();
      setData(json);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const openFormat = (format: 'markdown' | 'csv') => {
    if (!canDownload) return;
    const params = new URLSearchParams({
      repo,
      from,
      to,
      format,
    });
    window.open(`/api/reports/release-diff?` + params.toString(), '_blank');
  };

  return (
    <div style={{ padding: '1rem' }}>
      <h1>Relatório 2 - PRs entre tags</h1>

      <form onSubmit={handleSubmit} style={{ marginBottom: '1rem' }}>
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
          <label>From (tag/branch/SHA): </label>
          <input
            value={from}
            onChange={e => setFrom(e.target.value)}
            placeholder="v1.2.3"
            style={{ width: '200px' }}
          />
        </div>
        <div>
          <label>To (tag/branch/SHA): </label>
          <input
            value={to}
            onChange={e => setTo(e.target.value)}
            placeholder="v1.3.0"
            style={{ width: '200px' }}
          />
        </div>

        <button type="submit" style={{ marginTop: '0.5rem' }}>
          Gerar (JSON)
        </button>
      </form>

      {canDownload && (
        <div style={{ marginBottom: '1rem', display: 'flex', gap: '0.5rem' }}>
          <button type="button" onClick={() => openFormat('markdown')}>
            Abrir em Markdown
          </button>
          <button type="button" onClick={() => openFormat('csv')}>
            Baixar CSV
          </button>
        </div>
      )}

      {loading && <p>Carregando...</p>}
      {error && (
        <p style={{ color: 'red', whiteSpace: 'pre-wrap' }}>
          {error}
        </p>
      )}

      {data && (
        <pre
          style={{
            background: '#f5f5f5',
            padding: '1rem',
            maxHeight: '400px',
            overflow: 'auto',
          }}
        >
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  );
};
