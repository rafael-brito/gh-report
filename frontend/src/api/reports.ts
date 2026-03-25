import { useQuery } from '@tanstack/react-query';
import { apiFetch } from '../utils/apiClient';

export interface FileHistoryQueryParams {
  repo: string;
  file: string;
  limit?: number;
  mode?: 'commits' | 'prs';
}

function buildFileHistoryUrl(params: FileHistoryQueryParams, format: 'json' | 'markdown' | 'csv') {
  const { repo, file, limit = 10, mode = 'prs' } = params;

  const search = new URLSearchParams({
    repo,
    file,
    limit: String(limit),
    mode,
    format,
  });

  return `/api/reports/file-history?` + search.toString();
}

export function useFileHistoryReport(params: FileHistoryQueryParams) {
  const { repo, file, limit = 10, mode = 'prs' } = params;

  const enabled = !!repo && !!file;

  return useQuery({
    queryKey: ['file-history', { repo, file, limit, mode }],
    queryFn: async () => {
      const url = buildFileHistoryUrl({ repo, file, limit, mode }, 'json');
      const res = await apiFetch(url);
      if (!res.ok) {
        throw new Error(`Erro ao buscar relatório: ${res.status} ${res.statusText}`);
      }
      return res.json();
    },
    enabled,
    staleTime: 60_000, // 1min
  });
}

export function getFileHistoryDownloadUrl(
  params: FileHistoryQueryParams,
  format: 'markdown' | 'csv',
): string {
  return buildFileHistoryUrl(params, format);
}
