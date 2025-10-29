import "@typespec/http";
import "@typespec/rest";
import "../../resource/common/models.tsp";
import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

@tag("Colaboradores")
@route("/colaboradores")
interface ColaboradorOperations {
  @doc("Criar um novo colaborador com validação de CPF e RG")
  @post
  @summary("Criar colaborador")
  createColaborador(@body colaborador: ColaboradorCreateRequest): {
    @statusCode statusCode: 201;
    @body body: Colaborador;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  } | {
    @statusCode statusCode: 409;
    @body body: ErrorResponse;
  } | {
    @statusCode statusCode: 422;
    @body body: ErrorResponse;
  };

  @doc("Obter colaborador por ID com informações do gerente do departamento")
  @get
  @summary("Obter colaborador")
  getColaborador(@path id: string): {
    @statusCode statusCode: 200;
    @body body: ColaboradorWithGerente;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @doc("Atualizar informações do colaborador")
  @put
  @summary("Atualizar colaborador")
  updateColaborador(@path id: string, @body colaborador: ColaboradorUpdateRequest): {
    @statusCode statusCode: 200;
    @body body: Colaborador;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  } | {
    @statusCode statusCode: 409;
    @body body: ErrorResponse;
  } | {
    @statusCode statusCode: 422;
    @body body: ErrorResponse;
  };

  @doc("Remover um colaborador (soft delete)")
  @delete
  @summary("Deletar colaborador")
  deleteColaborador(@path id: string): {
    @statusCode statusCode: 200;
    @body body: SuccessResponse;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @doc("Listar colaboradores com filtros e paginação")
  @post
  @route("/listar")
  @summary("Listar colaboradores")
  listColaboradores(@body request: ListRequest<ColaboradorFilters>): {
    @statusCode statusCode: 200;
    @body body: ListResponse<Colaborador>;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

