import "@typespec/http";
import "@typespec/rest";
import "../../resource/common/models.tsp";
import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

@tag("Departamentos")
@route("/departamentos")
interface DepartamentoOperations {
  @doc("Criar um novo departamento com validação de gerente e prevenção de ciclos")
  @post
  @summary("Criar departamento")
  createDepartamento(@body departamento: DepartamentoCreateRequest): {
    @statusCode statusCode: 201;
    @body body: Departamento;
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

  @doc("Obter departamento por ID com árvore hierárquica completa de subdepartamentos")
  @get
  @summary("Obter departamento com hierarquia")
  getDepartamento(@path id: string): {
    @statusCode statusCode: 200;
    @body body: DepartamentoWithHierarchy;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @doc("Atualizar departamento com prevenção de ciclos na hierarquia")
  @put
  @summary("Atualizar departamento")
  updateDepartamento(@path id: string, @body departamento: DepartamentoUpdateRequest): {
    @statusCode statusCode: 200;
    @body body: Departamento;
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

  @doc("Deletar um departamento")
  @delete
  @summary("Deletar departamento")
  deleteDepartamento(@path id: string): {
    @statusCode statusCode: 200;
    @body body: SuccessResponse;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @doc("Listar departamentos com filtros e paginação")
  @post
  @route("/listar")
  @summary("Listar departamentos")
  listDepartamentos(@body request: ListRequest<DepartamentoFilters>): {
    @statusCode statusCode: 200;
    @body body: ListResponse<Departamento>;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

