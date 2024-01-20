import { execute } from "@/actions/execute";
import { useAppDispatch, useAppSelector } from "@/store";
import {
  selectExecutionProgress,
  selectExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { selectSelectedUlids } from "@/store/ui/selectedRowIdx";

export function useExecuteChecked() {
  const dispatch = useAppDispatch();
  const selectedUlids = useAppSelector(selectSelectedUlids);

  return () => {
    dispatch(async (dispatch, getState) => {
      // update status to "waitingQueue"
      const executionProgress = { ...selectExecutionProgressMap(getState()) };
      for (const ulid of selectedUlids) {
        if (executionProgress[ulid] === "waitingResponse") {
          continue;
        }
        executionProgress[ulid] = "waitingQueue";
      }
      dispatch(setExecutionProgressAll(executionProgress));

      // execute in order
      for (const ulid of selectedUlids) {
        const progress = selectExecutionProgress(ulid)(getState());
        if (progress !== "waitingQueue") continue;

        await dispatch(execute(ulid));
      }
    });
  };
}
