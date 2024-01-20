import { useAppDispatch, useAppSelector } from "@/store";
import {
  ExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { setRecordedTransactionList } from "@/store/recorded-transaction";
import {
  selectIsFetchingTransactions,
  setIsFetchingTransactions,
} from "@/store/ui/is-fetching-transactions";
import { setSelectedUlids } from "@/store/ui/selectedRowIdx";
import { RecordedTransaction } from "@/types";
import { getFetchListUrl } from "@/utils/get-url";

const MAX_ROW_LENGTH = 100;

export function useFetchTransactions() {
  const dispatch = useAppDispatch();
  const isFetchingTransactions = useAppSelector(selectIsFetchingTransactions);

  const fetchTransactions = async (page: number) => {
    const response = await fetch(
      getFetchListUrl((page - 1) * MAX_ROW_LENGTH, MAX_ROW_LENGTH)
    );
    const json: RecordedTransaction[] = await response.json();
    dispatch(setRecordedTransactionList(json));

    const progressMap: ExecutionProgressMap = {};
    for (const progress of json) {
      progressMap[progress.Ulid] = "init";
    }
    dispatch(setExecutionProgressAll(progressMap));

    dispatch(setSelectedUlids(Object.keys(progressMap)));
    dispatch(setIsFetchingTransactions(false));
  };

  return { fetchTransactions, isFetchingTransactions };
}
