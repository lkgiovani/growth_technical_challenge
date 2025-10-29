import "@typespec/http";
import "@typespec/rest";
import "../../resource/common/models.tsp";
import "../../resource/colaboradores/models.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

@tag("Gerentes")
@route("/gerentes")
interface GerenteOperations {
  @doc("Obter todos os colaboradores dos departamentos subordinados ao gerente recursivamente")
  @get
  @route("/{id}/colaboradores")
  @summary("Obter colaboradores do gerente")
  getGerenteColaboradores(@path id: string): {
    @statusCode statusCode: 200;
    @body body: Colaborador[];
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };
}

