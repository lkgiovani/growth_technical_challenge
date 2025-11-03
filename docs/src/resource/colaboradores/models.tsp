import "@typespec/http";
import "../../resource/departamentos/models.tsp";

using TypeSpec.Http;

namespace GrowthAPI;

@doc("Modelo de colaborador")
model Colaborador {
  @doc("ID do colaborador (UUID)")
  @format("uuid")
  id?: string;

  @doc("Nome do colaborador")
  nome: string;

  @doc("CPF brasileiro (11 dígitos)")
  @pattern("^[0-9]{11}$")
  @example("12345678901")
  cpf: string;

  @doc("RG brasileiro (opcional)")
  @minLength(1)
  @maxLength(20)
  rg?: string;

  @doc("Referência ao ID do departamento")
  @format("uuid")
  departamento_id: string;

  @doc("Departamento associado")
  departamento?: Departamento;

  @doc("Timestamp de criação")
  created_at?: utcDateTime;

  @doc("Timestamp da última atualização")
  updated_at?: utcDateTime;
}

@doc("Dados mínimos para criar um departamento para o colaborador")
model DepartamentoForColaborador {
  @doc("Nome do departamento")
  @example("Tecnologia")
  nome: string;

  @doc("ID do gerente do departamento")
  @format("uuid")
  @example("22222222-2222-2222-2222-222222222222")
  gerente_id: string;

  @doc("ID do departamento superior (opcional)")
  @format("uuid")
  @example("11111111-1111-1111-1111-111111111111")
  departamento_superior_id?: string;
}

@doc("""
Requisição para criar colaborador.

**REGRA OBRIGATÓRIA:** Forneça exatamente UM dos dois campos (departamento_id OU departamento).
""")
model ColaboradorCreateRequest {
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

  @doc("ID do departamento existente - Use APENAS esta opção OU 'departamento', nunca ambos")
  @format("uuid")
  @example("22222222-2222-2222-2222-222222222222")
  departamento_id?: string;

  @doc("Dados do novo departamento - Use APENAS esta opção OU 'departamento_id', nunca ambos")
  departamento?: DepartamentoForColaborador;
}

@doc("Requisição para atualizar colaborador")
model ColaboradorUpdateRequest {
  @doc("Nome do colaborador")
  nome?: string;

  @doc("CPF brasileiro (11 dígitos)")
  @pattern("^[0-9]{11}$")
  cpf?: string;

  @doc("RG brasileiro")
  @minLength(1)
  @maxLength(20)
  rg?: string;

  @doc("Referência ao ID do departamento")
  @format("uuid")
  departamento_id?: string;
}

@doc("Filtros de colaborador para operações de listagem")
model ColaboradorFilters {
  @doc("Filtrar por nome")
  nome?: string;

  @doc("Filtrar por CPF")
  cpf?: string;

  @doc("Filtrar por RG")
  rg?: string;

  @doc("Filtrar por ID do departamento")
  @format("uuid")
  departamento_id?: string;
}

@doc("Resposta de colaborador com informações do gerente")
model ColaboradorWithGerente {
  ...Colaborador;

  @doc("Nome do gerente do departamento")
  gerente_departamento_nome?: string;
}

