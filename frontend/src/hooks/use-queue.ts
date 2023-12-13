import {
  selectExecutionProgressMap,
  setExecutionProgress,
} from "@/store/executionProgressMap";
import { useAppDispatch, useAppSelector } from "@/store/main";
import { selectRecordedTransactions } from "@/store/recordedTransactions";
import { getExecuteUrl } from "@/utils/getUrl";

export function useOnExecute(ulid: string) {
  const dispatch = useAppDispatch();
  const recordedTransactions = useAppSelector(selectRecordedTransactions);
  const target = recordedTransactions.find((v) => v.Ulid === ulid);
  if (!target) {
    throw new Error(`ulid (${ulid}) is not found`);
  }

  return () => {
    dispatch(async (dispatch, getState) => {
      const statusMap = selectExecutionProgressMap(getState());
      if (statusMap[ulid] === "loading") return;
      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "loading",
        })
      );

      const res = await fetch(getExecuteUrl(ulid));
      const json = res.json();

      dispatch(
        setExecutionProgress({
          ulid,
          executionProgress: "success",
        })
      );
    });
  };
}
