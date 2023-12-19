import { execute } from "@/actions/execute";
import {
  selectExecutionProgress,
  selectExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { useAppDispatch } from "@/store/main";

export function useExecuteChecked() {
  const dispatch = useAppDispatch();

  return (ulidList: string[]) => {
    dispatch(async (dispatch, getState) => {
      // update status to "waitingQueue"
      const executionProgress = { ...selectExecutionProgressMap(getState()) };
      for (const ulid of ulidList) {
        if (executionProgress[ulid] === "waitingResponse") {
          continue;
        }
        executionProgress[ulid] = "waitingQueue";
      }
      dispatch(setExecutionProgressAll(executionProgress));

      // execute in order
      for (const ulid of ulidList) {
        const progress = selectExecutionProgress(ulid)(getState());
        if (progress !== "waitingQueue") continue;

        await dispatch(execute(ulid));
      }
    });
  };
}
