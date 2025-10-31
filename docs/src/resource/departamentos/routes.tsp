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
  @doc("""
  Criar um novo departamento com validação de gerente e prevenção de ciclos.
  
  **REGRA IMPORTANTE:** Você DEVE fornecer exatamente UM dos dois campos: gerente_id OU gerente.
  
  **Validações:**
  - Previne ciclos na hierarquia de departamentos
  - Valida CPF e RG únicos quando criar novo gerente
  - Verifica existência do departamento superior se fornecido
  - Verifica existência do gerente quando usar gerente_id
  """)
  @opExample(
    #{
      parameters: #{
        departamento: #{
          nome: "Tecnologia",
          gerente_id: "22222222-2222-2222-2222-222222222222"
        }
      }
    },
    #{
      title: "Opção 1 - Com gerente existente",
      description: "Use quando você já tem um colaborador cadastrado"
    }
  )
  @opExample(
    #{
      parameters: #{
        departamento: #{
          nome: "Recursos Humanos",
          gerente: #{
            nome: "Maria Santos",
            cpf: "98765432109",
            rg: "987654321"
          }
        }
      }
    },
    #{
      title: "Opção 2 - Criando gerente junto",
      description: "Use quando você quer criar o colaborador e o departamento ao mesmo tempo"
    }
  )
  @post
  @summary("Criar departamento")
  createDepartamento(@body departamento: DepartamentoCreateRequest): {
    @statusCode statusCode: 201;
    @body body: Departamento;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
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

