import React, { useState } from 'react';
import { apiFetch, ApiError } from '../utils/apiClient';
import { setGitHubPAT } from '../utils/githubToken';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { FormField, TextInput } from '../components/ui/FormField';
import { theme } from '../theme';

export const ReleaseDiffPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [authError, setAuthError] = useState(false);
  const [isReportReady, setIsReportReady] = useState(false);
  const [validationError, setValidationError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!repo || !from || !to) {
      setValidationError('Preencha repo, from e to antes de gerar o relatório.');
      setIsReportReady(false);
      return;
    }
    if (!repo.includes('/')) {
      setValidationError('Repo deve estar no formato "owner/repo", por exemplo: openai/openai-python.');
      setIsReportReady(false);
      return;
    }

    setValidationError(null);
    setError(null);
    setAuthError(false);
    setLoading(true);
    setData(null);
    setIsReportReady(false);

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
      setIsReportReady(true);
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

  const handleRepoChange = (value: string) => {
    setRepo(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
    setValidationError(null);
  };

  const handleFromChange = (value: string) => {
    setFrom(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
    setValidationError(null);
  };

  const handleToChange = (value: string) => {
    setTo(value);
    setIsReportReady(false);
    setError(null);
    setAuthError(false);
    setValidationError(null);
  };

  const canDownload = !!repo && !!from && !!to && isReportReady;

  const prs = data && (data as any).prs ? ((data as any).prs as any[]) : [];
  const prsCount = Array.isArray(prs) ? prs.length : 0;
  const hasNoResults = isReportReady && prsCount === 0;

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
    <div>
      <h1 style={{ marginBottom: '0.75rem', color: theme.colors.textDark }}>
        Relatório 2 – PRs entre tags
      </h1>

      <Card style={{ marginBottom: '1rem', maxWidth: 640 }}>
        <form onSubmit={handleSubmit}>
          <FormField label="Repo (owner/repo):">
            <TextInput
              value={repo}
              onChange={e => handleRepoChange(e.target.value)}
              placeholder="openai/openai-python"
            />
          </FormField>

          <FormField label="From (tag/branch/SHA):">
            <TextInput
              value={from}
              onChange={e => handleFromChange(e.target.value)}
              placeholder="v1.2.3"
            />
          </FormField>

          <FormField label="To (tag/branch/SHA):">
            <TextInput
              value={to}
              onChange={e => handleToChange(e.target.value)}
              placeholder="v1.3.0"
            />
          </FormField>

          {validationError && (
            <p style={{ color: '#a00', fontSize: '0.85rem', marginTop: '0.25rem' }}>
              {validationError}
            </p>
          )}

          <Button type="submit" style={{ marginTop: '0.5rem' }}>
            Gerar (JSON)
          </Button>
        </form>
      </Card>

      <div style={{ marginBottom: '0.75rem', display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
        <Button
          type="button"
          onClick={() => openFormat('markdown')}
          disabled={!canDownload}
          title={
            !canDownload ? 'Gere o relatório com sucesso para exportar em Markdown' : ''
          }
        >
          Abrir em Markdown
        </Button>
        <Button
          type="button"
          variant="secondary"
          onClick={() => openFormat('csv')}
          disabled={!canDownload}
          title={!canDownload ? 'Gere o relatório com sucesso para exportar em CSV' : ''}
        >
          Baixar CSV
        </Button>

        {isReportReady && (
          <span style={{ fontSize: '0.85rem', color: '#666' }}>
            {hasNoResults
              ? 'Relatório pronto (nenhum PR encontrado).'
              : `Relatório pronto (${prsCount} PRs).`}
          </span>
        )}
      </div>

      {loading && <p>Carregando...</p>}

      {error && (
        <Card
          style={{
            marginTop: '0.5rem',
            maxWidth: 640,
            borderColor: '#f99',
            background: '#fee',
          }}
        >
          <p style={{ color: '#a00', whiteSpace: 'pre-wrap', margin: 0 }}>{error}</p>
          {authError && (
            <div style={{ marginTop: '0.5rem' }}>
              <Button type="button" variant="secondary" onClick={handleClearPAT}>
                Limpar PAT salvo
              </Button>
            </div>
          )}
        </Card>
      )}

      {data && (
        <Card
          style={{
            marginTop: '0.75rem',
            maxWidth: 900,
            overflow: 'auto',
          }}
        >
          <pre
            style={{
              margin: 0,
              background: 'transparent',
              maxHeight: 400,
              overflow: 'auto',
            }}
          >
            {JSON.stringify(data, null, 2)}
          </pre>
        </Card>
      )}
    </div>
  );
};
