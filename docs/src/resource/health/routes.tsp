import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

@tag("Health")
@route("/health")
interface HealthOperations {
  @doc("Check API health status")
  @get
  @summary("Health check")
  healthCheck(): {
    @statusCode statusCode: 200;
    @body body: {
      status: "ok";
      timestamp: utcDateTime;
    };
  };
}

