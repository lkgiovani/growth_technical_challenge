# Growth Technical Challenge

API RESTful desenvolvida em Go usando arquitetura limpa com PostgreSQL, Redis e Flyway.

## 🚀 Stack Tecnológico

- **Linguagem**: Go 1.23+
- **Framework HTTP**: Gin
- **ORM**: GORM
- **Banco de Dados**: PostgreSQL 16
- **Migrações**: Flyway 10
- **Cache**: Redis 7
- **Documentação**: TypeSpec + OpenAPI 3.1 (Swagger UI, ReDoc, Scalar)
- **Containerização**: Docker + Docker Compose

## 🐳 Como Subir o Ambiente com Docker

### 1. Clone o repositório

```bash
git clone https://github.com/lkgiovani/growth_technical_challenge.git
cd growth_technical_challenge
```

### 2. Configure as variáveis de ambiente

```bash
cp env.example .env
```

### 3. Inicie todos os serviços

Com um único comando, o Docker Compose orquestrará todos os serviços:

1. **PostgreSQL** - Inicia e aguarda até estar saudável
2. **Flyway** - Executa as migrações automaticamente
3. **Redis** - Inicia o serviço de cache
4. **Aplicação** - Inicia após as migrações serem concluídas

```bash
docker-compose up -d
```

### 4. Verifique se está rodando

```bash
curl http://localhost:8080/health
```

Resposta esperada:

```json
{
  "status": "ok"
}
```

## 🗄️ Como Rodar Migrations com Flyway

### Automático

As migrações são executadas automaticamente quando você faz `docker-compose up`.

### Manual

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

### Verificar histórico de migrações no banco

```bash
docker exec -it growth_postgres psql -U postgres -d growth_db
```

```sql
SELECT * FROM flyway_schema_history;
```

### Criar nova migração

1. Crie um arquivo em `infra/database/migrations/`:

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

## 📖 Como Acessar a Documentação Swagger

Após iniciar a aplicação com `docker-compose up -d`, acesse no seu navegador:

### Página Principal (Todas as Documentações)

- **URL**: http://localhost:8080/
- Página HTML com links para todas as documentações disponíveis

### Swagger UI (Recomendado)

- **URL**: http://localhost:8080/docs/swagger
- Interface interativa para testar endpoints diretamente

### ReDoc

- **URL**: http://localhost:8080/docs/redoc
- Interface moderna para leitura da documentação

### Scalar

- **URL**: http://localhost:8080/docs/scalar
- UI moderna com exemplos de código

### OpenAPI Spec (YAML)

- **URL**: http://localhost:8080/docs/openapi.yaml
- Especificação OpenAPI para importar em ferramentas

## 🧪 Exemplos de Requests

> **💡 Nota Importante**:
>
> - Para **Departamentos**: Você pode fornecer `gerente_id` (gerente existente) OU `gerente` (criar novo gerente)
> - Para **Colaboradores**: Você pode fornecer `departamento_id` (departamento existente) OU `departamento` (criar novo departamento)
> - **Nunca envie ambos os campos ao mesmo tempo!**

### Via cURL

#### 1. Health Check

```bash
curl http://localhost:8080/health
```

#### 2. Criar Departamento

##### Opção 1: Com gerente existente

```bash
curl -X POST http://localhost:8080/api/v1/departamentos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "TI",
    "gerente_id": "01234567-89ab-7def-0123-456789abcdef"
  }'
```

##### Opção 2: Criando gerente junto com o departamento

```bash
curl -X POST http://localhost:8080/api/v1/departamentos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Recursos Humanos",
    "gerente": {
      "nome": "Maria Santos",
      "cpf": "85165167097",
      "rg": "987654321"
    },
    "departamento_superior_id": "01234567-89ab-7def-0123-456789abcdef"
  }'
```

#### 3. Criar Colaborador

##### Opção 1: Com departamento existente

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

##### Opção 2: Criando departamento junto com o colaborador

```bash
curl -X POST http://localhost:8080/api/v1/colaboradores \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Carlos Souza",
    "cpf": "98765432101",
    "rg": "111222333",
    "departamento": {
      "nome": "Vendas",
      "gerente_id": "01234567-89ab-7def-0123-456789abcdef",
      "departamento_superior_id": "22222222-2222-2222-2222-222222222222"
    }
  }'
```

#### 4. Buscar Colaborador por ID

```bash
curl http://localhost:8080/api/v1/colaboradores/01234567-89ab-7def-0123-456789abcdef
```

#### 5. Atualizar Colaborador

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

#### 6. Listar Colaboradores com Filtros

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

#### 7. Deletar Colaborador (Soft Delete)

```bash
curl -X DELETE http://localhost:8080/api/v1/colaboradores/01234567-89ab-7def-0123-456789abcdef
```

#### 8. Buscar Departamento com Hierarquia

```bash
curl http://localhost:8080/api/v1/departamentos/01234567-89ab-7def-0123-456789abcdef
```

#### 9. Listar Departamentos com Filtros

```bash
curl -X POST http://localhost:8080/api/v1/departamentos/list \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "nome": "TI",
      "nome_gerente": "João"
    },
    "page": 1,
    "limit": 10
  }'
```

#### 10. Buscar Colaboradores do Gerente

```bash
curl http://localhost:8080/api/v1/gerentes/01234567-89ab-7def-0123-456789abcdef/colaboradores
```

### Via Postman

1. Abra o Postman
2. Clique em **Import** → **Link**
3. Digite: `http://localhost:8080/docs/openapi.yaml`
4. Clique em **Continue** → **Import**
5. Todos os endpoints estarão disponíveis com exemplos pré-configurados

### Via Insomnia

1. Abra o Insomnia
2. Clique em **Create** → **Import From** → **URL**
3. Digite: `http://localhost:8080/docs/openapi.yaml`
4. Clique em **Fetch and Import**
5. Todos os endpoints prontos para usar

### Via Swagger UI

1. Acesse http://localhost:8080/docs/swagger
2. Escolha um endpoint (ex: `POST /api/v1/colaboradores`)
3. Clique em **"Try it out"**
4. Preencha o corpo da requisição
5. Clique em **"Execute"**
6. Veja a resposta com status, headers e body

## 👤 Autor

**lkgiovani**

- GitHub: [@lkgiovani](https://github.com/lkgiovani)
