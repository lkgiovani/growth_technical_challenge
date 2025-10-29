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

@doc("Requisição para criar colaborador")
model ColaboradorCreateRequest {
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

