import React from 'react';
import { Link } from 'react-router-dom';

export const HomePage: React.FC = () => {
  return (
    <div style={{ padding: '1rem' }}>
      <h1>GH Report</h1>
      <p>Ferramenta de relatórios de histórico de alterações e releases.</p>
      <ul>
        <li>
          <Link to="/file-history">Relatório de Histórico de Arquivo</Link>
        </li>
        <li>
          <Link to="/release-diff">Relatório de PRs entre Tags</Link>
        </li>
      </ul>
    </div>
  );
};