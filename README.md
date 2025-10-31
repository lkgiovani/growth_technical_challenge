# Growth Technical Challenge

API RESTful desenvolvida em Go usando boas práticas e tecnologias modernas com arquitetura em camadas.

## 🚀 Stack Tecnológico

- **Linguagem**: Go 1.23+
- **Framework HTTP**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Banco de Dados**: PostgreSQL 16
- **Migrações**: Flyway 10
- **Cache**: Redis 7
- **Documentação**: TypeSpec + OpenAPI 3.1 (com ReDoc, Swagger UI, Scalar)
- **Injeção de Dependências**: Uber FX
- **Containerização**: Docker + Docker Compose

## 🚀 Quick Start

```bash
# 1. Clone and enter directory
git clone https://github.com/lkgiovani/growth_technical_challenge.git
cd growth_technical_challenge

# 2. Copy environment file
cp env.example .env

# 3. Start everything with one command
docker-compose up -d

# 4. Wait a few seconds and check health
curl http://localhost:8080/health

# 5. Access documentation
# Open http://localhost:8080/docs/swagger in your browser
```

That's it! The application is ready with:

- ✅ PostgreSQL database running
- ✅ All migrations applied by Flyway
- ✅ Redis cache ready
- ✅ API running on port 8080

## 📐 Arquitetura

O projeto segue uma arquitetura limpa em camadas:

- **Models**: Entidades de domínio e schemas do banco de dados
- **Repositories**: Camada de acesso a dados com interfaces
- **Services**: Lógica de negócio e regras de validação
- **Handlers**: Manipulação de requisições/respostas HTTP
- **Routes**: Definições dos endpoints da API

## 📁 Estrutura do Projeto

```
.
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── database.go          # Database configurations
├── internal/
│   ├── database/
│   │   └── database.go      # Connection and migrations
│   ├── handlers/
│   │   ├── user_handler.go
│   │   ├── employee_handler.go
│   │   ├── department_handler.go
│   │   └── manager_handler.go
│   ├── models/
│   │   ├── user.go
│   │   ├── employee.go
│   │   └── department.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── employee_repository.go
│   │   └── department_repository.go
│   ├── services/
│   │   ├── user_service.go
│   │   ├── employee_service.go
│   │   └── department_service.go
│   ├── routes/
│   │   └── routes.go
│   └── utils/
│       └── validators.go     # CPF and RG validators
├── docs/
│   └── docs.go              # Swagger documentation
├── docker-compose.yml       # Container orchestration
├── Dockerfile               # Application Docker image
├── Makefile                 # Facilitated commands
└── go.mod                   # Project dependencies
```

## 🔧 Pré-requisitos

- Docker & Docker Compose
- Go 1.23+ (para desenvolvimento local)
- Make (opcional, para comandos facilitados)

## 🐳 Executando com Docker (Recomendado)

### 1. Clone o repositório

```bash
git clone https://github.com/lkgiovani/growth_technical_challenge.git
cd growth_technical_challenge
```

### 2. Configure as variáveis de ambiente

```bash
cp env.example .env
```

Edite o arquivo `.env` conforme necessário.

### 3. Inicie todos os serviços

Com um **único comando**, o Docker Compose orquestrará todos os serviços na ordem correta:

1. **PostgreSQL** - Inicia e aguarda até estar saudável
2. **Flyway** - Executa as migrações do banco automaticamente
3. **Redis** - Inicia o serviço de cache
4. **Aplicação** - Inicia apenas após as migrações serem concluídas com sucesso

```bash
docker-compose up -d
```

Ou usando Make:

```bash
make up
```

Para acompanhar o processo de inicialização em tempo real (recomendado na primeira execução):

```bash
docker-compose up
```

### 4. Verifique se os serviços estão rodando

Verifique o status de todos os serviços:

```bash
docker-compose ps
```

Verifique os logs da aplicação:

```bash
docker-compose logs -f app
```

Verifique os logs das migrações do Flyway:

```bash
docker-compose logs flyway
```

Ou usando Make:

```bash
make logs
make logs-flyway
```

### 5. Acesse a aplicação

A API estará disponível em:

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Documentação**:
  - ReDoc: http://localhost:8080/docs/redoc
  - Swagger UI: http://localhost:8080/docs/swagger
  - Scalar: http://localhost:8080/docs/scalar
  - OpenAPI Spec: http://localhost:8080/docs/openapi.yaml
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 💻 Desenvolvimento Local

### 1. Instale as dependências

```bash
go mod download
```

### 2. Configure o banco de dados

Certifique-se de ter o PostgreSQL rodando e configure as variáveis de ambiente:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=growth_db
export DB_SSLMODE=disable
export APP_PORT=8080
```

### 3. Execute a aplicação

```bash
go run cmd/main.go
```

## 📊 Modelo de Domínio

### Colaborador (Employee)

- **id** (UUIDv7) - Chave primária
- **nome** (obrigatório) - Nome do colaborador
- **cpf** (obrigatório, único, validado) - CPF brasileiro
- **rg** (opcional, único se fornecido) - RG brasileiro
- **departamento_id** (FK para Departamento, obrigatório) - Referência do departamento

### Departamento (Department)

- **id** (UUIDv7) - Chave primária
- **nome** (obrigatório) - Nome do departamento
- **gerente_id** (FK para Colaborador, obrigatório) - Referência do gerente
- **departamento_superior_id** (FK opcional para Departamento) - Departamento pai para hierarquia

## 📚 Regras de Negócio

1. **CPF deve ser único e válido**
2. **RG, se fornecido, também deve ser único**
3. **O gerente deve ser um Colaborador existente vinculado ao mesmo departamento**
4. **O Departamento Superior é opcional, mas não pode criar ciclos na hierarquia**

## 🔌 Endpoints da API

### Health Check

```bash
GET /health
```

### Colaboradores

#### Criar Colaborador

```bash
POST /api/v1/employees
Content-Type: application/json

{
  "name": "John Doe",
  "cpf": "12345678901",
  "rg": "123456789",
  "department_id": "uuid-here"
}
```

#### Buscar Colaborador por ID

```bash
GET /api/v1/colaboradores/:id
```

Retorna o colaborador e o nome do gerente do departamento.

#### Atualizar Colaborador

```bash
PUT /api/v1/employees/:id
Content-Type: application/json

{
  "name": "John Doe Updated",
  "cpf": "12345678901",
  "rg": "123456789",
  "department_id": "uuid-here"
}
```

#### Deletar Colaborador

```bash
DELETE /api/v1/colaboradores/:id
```

Soft delete (exclusão lógica).

#### Listar Colaboradores com Filtros

```bash
POST /api/v1/employees/list
Content-Type: application/json

{
  "filters": {
    "name": "John",
    "cpf": "12345678901",
    "rg": "123456789",
    "department_id": "uuid-here"
  },
  "page": 1,
  "limit": 10
}
```

### Departamentos

#### Criar Departamento

```bash
POST /api/v1/departamentos
Content-Type: application/json

{
  "nome": "TI",
  "gerente_id": "uuid-aqui",
  "departamento_superior_id": "uuid-aqui" // opcional
}
```

Valida gerente_id e previne ciclos.

#### Buscar Departamento por ID

```bash
GET /api/v1/departamentos/:id
```

Retorna departamento, gerente e árvore hierárquica completa de subdepartamentos.

#### Atualizar Departamento

```bash
PUT /api/v1/departamentos/:id
Content-Type: application/json

{
  "nome": "TI Atualizado",
  "gerente_id": "uuid-aqui",
  "departamento_superior_id": "uuid-aqui"
}
```

Previne ciclos na hierarquia.

#### Deletar Departamento

```bash
DELETE /api/v1/departamentos/:id
```

#### Listar Departamentos com Filtros

```bash
POST /api/v1/departamentos/list
Content-Type: application/json

{
  "filters": {
    "nome": "TI",
    "nome_gerente": "João",
    "departamento_superior_id": "uuid-aqui"
  },
  "page": 1,
  "limit": 10
}
```

### Gerentes

#### Buscar Colaboradores do Gerente

```bash
GET /api/v1/gerentes/:id/colaboradores
```

Retorna todos os colaboradores dos departamentos subordinados do gerente, recursivamente.

## ⚖️ Códigos de Status HTTP

A API usa respostas de erro consistentes:

- **200 OK** - Requisição bem-sucedida
- **201 Created** - Recurso criado com sucesso
- **400 Bad Request** - Filtros inválidos ou requisição malformada
- **404 Not Found** - Recurso não encontrado
- **409 Conflict** - Violação de restrição de unicidade (CPF, RG)
- **422 Unprocessable Entity** - Erro de validação de domínio (CPF inválido, regras de negócio)
- **500 Internal Server Error** - Erro do servidor

## 🧪 Testing the API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

### 2. Create Department

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "TI",
    "gerente_id": "01234567-89ab-7def-0123-456789abcdef"
  }'
```

### 3. Criar Colaborador

```bash
curl -X POST http://localhost:8080/api/v1/colaboradores \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva",
    "cpf": "12345678901",
    "rg": "123456789",
    "departamento_id": "01234567-89ab-7def-0123-456789abcdef"
  }'
```

### 4. Buscar Colaborador por ID

```bash
curl http://localhost:8080/api/v1/colaboradores/01234567-89ab-7def-0123-456789abcdef
```

### 5. Atualizar Colaborador

```bash
curl -X PUT http://localhost:8080/api/v1/colaboradores/01234567-89ab-7def-0123-456789abcdef \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva Updated",
    "cpf": "12345678901",
    "rg": "987654321",
    "departamento_id": "01234567-89ab-7def-0123-456789abcdef"
  }'
```

### 6. Listar Colaboradores com Filtros

```bash
curl -X POST http://localhost:8080/api/v1/colaboradores/list \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "nome": "João"
    },
    "page": 1,
    "limit": 10
  }'
```

### 7. Deletar Colaborador (Soft Delete)

```bash
curl -X DELETE http://localhost:8080/api/v1/colaboradores/01234567-89ab-7def-0123-456789abcdef
```

### 8. Buscar Departamento com Hierarquia

```bash
curl http://localhost:8080/api/v1/departamentos/01234567-89ab-7def-0123-456789abcdef
```

### 9. Buscar Colaboradores do Gerente

```bash
curl http://localhost:8080/api/v1/gerentes/01234567-89ab-7def-0123-456789abcdef/colaboradores
```

### Using Postman or Insomnia

1. **Import OpenAPI Spec**:
   - URL: `http://localhost:8080/docs/openapi.yaml`
2. **Set Base URL**:

   - `http://localhost:8080`

3. **Use Interactive Swagger UI**:
   - Access: http://localhost:8080/docs/swagger
   - Click "Try it out" on any endpoint
   - Fill in the parameters
   - Click "Execute"

## 🛠️ Comandos Make

Se você tiver o `make` instalado, pode usar estes atalhos:

```bash
make up          # Inicia todos os serviços em background
make down        # Para todos os serviços
make run         # Inicia todos os serviços e mostra logs
make logs        # Mostra logs da aplicação
make logs-flyway # Mostra logs das migrações Flyway
make restart     # Reinicia todos os serviços
make clean       # Remove containers e volumes (início limpo)
make build       # Reconstrói imagens Docker
make rebuild     # Clean + build + start tudo
make status      # Mostra status de todos os serviços
```

## 🗄️ Banco de Dados & Migrações

### Migrações de Banco de Dados com Flyway

Este projeto usa **Flyway** para controle de versão do banco de dados e migrações. Todos os arquivos de migração estão localizados em `infra/database/migrations/`.

#### Convenção de Nomenclatura dos Arquivos de Migração

O Flyway usa um padrão específico de nomenclatura:

```
V{versão}__{descrição}.sql
```

Exemplos:

- `V1__create_departamentos.sql`
- `V2__create_colaboradores.sql`
- `V3__add_fk_colaboradores_departamento.sql`

#### Como as Migrações Funcionam

1. Quando você executa `docker-compose up`, o Flyway automaticamente:

   - Conecta ao PostgreSQL
   - Cria uma tabela `flyway_schema_history` para rastrear migrações
   - Executa migrações pendentes em ordem (V1, V2, V3, etc.)
   - Marca cada migração como aplicada
   - Pula migrações já aplicadas

2. A aplicação só inicia **após** todas as migrações serem concluídas com sucesso

#### Executar Migrações Manualmente

Para executar migrações separadamente:

```bash
docker-compose up flyway
```

Para ver o status das migrações:

```bash
docker-compose run --rm flyway info
```

Para validar migrações:

```bash
docker-compose run --rm flyway validate
```

#### Criar uma Nova Migração

1. Crie um novo arquivo em `infra/database/migrations/`:

   ```bash
   touch infra/database/migrations/V6__add_new_feature.sql
   ```

2. Escreva seu SQL:

   ```sql
   ALTER TABLE colaboradores ADD COLUMN email VARCHAR(255);
   ```

3. Reinicie os serviços:
   ```bash
   docker-compose restart flyway app
   ```

### Conectar ao PostgreSQL

```bash
docker exec -it growth_postgres psql -U postgres -d growth_db
```

Verificar histórico de migrações:

```sql
SELECT * FROM flyway_schema_history;
```

## 🔐 Variáveis de Ambiente

| Variável            | Descrição                | Padrão      |
| ------------------- | ------------------------ | ----------- |
| `POSTGRES_USER`     | Usuário PostgreSQL       | `postgres`  |
| `POSTGRES_PASSWORD` | Senha PostgreSQL         | `postgres`  |
| `POSTGRES_DB`       | Nome do banco de dados   | `growth_db` |
| `POSTGRES_PORT`     | Porta PostgreSQL         | `5432`      |
| `DB_HOST`           | Host do banco de dados   | `localhost` |
| `DB_PORT`           | Porta do banco de dados  | `5432`      |
| `DB_SSLMODE`        | Modo SSL PostgreSQL      | `disable`   |
| `APP_PORT`          | Porta da aplicação       | `8080`      |
| `GIN_MODE`          | Modo Gin (debug/release) | `debug`     |

## 📖 Documentação da API

Este projeto usa **TypeSpec** para gerar documentação OpenAPI 3.1, fornecendo múltiplos visualizadores interativos.

### 📚 Visualizadores de Documentação

Após iniciar a aplicação, acesse a documentação no seu navegador:

#### 1. **Swagger UI** (Recomendado para Testes)

- **URL**: http://localhost:8080/docs/swagger
- **Recursos**:
  - Explorador interativo da API
  - Teste endpoints diretamente no navegador
  - Veja exemplos de requisição/resposta
  - Não precisa de ferramentas adicionais

#### 2. **ReDoc**

- **URL**: http://localhost:8080/docs/redoc
- **Recursos**:
  - Interface moderna e limpa
  - Ótimo para ler e entender a API
  - Funcionalidade de busca
  - Amigável para impressão

#### 3. **Scalar**

- **URL**: http://localhost:8080/docs/scalar
- **Recursos**:
  - UI bonita e moderna
  - Exemplos de código em múltiplas linguagens
  - Playground da API

#### 4. **OpenAPI Spec (YAML)**

- **URL**: http://localhost:8080/docs/openapi.yaml
- **Use para**:
  - Importar no Postman/Insomnia
  - Gerar SDKs de cliente
  - Integração CI/CD

### 🔧 Como Usar o Swagger UI

1. Abra http://localhost:8080/docs/swagger no seu navegador

2. Escolha um endpoint (ex: `POST /api/v1/colaboradores`)

3. Clique em **"Try it out"**

4. Preencha o corpo da requisição:

   ```json
   {
     "nome": "João Silva",
     "cpf": "12345678901",
     "rg": "123456789",
     "departamento_id": "01234567-89ab-7def-0123-456789abcdef"
   }
   ```

5. Clique em **"Execute"**

6. Veja a resposta abaixo com código de status, headers e body

### 📥 Importar para o Postman

1. Abra o Postman

2. Clique em **Import** → **Link**

3. Digite: `http://localhost:8080/docs/openapi.yaml`

4. Clique em **Continue** → **Import**

5. Todos os endpoints estarão disponíveis no Postman com:
   - Requisições pré-configuradas
   - Payloads de exemplo
   - Descrições

### 📥 Importar para o Insomnia

1. Abra o Insomnia

2. Clique em **Create** → **Import From** → **URL**

3. Digite: `http://localhost:8080/docs/openapi.yaml`

4. Clique em **Fetch and Import**

5. Todos os endpoints prontos para usar!

### 🔄 Regenerar Documentação

Se você modificar os arquivos TypeSpec em `docs/`:

```bash
cd docs
npm install
npm run compile
```

Ou usando Make:

```bash
cd docs
make install
make compile
```

Então reinicie a aplicação:

```bash
docker-compose restart app
```

Veja [docs/README.md](docs/README.md) para mais informações sobre TypeSpec.

## 🚀 Build de Produção

```bash
# Build da aplicação
go build -o app cmd/main.go

# Executar
./app
```

Ou com Docker:

```bash
docker build -t growth-app .
docker run -p 8080:8080 growth-app
```

## ✅ Funcionalidades

- ✅ Arquitetura limpa em camadas (Models, Repositories, Services, Handlers)
- ✅ UUID v7 para chaves primárias
- ✅ Validação de CPF e RG
- ✅ Prevenção de ciclos na hierarquia de departamentos
- ✅ Soft delete em registros
- ✅ Paginação e filtragem
- ✅ Cache com Redis
- ✅ Respostas de erro consistentes
- ✅ Migrações de banco de dados com Flyway
- ✅ Documentação TypeSpec + OpenAPI 3.1
- ✅ Múltiplos visualizadores de documentação (Swagger UI, ReDoc, Scalar)
- ✅ Docker multi-stage build
- ✅ Orquestração de serviços com Docker Compose
- ✅ Injeção de dependências com Uber FX
- ✅ README completo

## 📝 Licença

Este projeto foi desenvolvido como um desafio técnico.

## 👤 Autor

**lkgiovani**

- GitHub: [@lkgiovani](https://github.com/lkgiovani)
