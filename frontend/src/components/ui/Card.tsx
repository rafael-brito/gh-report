import React from 'react';
import { theme } from '../../theme';

interface CardProps {
  children: React.ReactNode;
  style?: React.CSSProperties;
}

export const Card: React.FC<CardProps> = ({ children, style }) => {
  return (
    <div
      style={{
        padding: '1rem',
        borderRadius: 8,
        background: theme.colors.white,
        border: `1px solid ${theme.colors.border}`,
        ...style,
      }}
    >
      {children}
    </div>
  );
};
