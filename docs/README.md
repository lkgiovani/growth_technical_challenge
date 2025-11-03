# API Documentation with TypeSpec

Este diretório contém a especificação TypeSpec da API que gera a documentação OpenAPI 3.1.

## O que é TypeSpec?

TypeSpec é uma linguagem para descrever APIs e gerar especificações OpenAPI, código de cliente, documentação e outros assets.

## 🚀 Como Rodar

### 1. Instalar dependências

```bash
cd docs
npm install
```

### 2. Compilar TypeSpec para OpenAPI

```bash
npm run compile
```

Isso irá gerar a especificação OpenAPI no diretório `tsp-output/@typespec/openapi3/openapi.yaml`.

### 3. Copiar para a aplicação

O arquivo gerado precisa ser copiado para onde a aplicação Go pode lê-lo:

```bash
cp tsp-output/@typespec/openapi3/openapi.yaml ../internal/delivery/http/resources/openapi.yaml
```

### 4. Reiniciar a aplicação

```bash
cd ..
docker-compose restart app
```

### 5. Acessar a documentação

Abra no navegador:

- http://localhost:8080/
- http://localhost:8080/docs/swagger
- http://localhost:8080/docs/redoc
- http://localhost:8080/docs/scalar

## 🔄 Modo Watch (Desenvolvimento)

Para recompilar automaticamente quando os arquivos mudarem:

```bash
npm run watch
```

## 📁 Estrutura do Projeto

```
docs/
├── src/
│   ├── main.tsp                    # Entry point principal
│   └── resource/
│       ├── colaboradores/
│       │   ├── models.tsp          # Modelos de colaboradores
│       │   └── routes.tsp          # Rotas de colaboradores
│       └── departamentos/
│           ├── models.tsp          # Modelos de departamentos
│           └── routes.tsp          # Rotas de departamentos
├── tspconfig.yaml                  # Configuração do compilador
├── package.json                    # Dependências Node.js
└── tsp-output/                     # OpenAPI gerado (gitignored)
```

## 📝 Editar a Documentação

1. Edite os arquivos `.tsp` em `src/resource/`
2. Execute `npm run compile` para regenerar o OpenAPI
3. Copie o arquivo gerado para `../internal/delivery/http/resources/openapi.yaml`
4. Reinicie a aplicação

## ✨ Benefícios do TypeSpec

- ✅ Type Safety: Detecta erros em tempo de compilação
- ✅ Reutilização: Compartilhe modelos entre operações
- ✅ Manutenibilidade: Fonte única da verdade para contratos da API
- ✅ Consistência: Força padrões consistentes na API
- ✅ Geração de código: Pode gerar SDKs de cliente automaticamente

## 📖 Referência

- [TypeSpec Documentation](https://typespec.io/)
- [OpenAPI Specification](https://www.openapis.org/)
