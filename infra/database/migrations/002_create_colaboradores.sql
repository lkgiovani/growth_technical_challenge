CREATE TABLE IF NOT EXISTS colaboradores (
    id UUID PRIMARY KEY,
    nome TEXT NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    rg VARCHAR(20),
    departamento_id UUID NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_colaboradores_cpf ON colaboradores(cpf) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_colaboradores_rg ON colaboradores(rg) WHERE deleted_at IS NULL AND rg IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_colaboradores_deleted_at ON colaboradores(deleted_at);

