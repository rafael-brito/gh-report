import React from 'react';
import { Link } from 'react-router-dom';
import { theme } from '../theme';

export const HomePage: React.FC = () => {
  return (
    <div>
      <h1 style={{ marginBottom: '0.5rem', color: theme.colors.textDark }}>Painel de Relatórios</h1>
      <p style={{ marginBottom: '1.5rem', maxWidth: 620 }}>
        O <strong>gh-report</strong> gera relatórios a partir de repositórios GitHub, ajudando times a
        auditar histórico de arquivos e montar notas de release entre versões.
      </p>

      <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
        <Link
          to="/file-history"
          style={{
            flex: '1 1 260px',
            minWidth: 260,
            textDecoration: 'none',
            color: 'inherit',
          }}
        >
          <div
            style={{
              padding: '1rem',
              borderRadius: 8,
              background: theme.colors.white,
              border: `1px solid ${theme.colors.border}`,
            }}
          >
            <h2 style={{ marginTop: 0, color: theme.colors.textDark }}>Relatório 1 – Histórico de Arquivo</h2>
            <p style={{ fontSize: '0.9rem', marginBottom: '0.75rem' }}>
              Veja quais commits ou PRs tocaram um arquivo específico, com limite configurável e saída
              em JSON, Markdown ou CSV.
            </p>
            <span
              style={{
                display: 'inline-block',
                marginTop: '0.5rem',
                padding: '0.25rem 0.75rem',
                borderRadius: 999,
                background: theme.colors.primary,
                color: theme.colors.white,
                fontSize: '0.8rem',
              }}
            >
              Auditar histórico
            </span>
          </div>
        </Link>

        <Link
          to="/release-diff"
          style={{
            flex: '1 1 260px',
            minWidth: 260,
            textDecoration: 'none',
            color: 'inherit',
          }}
        >
          <div
            style={{
              padding: '1rem',
              borderRadius: 8,
              background: theme.colors.white,
              border: `1px solid ${theme.colors.border}`,
            }}
          >
            <h2 style={{ marginTop: 0, color: theme.colors.textDark }}>Relatório 2 – PRs entre tags</h2>
            <p style={{ fontSize: '0.9rem', marginBottom: '0.75rem' }}>
              Liste os PRs incluídos entre duas tags/refs (por exemplo, entre versões) e classifique por
              tipo: feature, bugfix, technical, improvement.
            </p>
            <span
              style={{
                display: 'inline-block',
                marginTop: '0.5rem',
                padding: '0.25rem 0.75rem',
                borderRadius: 999,
                background: theme.colors.primary,
                color: theme.colors.white,
                fontSize: '0.8rem',
              }}
            >
              Gerar release notes
            </span>
          </div>
        </Link>
      </div>
    </div>
  );
};