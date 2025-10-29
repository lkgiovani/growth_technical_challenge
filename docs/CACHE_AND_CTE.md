# Cache e CTE Recursivo para Consultas Hierárquicas

## Visão Geral

Esta implementação adiciona cache com Redis e otimiza consultas hierárquicas usando Common Table Expressions (CTEs) recursivos do PostgreSQL para melhorar significativamente a performance de operações com departamentos.

## Arquitetura

### 1. Cache com Redis

#### Estrutura de Camadas

```
UseCase → CachedDepartmentRepository → BaseDepartmentRepository → Database
                 ↓
              Redis Cache
```

#### Implementação

**Redis Client** (`infra/cache/redis.go`)

- Conexão e gerenciamento do cliente Redis
- Operações básicas: Get, Set, Del, DelPattern
- TTL configurável via variável de ambiente
- Graceful degradation: aplicação continua funcionando sem cache se Redis não estiver disponível

**Cache Interface** (`pkg/cache/cache.go`)

- Interface genérica para operações de cache
- Helpers para serialização/deserialização JSON
- Facilita testes e possibilita múltiplas implementações

**Cached Repository** (`internal/repository/cached_department_repository.go`)

- Decorator pattern sobre o repository base
- Cache automático em operações de leitura
- Invalidação inteligente em operações de escrita
- Fallback para repository base em caso de falha do cache

#### Estratégia de Cache

**Chaves de Cache:**

- `department:{id}` - Dados básicos de um departamento
- `department:hierarchy:{id}` - Árvore completa de departamentos
- `department:subdepts:{id}` - Sub-departamentos diretos
- `department:subdept_ids:{id}` - IDs de todos os sub-departamentos

**Invalidação de Cache:**

- Em operações de criação: invalida cache do departamento pai
- Em operações de atualização: invalida cache do departamento e do pai
- Em operações de deleção: invalida cache do departamento e do pai
- Usa pattern matching para garantir limpeza completa

**TTL (Time To Live):**

- Configurável via `CACHE_TTL` (padrão: 300 segundos)
- Previne dados desatualizados
- Reduz uso de memória

### 2. CTE Recursivo para Consultas Hierárquicas

#### FindByIDWithHierarchy

Recupera um departamento e toda sua sub-árvore usando CTE recursivo.

**Query SQL:**

```sql
WITH RECURSIVE department_tree AS (
    -- Caso base: departamento raiz
    SELECT
        d.id,
        d.name,
        d.manager_id,
        e.name as manager_name,
        d.parent_department_id,
        0 as level
    FROM departments d
    LEFT JOIN employees e ON d.manager_id = e.id
    WHERE d.id = ? AND d.deleted_at IS NULL

    UNION ALL

    -- Caso recursivo: sub-departamentos
    SELECT
        d.id,
        d.name,
        d.manager_id,
        e.name as manager_name,
        d.parent_department_id,
        dt.level + 1
    FROM departments d
    LEFT JOIN employees e ON d.manager_id = e.id
    INNER JOIN department_tree dt ON d.parent_department_id = dt.id
    WHERE d.deleted_at IS NULL
)
SELECT * FROM department_tree ORDER BY level, name
```

**Vantagens:**

- Uma única query ao banco de dados
- Performance O(n) em vez de O(n²)
- Elimina problema de N+1 queries
- Retorna dados ordenados por nível hierárquico

#### FindAllSubDepartmentIDs

Recupera todos os IDs de sub-departamentos usando CTE recursivo.

**Query SQL:**

```sql
WITH RECURSIVE subdepartments AS (
    -- Caso base: departamento inicial
    SELECT id FROM departments
    WHERE id = ? AND deleted_at IS NULL

    UNION ALL

    -- Caso recursivo: descendentes
    SELECT d.id
    FROM departments d
    INNER JOIN subdepartments sd ON d.parent_department_id = sd.id
    WHERE d.deleted_at IS NULL
)
SELECT id FROM subdepartments
```

**Vantagens:**

- Substitui BFS (Breadth-First Search) em memória
- Evita múltiplas queries recursivas
- Performance superior em árvores profundas
- Aproveitamento de índices do PostgreSQL

## Benefícios de Performance

### Antes vs Depois

#### FindByIDWithHierarchy

**Antes (Preload aninhado):**

- Limitado a 3-4 níveis de profundidade
- 1 + N + N² + N³ queries (problema N+1)
- Performance degradada com árvores profundas
- Timeout em hierarquias grandes

**Depois (CTE Recursivo + Cache):**

- Profundidade ilimitada
- 1 query única ao banco
- Cache reduz latência para ~1ms (cache hit)
- Escalável para hierarquias complexas

#### Exemplo Numérico

Para uma árvore com 100 departamentos:

**Sem otimização:**

- Queries: ~100-400 (N+1 problem)
- Tempo: ~500-2000ms
- Cache: não utilizado

**Com CTE + Cache:**

- Queries: 1 (primeira vez), 0 (cache hit)
- Tempo: ~50-100ms (miss), ~1-5ms (hit)
- Cache hit rate: ~80-95% em produção

### Cache Hit Rates Esperados

- **GetById**: 85-95% (dados consultados frequentemente)
- **GetHierarchy**: 70-85% (menos volátil)
- **SubDepartments**: 75-90% (estrutura estável)

## Configuração

### Variáveis de Ambiente

```bash
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Cache
CACHE_TTL=300
```

### Docker Compose

O serviço Redis já está configurado no `docker-compose.yml`:

```yaml
redis:
  image: redis:7-alpine
  container_name: growth_redis
  ports:
    - "6379:6379"
  volumes:
    - redis_data:/data
  command: redis-server --appendonly yes
  healthcheck:
    test: ["CMD", "redis-cli", "ping"]
```

## Monitoramento

### Logs

O sistema registra:

- Conexão/desconexão do Redis
- Falhas de cache (não bloqueantes)
- Invalidações de cache
- Fallback para repository sem cache

### Métricas Importantes

1. **Cache Hit Rate**: % de requisições servidas do cache
2. **Query Time**: tempo de execução das queries CTE
3. **Cache Invalidation Rate**: frequência de invalidações
4. **Redis Memory Usage**: uso de memória do Redis

### Como Monitorar

**Redis CLI:**

```bash
redis-cli INFO stats
redis-cli KEYS "department:*"
redis-cli TTL "department:{id}"
```

**Application Logs:**

```bash
docker-compose logs -f app | grep -i "cache\|redis"
```

## Testes

### Testar Cache

```bash
# 1. Buscar departamento (deve fazer query)
curl http://localhost:8080/api/departments/{id}

# 2. Buscar novamente (deve vir do cache)
curl http://localhost:8080/api/departments/{id}

# 3. Verificar no Redis
redis-cli GET "department:{id}"
```

### Testar Hierarquia

```bash
# Buscar com hierarquia completa
curl http://localhost:8080/api/departments/{id}/hierarchy

# Verificar logs para ver query CTE
docker-compose logs app | grep "department_tree"
```

### Testar Invalidação

```bash
# 1. Buscar departamento
curl http://localhost:8080/api/departments/{id}

# 2. Atualizar departamento
curl -X PUT http://localhost:8080/api/departments/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}'

# 3. Verificar que cache foi invalidado
redis-cli GET "department:{id}"  # Deve retornar nil

# 4. Buscar novamente (deve fazer nova query)
curl http://localhost:8080/api/departments/{id}
```

## Troubleshooting

### Redis não conecta

**Sintoma:** Logs mostram "Failed to connect to Redis"

**Solução:**

- Verificar se Redis está rodando: `docker-compose ps redis`
- Verificar logs do Redis: `docker-compose logs redis`
- Aplicação continua funcionando sem cache

### Cache não invalida

**Sintoma:** Dados desatualizados retornados

**Solução:**

- Verificar TTL: pode estar expirando naturalmente
- Limpar manualmente: `redis-cli FLUSHDB`
- Verificar logs de invalidação

### Performance não melhorou

**Sintoma:** Queries ainda lentas

**Solução:**

1. Verificar cache hit rate no Redis
2. Confirmar que CTE está sendo usado (logs)
3. Verificar índices no PostgreSQL
4. Aumentar TTL se dados mudam pouco

## Próximos Passos

### Possíveis Melhorias

1. **Warming de Cache**: Popular cache em background
2. **Métricas Detalhadas**: Integrar com Prometheus/Grafana
3. **Cache Distribuído**: Redis Cluster para alta disponibilidade
4. **Compressão**: Comprimir valores grandes no cache
5. **Prefetch Inteligente**: Carregar hierarquias relacionadas
6. **Rate Limiting**: Proteger queries pesadas

### Performance Adicional

1. **Índices Compostos**: Otimizar queries frequentes
2. **Materialized Views**: Para agregações pesadas
3. **Particionamento**: Se volume de dados crescer muito
4. **Read Replicas**: Separar leitura de escrita

## Referências

- [PostgreSQL Recursive Queries](https://www.postgresql.org/docs/current/queries-with.html)
- [Redis Best Practices](https://redis.io/docs/manual/patterns/)
- [Cache Invalidation Patterns](https://docs.microsoft.com/en-us/azure/architecture/patterns/cache-aside)
