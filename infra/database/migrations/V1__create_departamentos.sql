CREATE TABLE IF NOT EXISTS departamentos (
    id UUID PRIMARY KEY,
    nome TEXT NOT NULL,
    gerente_id UUID NOT NULL,
    departamento_superior_id UUID,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_departamentos_deleted_at ON departamentos(deleted_at);

