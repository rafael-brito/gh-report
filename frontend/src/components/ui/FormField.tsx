import React from 'react';
import { theme } from '../../theme';

interface FormFieldProps {
  label: string;
  children: React.ReactNode; // input/select/etc.
  helpText?: string;
  style?: React.CSSProperties;
}

export const FormField: React.FC<FormFieldProps> = ({ label, children, helpText, style }) => {
  return (
    <div style={{ marginBottom: '0.75rem', ...style }}>
      <label
        style={{
          display: 'block',
          marginBottom: '0.25rem',
          fontSize: '0.9rem',
          color: theme.colors.textDark,
          fontWeight: 500,
        }}
      >
        {label}
      </label>
      {children}
      {helpText && (
        <div style={{ marginTop: '0.25rem', fontSize: '0.8rem', color: '#666' }}>{helpText}</div>
      )}
    </div>
  );
};

export const TextInput: React.FC<
  React.InputHTMLAttributes<HTMLInputElement>
> = props => {
  return (
    <input
      {...props}
      style={{
        width: '100%',
        padding: '0.4rem 0.5rem',
        borderRadius: 4,
        border: `1px solid ${theme.colors.border}`,
        fontSize: '0.9rem',
        ...props.style,
      }}
    />
  );
};

export const SelectInput: React.FC<
  React.SelectHTMLAttributes<HTMLSelectElement>
> = props => {
  return (
    <select
      {...props}
      style={{
        width: '100%',
        padding: '0.4rem 0.5rem',
        borderRadius: 4,
        border: `1px solid ${theme.colors.border}`,
        fontSize: '0.9rem',
        background: theme.colors.white,
        ...props.style,
      }}
    />
  );
};
