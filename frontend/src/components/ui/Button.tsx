import React from 'react';
import { theme } from '../../theme';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary';
}

export const Button: React.FC<ButtonProps> = ({ variant = 'primary', style, ...rest }) => {
  const base: React.CSSProperties = {
    padding: '0.4rem 0.9rem',
    borderRadius: 4,
    fontSize: '0.9rem',
    fontWeight: 500,
    cursor: 'pointer',
    border: '1px solid transparent',
    fontFamily: 'inherit',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '0.25rem',
  };

  let variantStyle: React.CSSProperties;

  if (variant === 'primary') {
    variantStyle = {
      background: theme.colors.primary,
      color: theme.colors.white,
      borderColor: theme.colors.primary,
    };
  } else {
    variantStyle = {
      background: theme.colors.white,
      color: theme.colors.textDark,
      borderColor: theme.colors.border,
    };
  }

  const disabledStyle: React.CSSProperties = rest.disabled
    ? {
        opacity: 0.6,
        cursor: 'not-allowed',
      }
    : {};

  return <button style={{ ...base, ...variantStyle, ...disabledStyle, ...style }} {...rest} />;
};
