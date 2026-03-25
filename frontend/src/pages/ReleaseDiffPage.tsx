import React, { useState } from 'react';
import { GitHubTokenConfig } from '../components/GitHubTokenConfig';
import { apiFetch, ApiError } from '../utils/apiClient';
import { setGitHubPAT } from '../utils/githubToken';

export const ReleaseDiffPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [authError, setAuthError] = useState(false);

  // Estado estratégico: relatório pronto para exportar
  const [isReportReady, setIsReportReady] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setAuthError(false);
    setLoading(true);
    setData(null);
    setIsReportReady(false); // novo ciclo de geração → ainda não está pronto

    try {
      const params = new URLSearchParams({
        repo,
        from,
        to,
        format: 'json',
      });

      const res = await apiFetch(`/api/reports/release-diff?` + params.toString());
      const json = await res.json();
      setData(json);
      setIsReportReady(true); // sucesso → relatório pronto para exportação
    } catch (err) {
      if (err instanceof ApiError && err.isAuthError) {
        setAuthError(true);
        setError(err.message);
      } else {
        setError((err as Error).message);
      }
      setIsReportReady(false);
    } finally {
      setLoading(false);
    }
  };

  // Reset estratégico: qualquer mudança de parâmetro invalida o relatório atual
  const handleRepoChange = (value: string) => {
    setRepo(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
  };

  const handleFromChange = (value: string) => {
    setFrom(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
  };

  const handleToChange = (value: string) => {
    setTo(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
  };

  const canDownload = !!repo && !!from && !!to && isReportReady;

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

  const handleClearPAT = () => {
    setGitHubPAT(null);
    setAuthError(false);
    setError(null);
    setIsReportReady(false);
  };

  return (
    <div style={{ padding: '1rem' }}>
      <GitHubTokenConfig />
      <h1>Relatório 2 - PRs entre tags</h1>

      <form onSubmit={handleSubmit} style={{ marginBottom: '1rem' }}>
        <div>
          <label>Repo (owner/repo): </label>
          <input
            value={repo}
            onChange={e => handleRepoChange(e.target.value)}
            placeholder="org/projeto-x"
            style={{ width: '300px' }}
          />
        </div>
        <div>
          <label>From (tag/branch/SHA): </label>
          <input
            value={from}
            onChange={e => handleFromChange(e.target.value)}
            placeholder="v1.2.3"
            style={{ width: '200px' }}
          />
        </div>
        <div>
          <label>To (tag/branch/SHA): </label>
          <input
            value={to}
            onChange={e => handleToChange(e.target.value)}
            placeholder="v1.3.0"
            style={{ width: '200px' }}
          />
        </div>

        <button type="submit" style={{ marginTop: '0.5rem' }}>
          Gerar (JSON)
        </button>
      </form>

      <div style={{ marginBottom: '1rem', display: 'flex', gap: '0.5rem' }}>
        <button
          type="button"
          onClick={() => openFormat('markdown')}
          disabled={!canDownload}
          title={
            !canDownload ? 'Gere o relatório com sucesso para exportar em Markdown' : ''
          }
        >
          Abrir em Markdown
        </button>
        <button
          type="button"
          onClick={() => openFormat('csv')}
          disabled={!canDownload}
          title={
            !canDownload ? 'Gere o relatório com sucesso para exportar em CSV' : ''
          }
        >
          Baixar CSV
        </button>
      </div>

      {loading && <p>Carregando...</p>}
      {error && (
        <div
          style={{
            marginTop: '0.5rem',
            padding: '0.5rem',
            border: '1px solid #f99',
            background: '#fee',
          }}
        >
          <p style={{ color: '#a00', whiteSpace: 'pre-wrap' }}>{error}</p>
          {authError && (
            <button type="button" onClick={handleClearPAT}>
              Limpar PAT salvo
            </button>
          )}
        </div>
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
