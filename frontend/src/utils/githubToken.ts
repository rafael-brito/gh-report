const STORAGE_KEY = 'gh-report.github_pat';

export function getGitHubPAT(): string | null {
  const value = window.localStorage.getItem(STORAGE_KEY);
  if (!value) return null;
  return value;
}

export function setGitHubPAT(token: string | null) {
  if (!token) {
    window.localStorage.removeItem(STORAGE_KEY);
    return;
  }
  window.localStorage.setItem(STORAGE_KEY, token);
}
