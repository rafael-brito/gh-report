import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { GitHubTokenConfig } from './GitHubTokenConfig';
import { theme } from '../theme';

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const [showPatConfig, setShowPatConfig] = useState(false);

  return (
    <div
      style={{
        minHeight: '100vh',
        background: theme.colors.background,
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      {/* Header */}
      <header
        style={{
          height: 56,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '0 1.5rem',
          background: theme.colors.textDark,
          color: theme.colors.white,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <div
            style={{
              width: 28,
              height: 28,
              borderRadius: '4px',
              background: theme.colors.primary,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 700,
              fontSize: '0.9rem',
            }}
          >
            gh
          </div>
          <div>
            <div style={{ fontWeight: 600 }}>gh-report</div>
            <div style={{ fontSize: '0.75rem', opacity: 0.8 }}>
              Relatórios sobre repositórios GitHub
            </div>
          </div>
        </div>

        {/* Navegação simples */}
        <nav style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <Link
            to="/"
            style={{ color: theme.colors.white, textDecoration: 'none', fontSize: '0.9rem' }}
          >
            Início
          </Link>
          <Link
            to="/file-history"
            style={{ color: theme.colors.white, textDecoration: 'none', fontSize: '0.9rem' }}
          >
            Histórico de Arquivo
          </Link>
          <Link
            to="/release-diff"
            style={{ color: theme.colors.white, textDecoration: 'none', fontSize: '0.9rem' }}
          >
            PRs entre tags
          </Link>

          {/* Botão de engrenagem */}
          <button
            type="button"
            onClick={() => setShowPatConfig(true)}
            title="Configurações de autenticação (PAT)"
            style={{
              marginLeft: '0.5rem',
              border: 'none',
              background: 'transparent',
              color: theme.colors.white,
              cursor: 'pointer',
              fontSize: '1.1rem',
            }}
          >
            ⚙
          </button>
        </nav>
      </header>

      {/* Conteúdo principal */}
      <main style={{ flex: 1, padding: '1.5rem' }}>{children}</main>

      {/* Modal simples para PAT */}
      {showPatConfig && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0, 0, 0, 0.35)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
          }}
          onClick={() => setShowPatConfig(false)}
        >
          <div onClick={e => e.stopPropagation()}>
            <GitHubTokenConfig onClose={() => setShowPatConfig(false)} />
          </div>
        </div>
      )}
    </div>
  );
};
