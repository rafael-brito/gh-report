import React, { useState } from 'react';
import { useFileHistoryReport, getFileHistoryDownloadUrl } from '../api/reports';
import { ApiError } from '../utils/apiClient';
import { setGitHubPAT } from '../utils/githubToken';
import { Card } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { FormField, TextInput, SelectInput } from '../components/ui/FormField';
import { theme } from '../theme';

export const FileHistoryPage: React.FC = () => {
  const [repo, setRepo] = useState('');
  const [file, setFile] = useState('');
  const [limit, setLimit] = useState(10);
  const [mode, setMode] = useState<'commits' | 'prs'>('prs');
  const [submitted, setSubmitted] = useState(false);

  const [isReportReady, setIsReportReady] = useState(false);
  const [validationError, setValidationError] = useState<string | null>(null);

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

  if (error && isReportReady) {
    setIsReportReady(false);
  }

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // validação leve antes de disparar o hook
    if (!repo || !file) {
      setValidationError('Preencha repo e file path antes de gerar o relatório.');
      setSubmitted(false);
      setIsReportReady(false);
      return;
    }
    if (!repo.includes('/')) {
      setValidationError('Repo deve estar no formato "owner/repo", por exemplo: openai/openai-python.');
      setSubmitted(false);
      setIsReportReady(false);
      return;
    }

    setValidationError(null);
    setSubmitted(true);
    setIsReportReady(false);
  };

  const handleRepoChange = (value: string) => {
    setRepo(value);
    setSubmitted(false);
    setIsReportReady(false);
    setValidationError(null);
  };

  const handleFileChange = (value: string) => {
    setFile(value);
    setSubmitted(false);
    setIsReportReady(false);
    setValidationError(null);
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
    setIsReportReady(false);
  };

  const items =
    data && typeof data === 'object'
      ? ((data as any).prs as any[]) || ((data as any).commits as any[])
      : null;
  const itemsCount = Array.isArray(items) ? items.length : 0;
  const hasNoResults = isReportReady && itemsCount === 0;

  return (
    <div>
      <h1 style={{ marginBottom: '0.75rem', color: theme.colors.textDark }}>
        Relatório de Histórico de Arquivo
      </h1>

      <Card style={{ marginBottom: '1rem', maxWidth: 640 }}>
        <form onSubmit={onSubmit}>
          <FormField label="Repo (owner/repo):">
            <TextInput
              value={repo}
              onChange={e => handleRepoChange(e.target.value)}
              placeholder="openai/openai-python"
            />
          </FormField>

          <FormField label="File path:">
            <TextInput
              value={file}
              onChange={e => handleFileChange(e.target.value)}
              placeholder="src/main/.../MinhaClasse.java"
            />
          </FormField>

          <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
            <FormField label="Limit:" style={{ flex: '0 0 120px' }}>
              <TextInput
                type="number"
                value={limit}
                onChange={e => handleLimitChange(Number(e.target.value) || 10)}
              />
            </FormField>

            <FormField label="Mode:" style={{ flex: '0 0 160px' }}>
              <SelectInput
                value={mode}
                onChange={e => handleModeChange(e.target.value as 'commits' | 'prs')}
              >
                <option value="prs">PRs</option>
                <option value="commits">Commits</option>
              </SelectInput>
            </FormField>
          </div>

          {validationError && (
            <p style={{ color: '#a00', fontSize: '0.85rem', marginTop: '0.25rem' }}>
              {validationError}
            </p>
          )}

          <Button type="submit" style={{ marginTop: '0.5rem' }}>
            Gerar
          </Button>
        </form>
      </Card>

      <div style={{ marginBottom: '0.75rem', display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
        <Button
          type="button"
          onClick={handleOpenMarkdown}
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
          onClick={handleOpenCSV}
          disabled={!canDownload}
          title={!canDownload ? 'Gere o relatório com sucesso para exportar em CSV' : ''}
        >
          Baixar CSV
        </Button>

        {isReportReady && (
          <span style={{ fontSize: '0.85rem', color: '#666' }}>
            {hasNoResults
              ? 'Relatório pronto (nenhum resultado encontrado).'
              : `Relatório pronto (${itemsCount} itens).`}
          </span>
        )}
      </div>

      {isLoading && <p>Carregando...</p>}

      {errorMessage && (
        <Card
          style={{
            marginTop: '0.5rem',
            maxWidth: 640,
            borderColor: '#f99',
            background: '#fee',
          }}
        >
          <p style={{ color: '#a00', whiteSpace: 'pre-wrap', margin: 0 }}>{errorMessage}</p>
          {isAuthErr && (
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
