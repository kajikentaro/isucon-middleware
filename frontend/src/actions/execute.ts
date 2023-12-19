import { setExecutionProgress } from "@/store/execution-progress";
import { setExecutionResponse } from "@/store/execution-response";
import { ExecutionResponse } from "@/types";
import { getReproduceUrl } from "@/utils/get-url";
import { AppDispatch, GetState } from "@/store/main";

export function execute(ulid: string) {
  return async (dispatch: AppDispatch, getState: GetState) => {
    try {
      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "waitingResponse",
        })
      );

      const res = await fetch(getReproduceUrl(ulid));
      const json = (await res.json()) as ExecutionResponse;

      dispatch(
        setExecutionResponse({
          ulid,
          executeResponse: json,
        })
      );

      if (!json.IsSameStatusCode) {
        dispatch(
          setExecutionProgress({
            ulid,
            executionProgress: "statusCodeNotSame",
          })
        );
        return;
      }
      if (!json.IsSameResBody) {
        dispatch(
          setExecutionProgress({
            ulid,
            executionProgress: "bodyNotSame",
          })
        );
        return;
      }
      if (!json.IsSameResHeader) {
        dispatch(
          setExecutionProgress({
            ulid,
            executionProgress: "headerNotSame",
          })
        );
        return;
      }
      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "success",
        })
      );
    } catch (e) {
      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "fail",
        })
      );
    }
  };
}
