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

@doc("Dados mínimos para criar um colaborador como gerente")
model ColaboradorForGerente {
  @doc("Nome do colaborador")
  @example("João Silva")
  nome: string;

  @doc("CPF brasileiro (11 dígitos)")
  @pattern("^[0-9]{11}$")
  @example("12345678901")
  cpf: string;

  @doc("RG brasileiro (opcional)")
  @minLength(1)
  @maxLength(20)
  @example("123456789")
  rg?: string;
}

@doc("""
Requisição para criar departamento.

**REGRA OBRIGATÓRIA:** Forneça exatamente UM dos dois campos (gerente_id OU gerente).
""")
model DepartamentoCreateRequest {
  @doc("Nome do departamento")
  @example("Tecnologia")
  nome: string;

  @doc("ID do gerente existente - Use APENAS esta opção OU 'gerente', nunca ambos")
  @format("uuid")
  @example("22222222-2222-2222-2222-222222222222")
  gerente_id?: string;

  @doc("Dados do novo gerente - Use APENAS esta opção OU 'gerente_id', nunca ambos")
  gerente?: ColaboradorForGerente;

  @doc("ID do departamento superior (opcional)")
  @format("uuid")
  @example("11111111-1111-1111-1111-111111111111")
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

