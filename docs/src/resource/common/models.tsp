import "@typespec/http";

namespace GrowthAPI;

model ErrorResponse {
  error: string;
  message?: string;
}

model SuccessResponse {
  message: string;
}

model ListRequest<T> {
  filters?: T;
  page?: int32 = 1;
  limit?: int32 = 10;
}

model ListResponse<T> {
  data: T[];
  total: int64;
  page: int32;
  limit: int32;
  total_pages: int32;
}

