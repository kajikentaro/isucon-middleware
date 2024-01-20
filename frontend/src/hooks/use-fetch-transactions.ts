import {
  ExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { useAppDispatch, useAppSelector } from "@/store/main";
import { setRecordedTransactionList } from "@/store/recorded-transaction";
import {
  selectIsFetchingTransactions,
  setIsFetchingTransactions,
} from "@/store/ui";
import { RecordedTransaction } from "@/types";
import { getFetchListUrl } from "@/utils/get-url";
import { useState } from "react";

const MAX_ROW_LENGTH = 100;

export function useFetchTransactions() {
  const dispatch = useAppDispatch();
  const isFetchingTransactions = useAppSelector(selectIsFetchingTransactions);

  const [selected, setSelected] = useState<boolean[]>([]);

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

    setSelected(Array(json.length).fill(true));
    dispatch(setIsFetchingTransactions(false));
  };

  return { fetchTransactions, isFetchingTransactions };
}
