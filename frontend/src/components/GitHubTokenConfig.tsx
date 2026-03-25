import React, { useState } from 'react';
import { getGitHubPAT, setGitHubPAT } from '../utils/githubToken';
import { theme } from '../theme';

interface GitHubTokenConfigProps {
  onClose?: () => void;
}

export const GitHubTokenConfig: React.FC<GitHubTokenConfigProps> = ({ onClose }) => {
  const [value, setValue] = useState(() => getGitHubPAT() || '');
  const [status, setStatus] = useState<string | null>(null);

  const handleSave = () => {
    const trimmed = value.trim();
    setGitHubPAT(trimmed || null);
    setStatus(trimmed ? 'Token salvo localmente.' : 'Token removido.');
    setTimeout(() => setStatus(null), 2000);
  };

  const handleClear = () => {
    setValue('');
    setGitHubPAT(null);
    setStatus('Token removido.');
    setTimeout(() => setStatus(null), 2000);
  };

  return (
    <div
      style={{
        minWidth: 360,
        padding: '1rem',
        background: theme.colors.white,
        borderRadius: 8,
        boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
        border: `1px solid ${theme.colors.border}`,
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
        <strong style={{ color: theme.colors.textDark }}>Configuração de PAT (GitHub)</strong>
        {onClose && (
          <button
            type="button"
            onClick={onClose}
            style={{
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              fontSize: '1rem',
            }}
          >
            ✕
          </button>
        )}
      </div>

      <p style={{ fontSize: '0.85rem', color: '#666', marginBottom: '0.5rem' }}>
        Informe um GitHub Personal Access Token para usar nas chamadas da API.
        Se deixar vazio, o backend usará o <code>GITHUB_TOKEN</code> global (modo demo).
      </p>

      <input
        type="password"
        placeholder="ghp_..."
        value={value}
        onChange={e => setValue(e.target.value)}
        style={{
          width: '100%',
          padding: '0.5rem',
          marginBottom: '0.5rem',
          borderRadius: 4,
          border: `1px solid ${theme.colors.border}`,
        }}
      />

      <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.25rem' }}>
        <button
          type="button"
          onClick={handleSave}
          style={{
            flex: 1,
            padding: '0.4rem 0.75rem',
            borderRadius: 4,
            border: 'none',
            background: theme.colors.primary,
            color: theme.colors.white,
            cursor: 'pointer',
          }}
        >
          Salvar
        </button>
        <button
          type="button"
          onClick={handleClear}
          style={{
            padding: '0.4rem 0.75rem',
            borderRadius: 4,
            border: `1px solid ${theme.colors.border}`,
            background: theme.colors.white,
            cursor: 'pointer',
          }}
        >
          Limpar
        </button>
      </div>

      {status && (
        <p style={{ fontSize: '0.8rem', color: theme.colors.success, marginTop: '0.25rem' }}>
          {status}
        </p>
      )}
    </div>
  );
};
