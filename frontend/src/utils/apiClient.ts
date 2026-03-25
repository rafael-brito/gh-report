import { getGitHubPAT } from './githubToken';

export class ApiError extends Error {
  status: number;
  isAuthError: boolean;

  constructor(message: string, status: number, isAuthError: boolean) {
    super(message);
    this.status = status;
    this.isAuthError = isAuthError;
  }
}

export async function apiFetch(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
  const pat = getGitHubPAT();

  const headers = new Headers(init.headers || {});
  if (pat) {
    headers.set('X-GitHub-Token', pat);
  }

  const res = await fetch(input, {
    ...init,
    headers,
  });

  if (!res.ok) {
    const text = await res.text().catch(() => '');
    const status = res.status;
    const isAuthError = status === 401 || status === 403;

    // Mensagem mais amigável para auth error
    const defaultMsg = isAuthError
      ? 'Token inválido ou expirado. Verifique suas configurações de PAT.'
      : `Erro: ${status} ${res.statusText}`;

    throw new ApiError(text || defaultMsg, status, isAuthError);
  }

  return res;
}
