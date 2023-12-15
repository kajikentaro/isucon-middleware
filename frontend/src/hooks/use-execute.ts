import {
    selectExecutionProgress,
    setExecutionProgress,
} from "@/store/execution-progress";
import { setExecutionResponse } from "@/store/execution-response";
import { useAppDispatch, useAppSelector } from "@/store/main";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { ExecutionResponse } from "@/types";
import { getReproduceUrl } from "@/utils/get-url";

export function useOnExecute(ulid: string) {
  const dispatch = useAppDispatch();
  const target = useAppSelector(selectRecordedTransaction(ulid));
  if (!target) {
    throw new Error(`ulid (${ulid}) is not found`);
  }

  return () => {
    dispatch(async (dispatch, getState) => {
      const progress = selectExecutionProgress(ulid)(getState());
      if (progress === "loading") return;
      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "loading",
        })
      );

      try {
        const res = await fetch(getReproduceUrl(ulid));
        const json = (await res.json()) as ExecutionResponse;

        dispatch(
          setExecutionResponse({
            ulid,
            executeResponse: json,
          })
        );
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
    });
  };
}
