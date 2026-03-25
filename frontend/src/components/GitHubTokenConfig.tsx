import React, {useState } from 'react';
import { getGitHubPAT, setGitHubPAT } from '../utils/githubToken';

export const GitHubTokenConfig: React.FC = () => {
  const [value, setValue] = useState(() => getGitHubPAT() || '');
  const [status, setStatus] = useState<string | null>(null);

  const handleSave = () => {
    const trimmed = value.trim();
    setGitHubPAT(trimmed || null);
    setStatus(trimmed ? 'Token salvo localmente.' : 'Token removido.');
    setTimeout(() => setStatus(null), 2000);
  };

  return (
    <div style={{ marginBottom: '1rem', padding: '0.5rem', border: '1px solid #ddd' }}>
      <strong>GitHub Personal Access Token (opcional)</strong>
      <div style={{ marginTop: '0.5rem' }}>
        <input
          type="password"
          placeholder="ghp_..."
          value={value}
          onChange={e => setValue(e.target.value)}
          style={{ width: '320px', marginRight: '0.5rem' }}
        />
        <button type="button" onClick={handleSave}>
          Salvar
        </button>
      </div>
      <p style={{ fontSize: '0.8rem', color: '#666', marginTop: '0.25rem' }}>
        Se vazio, o backend usará o GITHUB_TOKEN de servidor (modo demo).
      </p>
      {status && (
        <p style={{ fontSize: '0.8rem', color: '#007700', marginTop: '0.25rem' }}>
          {status}
        </p>
      )}
    </div>
  );
};
