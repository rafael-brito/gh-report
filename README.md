# gh-report

Geração de relatórios com base nos dados da GitHub API.

O objetivo é ajudar devs, squads e times de plataforma a enxergar melhor o histórico de código e o conteúdo de releases, sem precisar montar queries manuais na API ou ficar navegando por telas do GitHub.

Atualmente o projeto oferece **dois relatórios principais**:

- **Relatório 1 – Histórico de Arquivo (File History)**  
  Mostra commits ou PRs que tocaram um arquivo específico.
- **Relatório 2 – PRs entre tags (Release Diff)**  
  Mostra os PRs incluídos entre duas refs (tags/branches/SHAs), útil para montar notas de release.

---

## Relatórios

### 1. Relatório 1 – Histórico de Arquivo

Endpoint:

```http
GET /api/reports/file-history
```

Esse relatório responde perguntas do tipo:

- “Quais PRs tocaram o arquivo `path/to/file` recentemente?”
- “Quais commits mexeram nesse arquivo e quando?”

#### Parâmetros

Query params:

- `repo` (obrigatório)  
  Repositório no formato `owner/repo`, por exemplo: `org/projeto-x`.

- `file` (obrigatório)  
  Caminho do arquivo dentro do repositório, por exemplo:  
  `src/main/java/com/example/MyService.java`.

- `mode` (opcional, default: `prs`)  
  Modo de agregação:
  - `prs` – Agrupa por PRs que tocaram o arquivo.
  - `commits` – Lista commits diretamente.

- `limit` (opcional, default: `10`)  
  Limite de itens buscados na API do GitHub.

- `format` (opcional, default: `json`)  
  Formato de saída:
  - `json`
  - `markdown` (ou `md`)
  - `csv`

#### Exemplos de uso (curl)

**JSON (default)**

```bash
curl "http://localhost:8080/api/reports/file-history?repo=owner/repo&file=path/to/file.go"
```

**JSON explícito com modo e limite**

```bash
curl "http://localhost:8080/api/reports/file-history?repo=owner/repo&file=path/to/file.go&mode=prs&limit=20&format=json"
```

**Markdown**

```bash
curl "http://localhost:8080/api/reports/file-history?repo=owner/repo&file=path/to/file.go&mode=prs&format=markdown"
```

**CSV**

```bash
curl -o file-history.csv \
  "http://localhost:8080/api/reports/file-history?repo=owner/repo&file=path/to/file.go&mode=prs&format=csv"
```

---

### 2. Relatório 2 – PRs entre tags (Release Diff)

Endpoint:

```http
GET /api/reports/release-diff
```

Esse relatório responde perguntas do tipo:

- “Quais PRs entraram entre a `v1.2.3` e a `v1.3.0`?”
- “Quais tipos de mudança (feature, bugfix, technical, etc.) foram incluídos nesse intervalo?”

#### Parâmetros

Query params:

- `repo` (obrigatório)  
  Repositório no formato `owner/repo`.

- `from` (obrigatório)  
  Ref de base: pode ser uma tag, branch ou SHA.

- `to` (obrigatório)  
  Ref de destino: pode ser uma tag, branch ou SHA.

- `format` (opcional, default: `json`)  
  Formato de saída:
  - `json`
  - `markdown` (ou `md`)
  - `csv`

#### Exemplos de uso (curl)

**JSON (default)**

```bash
curl "http://localhost:8080/api/reports/release-diff?repo=owner/repo&from=v1.2.3&to=v1.3.0"
```

**Markdown**

```bash
curl "http://localhost:8080/api/reports/release-diff?repo=owner/repo&from=v1.2.3&to=v1.3.0&format=markdown"
```

**CSV**

```bash
curl -o release-diff.csv \
  "http://localhost:8080/api/reports/release-diff?repo=owner/repo&from=v1.2.3&to=v1.3.0&format=csv"
```

Se `from` ou `to` forem inválidos (tag/ref inexistente), a API retorna `400 Bad Request` com mensagem descritiva.

---

## Autenticação com GitHub (PAT e fallback)

O backend chama a API do GitHub usando um token, com duas possibilidades:

1. **Token global (modo demo/dev)**  
   Variável de ambiente `GITHUB_TOKEN` (como já era no MVP).
2. **PAT por usuário (recomendado)**  
   O usuário fornece um Personal Access Token no frontend, que é enviado em cada request via header:

   ```http
   X-GitHub-Token: <seu_pat>
   ```

O fluxo é:

- Se o header `X-GitHub-Token` estiver presente na request:
  - O backend usa **esse token** para falar com a API do GitHub.
- Se não houver header:
  - O backend usa o `GITHUB_TOKEN` do ambiente (se configurado).

No frontend, existe um componente de configuração de PAT (`GitHubTokenConfig`) nas páginas:

- Home
- Relatório de Histórico de Arquivo
- Relatório de Release Diff

Esse componente:

- Lê/grava o PAT no `localStorage`.
- Envia o PAT automaticamente em todas as chamadas de API.
- Permite limpar o PAT se ele estiver inválido.

Quando o token é inválido ou não tem permissão:

- A API do backend retorna erro 401/403.
- O frontend mostra mensagem clara (“Token inválido ou expirado…”) e oferece um botão para limpar o PAT salvo.

---

## Como rodar o projeto

### Backend

Pré‑requisitos:

- Go instalado.
- Um token GitHub com permissão de leitura de repositórios (PAT ou classic token).

Rodando localmente:

```bash
cd backend
export GITHUB_TOKEN=<seu_token_global_opcional>
go run ./cmd/server/
```

Por padrão, o servidor sobe na porta `8080`:

```text
http://localhost:8080
```

Endpoints principais de API:

- `/api/reports/file-history`
- `/api/reports/release-diff`

### Frontend

Pré‑requisitos:

- Node.js (recomendado versão LTS).
- npm ou yarn.

Rodando localmente:

```bash
cd frontend
npm install
npm run dev
```

O frontend sobe normalmente em algo como:

```text
http://localhost:5173
```

(veja o output do `npm run dev` para a porta exata).

O frontend está configurado para falar com o backend em `http://localhost:8080` para os endpoints de `/api`.

---

## Arquitetura

### Visão geral de camadas

O projeto está dividido em backend (Go) e frontend (React), com as seguintes responsabilidades principais:

#### Backend

- `internal/githubclient`  
  Client HTTP para a API do GitHub (REST).  
  Responsável por:
  - Listar commits por arquivo.
  - Buscar PRs associadas a um commit.
  - Comparar ranges de commits (`base...head`).
  - Aplicar token correto (PAT ou `GITHUB_TOKEN`) em cada requisição.

- `internal/reports`  
  Camada de domínio dos relatórios.  
  Contém:
  - Modelos de relatório (`FileHistoryReport`, `ReleaseDiffReport`, etc.).
  - Serviços:
    - `FileHistoryService` – monta o relatório de histórico de arquivo.
    - `ReleaseDiffService` – monta o relatório de PRs entre tags.
  - Renderizadores:
    - Métodos `ToMarkdown()`, `ToCSV()` para cada relatório.

- `internal/api`  
  Camada HTTP (handlers).  
  Responsável por:
  - Parse de query params.
  - Validação básica de entrada.
  - Escolha do formato de saída (JSON/Markdown/CSV).
  - Uso do `TokenProvider` e `GitHubClientFactory` para criar serviços por requisição.
  - Handlers principais:
    - `HandleFileHistory`
    - `HandleReleaseDiff`

- `cmd/server`  
  Ponto de entrada da aplicação backend:
  - Lê variáveis de ambiente (ex: `PORT`, `GITHUB_TOKEN`).
  - Inicializa `TokenProvider` e `GitHubClientFactory`.
  - Monta o router com os handlers.

#### Frontend

- `src/components/GitHubTokenConfig.tsx`  
  Componente de UI para configurar o PAT (token do usuário), armazenado em `localStorage`.

- `src/utils/apiClient.ts` (ou equivalente)  
  Wrapper para `fetch`:
  - Lê o PAT do `localStorage`.
  - Envia o header `X-GitHub-Token` quando houver PAT.
  - Converte erros HTTP em `ApiError`, marcando 401/403 como erros de autenticação.

- `src/api/reports.ts`  
  Funções e hooks para chamar os endpoints de relatório:
  - `useFileHistoryReport(...)`
  - `useReleaseDiffReport(...)`
  - Helpers para montar URLs de download (Markdown/CSV).

- `src/pages/FileHistoryPage.tsx`  
  Página do Relatório 1 (Histórico de Arquivo):
  - Campos: `repo`, `file`, `limit`, `mode`.
  - Botão “Gerar”.
  - Exibe JSON resultante.
  - Botões de exportar em Markdown/CSV (habilitados apenas quando o relatório está “pronto”).

- `src/pages/ReleaseDiffPage.tsx`  
  Página do Relatório 2 (Release Diff):
  - Campos: `repo`, `from`, `to`.
  - Botão “Gerar (JSON)”.
  - Exibe JSON resultante.
  - Botões de exportar em Markdown/CSV (habilitados apenas quando o relatório está “pronto”).

---

## Fluxo de uma requisição (Arquitetura em movimento)

### Relatório 1 – Histórico de Arquivo

1. **Frontend**  
   - Usuário acessa `/file-history` (ou rota equivalente).
   - Preenche `repo` e `file`.
   - Clica em “Gerar”.
   - A página usa o hook `useFileHistoryReport`, que chama:

     ```http
     GET /api/reports/file-history?repo=owner/repo&file=path/to/file.go&mode=prs&limit=10&format=json
     ```

   - Se houver PAT salvo, o frontend envia `X-GitHub-Token`.

2. **API Handler (`HandleFileHistory`)**  
   - Valida `repo`, `file`, `limit`, `mode`, `format`.
   - Usa o `TokenProvider` para obter o token correto para a request.
   - Usa o `GitHubClientFactory` para criar um `githubclient.Client` com esse token.
   - Cria um `FileHistoryService` para esse client.
   - Chama `GetFileHistoryReport(ctx, params)`.

3. **Service (`FileHistoryService`)**  
   - Usa o `githubclient.Client` para:
     - Listar commits que tocaram o arquivo.
     - Opcionalmente, buscar PRs relacionadas.
   - Monta um `FileHistoryReport` com os dados.
   - Converte para o formato solicitado (ou deixa em struct para JSON).

4. **Resposta**  
   - Handler devolve:
     - JSON, ou
     - Markdown (`text/markdown`), ou
     - CSV (`text/csv`).

### Relatório 2 – PRs entre tags (Release Diff)

1. **Frontend**  
   - Usuário acessa `/release-diff`.
   - Preenche `repo`, `from`, `to`.
   - Clica em “Gerar (JSON)”.
   - A página chama:

     ```http
     GET /api/reports/release-diff?repo=owner/repo&from=v1.2.3&to=v1.3.0&format=json
     ```

   - Se houver PAT, envia `X-GitHub-Token`.

2. **API Handler (`HandleReleaseDiff`)**  
   - Valida `repo`, `from`, `to`, `format`.
   - Usa `TokenProvider` + `GitHubClientFactory` para criar um `githubclient.Client` com o token da request.
   - Cria um `ReleaseDiffService` com esse client.
   - Chama `GetReleaseDiffReport(ctx, params)`.

3. **Service (`ReleaseDiffService`)**  

   Internamente:

   - Chama `CompareCommits` no `githubclient`:
     - Usa endpoint GitHub `/repos/{owner}/{repo}/compare/{from}...{to}`.
   - Para cada commit retornado:
     - Chama `ListPRsByCommit`:
       - `/repos/{owner}/{repo}/commits/{sha}/pulls`.
   - Deduplica PRs.
   - Classifica cada PR por tipo (feature, bugfix, technical, improvement, unknown) com base em labels.
   - Monta `ReleaseDiffReport` com:
     - Lista de PRs.
     - Summary (`total_prs`, contagem por tipo).
   - Devolve para o handler.

4. **Resposta**  
   - Handler devolve JSON, Markdown ou CSV conforme `format`.
