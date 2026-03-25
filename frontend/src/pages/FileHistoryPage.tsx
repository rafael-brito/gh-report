import React, { useState } from 'react';
import { useFileHistoryReport, getFileHistoryDownloadUrl } from '../api/reports';
import { GitHubTokenConfig } from '../components/GitHubTokenConfig';
import { ApiError } from '../utils/apiClient';
import { setGitHubPAT } from '../utils/githubToken';

export const FileHistoryPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [file, setFile] = useState('');
  const [limit, setLimit] = useState(10);
  const [mode, setMode] = useState<'commits' | 'prs'>('prs');
  const [submitted, setSubmitted] = useState(false);

  // Estado estratégico: só fica true após um "Gerar" bem-sucedido
  const [isReportReady, setIsReportReady] = useState(false);

  const queryParams =
    submitted && repo && file
      ? { repo, file, limit, mode }
      : { repo: '', file: '' };

  const { data, isLoading, error } = useFileHistoryReport(queryParams as any);

  let errorMessage: string | null = null;
  let isAuthErr = false;
  if (error) {
    const e = error as any;
    if (e instanceof ApiError && e.isAuthError) {
      errorMessage = e.message;
      isAuthErr = true;
    } else {
      errorMessage = (e as Error).message;
    }
  }

  // Sempre que o hook tiver erro, consideramos que o relatório NÃO está pronto
  if (error && isReportReady) {
    setIsReportReady(false);
  }

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
    // Ao gerar, ainda não sabemos se vai dar certo → marca como não pronto
    setIsReportReady(false);
  };

  // Reset estratégico: qualquer mudança de input invalida o "relatório pronto"
  const handleRepoChange = (value: string) => {
    setRepo(value);
    setSubmitted(false);
    setIsReportReady(false);
  };

  const handleFileChange = (value: string) => {
    setFile(value);
    setSubmitted(false);
    setIsReportReady(false);
  };

  const handleLimitChange = (value: number) => {
    setLimit(value);
    setSubmitted(false);
    setIsReportReady(false);
  };

  const handleModeChange = (value: 'commits' | 'prs') => {
    setMode(value);
    setSubmitted(false);
    setIsReportReady(false);
  };

  // Quando tivermos dados e nenhum erro, consideramos que o relatório está pronto
  const hasValidData = submitted && !!repo && !!file && !!data && !error;
  if (hasValidData && !isReportReady) {
    setIsReportReady(true);
  }

  const canDownload = isReportReady;

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

  const handleClearPAT = () => {
    setGitHubPAT(null);
    // limpamos também o estado de erro e de "pronto"
    setIsReportReady(false);
  };

  return (
    <div style={{ padding: '1rem' }}>
      <GitHubTokenConfig />
      <h1>Relatório de Histórico de Arquivo</h1>

      <form onSubmit={onSubmit} style={{ marginBottom: '1rem' }}>
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
          <label>File path: </label>
          <input
            value={file}
            onChange={e => handleFileChange(e.target.value)}
            placeholder="src/main/.../MinhaClasse.java"
            style={{ width: '300px' }}
          />
        </div>
        <div>
          <label>Limit: </label>
          <input
            type="number"
            value={limit}
            onChange={e => handleLimitChange(Number(e.target.value) || 10)}
            style={{ width: '80px' }}
          />
        </div>
        <div>
          <label>Mode: </label>
          <select
            value={mode}
            onChange={e => handleModeChange(e.target.value as 'commits' | 'prs')}
          >
            <option value="prs">PRs</option>
            <option value="commits">Commits</option>
          </select>
        </div>
        <button type="submit" style={{ marginTop: '0.5rem' }}>
          Gerar
        </button>
      </form>

      <div style={{ marginBottom: '1rem', display: 'flex', gap: '0.5rem' }}>
        <button
          type="button"
          onClick={handleOpenMarkdown}
          disabled={!canDownload}
          title={
            !canDownload ? 'Gere o relatório com sucesso para exportar em Markdown' : ''
          }
        >
          Abrir em Markdown
        </button>
        <button
          type="button"
          onClick={handleOpenCSV}
          disabled={!canDownload}
          title={
            !canDownload ? 'Gere o relatório com sucesso para exportar em CSV' : ''
          }
        >
          Baixar CSV
        </button>
      </div>

      {isLoading && <p>Carregando...</p>}
      {errorMessage && (
        <div
          style={{
            marginTop: '0.5rem',
            padding: '0.5rem',
            border: '1px solid #f99',
            background: '#fee',
          }}
        >
          <p style={{ color: '#a00', whiteSpace: 'pre-wrap' }}>{errorMessage}</p>
          {isAuthErr && (
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
