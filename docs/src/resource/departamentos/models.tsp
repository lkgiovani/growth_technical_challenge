import "@typespec/http";

using TypeSpec.Http;

namespace GrowthAPI;

@doc("Modelo de departamento")
model Departamento {
  @doc("ID do departamento (UUID)")
  @format("uuid")
  id?: string;

  @doc("Nome do departamento")
  nome: string;

  @doc("Referência ao ID do gerente")
  @format("uuid")
  gerente_id: string;

  @doc("Informações do gerente")
  gerente?: Colaborador;

  @doc("ID do departamento superior (opcional)")
  @format("uuid")
  departamento_superior_id?: string;

  @doc("Departamento superior")
  departamento_superior?: Departamento;

  @doc("Subdepartamentos")
  subdepartamentos?: Departamento[];

  @doc("Colaboradores neste departamento")
  colaboradores?: Colaborador[];

  @doc("Timestamp de criação")
  created_at?: utcDateTime;

  @doc("Timestamp da última atualização")
  updated_at?: utcDateTime;
}

@doc("Requisição para criar departamento")
model DepartamentoCreateRequest {
  @doc("Nome do departamento")
  nome: string;

  @doc("Referência ao ID do gerente")
  @format("uuid")
  gerente_id: string;

  @doc("ID do departamento superior (opcional)")
  @format("uuid")
  departamento_superior_id?: string;
}

@doc("Requisição para atualizar departamento")
model DepartamentoUpdateRequest {
  @doc("Nome do departamento")
  nome?: string;

  @doc("Referência ao ID do gerente")
  @format("uuid")
  gerente_id?: string;

  @doc("ID do departamento superior")
  @format("uuid")
  departamento_superior_id?: string;
}

@doc("Filtros de departamento para operações de listagem")
model DepartamentoFilters {
  @doc("Filtrar por nome")
  nome?: string;

  @doc("Filtrar por nome do gerente")
  gerente_nome?: string;

  @doc("Filtrar por ID do departamento superior")
  @format("uuid")
  departamento_superior_id?: string;
}

@doc("Resposta de departamento com hierarquia")
model DepartamentoWithHierarchy {
  ...Departamento;

  @doc("Árvore hierárquica completa de subdepartamentos")
  hierarquia?: Departamento[];
}

