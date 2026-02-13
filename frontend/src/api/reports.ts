import { useQuery } from '@tanstack/react-query';

export interface FileHistoryQueryParams {
  repo: string;
  file: string;
  limit?: number;
  mode?: 'commits' | 'prs';
}

export function useFileHistoryReport(params: FileHistoryQueryParams) {
  const { repo, file, limit = 10, mode = 'prs' } = params;

  const enabled = !!repo && !!file;

  return useQuery({
    queryKey: ['file-history', { repo, file, limit, mode }],
    queryFn: async () => {
      const search = new URLSearchParams({
        repo,
        file,
        limit: String(limit),
        mode,
        format: 'json',
      });

      const res = await fetch(`/api/reports/file-history?` + search.toString());
      if (!res.ok) {
        throw new Error(`Erro ao buscar relatório: ${res.status} ${res.statusText}`);
      }
      return res.json();
    },
    enabled,
    staleTime: 60_000, // 1min
  });
}
