import { execute } from "@/actions/execute";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectExecutionProgress } from "@/store/execution-progress";
import { selectRecordedTransaction } from "@/store/recorded-transaction";

export function useExecute(ulid: string) {
  const dispatch = useAppDispatch();
  const target = useAppSelector(selectRecordedTransaction(ulid));
  if (!target) {
    throw new Error(`ulid (${ulid}) is not found`);
  }

  return () => {
    dispatch(async (dispatch, getState) => {
      const progress = selectExecutionProgress(ulid)(getState());
      if (progress === "waitingQueue" || progress === "waitingResponse") return;

      dispatch(execute(ulid));
    });
  };
}
