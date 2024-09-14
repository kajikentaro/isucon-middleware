import { MAX_ROW_LENGTH } from "@/constants";
import { getPageParams } from "@/hooks/use-page-params";
import { AppDispatch, GetState } from "@/store";
import {
  ExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { setRecordedTransactionList } from "@/store/recorded-transaction";
import { setTotalTransactions } from "@/store/total-transactions";
import { setIsFetchingTransactions } from "@/store/ui/is-fetching-transactions";
import { setSelectedUlids } from "@/store/ui/selected-ulids";
import { SearchResponse } from "@/types";
import { getSearchUrl } from "@/utils/get-url";

export function fetchTransactions(
  optionalPageNum?: number,
  optionalQuery?: string
) {
  return async (dispatch: AppDispatch, getState: GetState) => {
    const params = getPageParams();
    const page = optionalPageNum || params.page;
    // should keep "" if optionalQuery is ""
    const query = optionalQuery ?? params.query;

    const response = await fetch(
      getSearchUrl((page - 1) * MAX_ROW_LENGTH, MAX_ROW_LENGTH, query)
    );
    const { transactions, totalHit } =
      (await response.json()) as SearchResponse;
    dispatch(setRecordedTransactionList(transactions));
    dispatch(setTotalTransactions(totalHit));

    const progressMap: ExecutionProgressMap = {};
    for (const t of transactions) {
      progressMap[t.ulid] = "init";
    }
    dispatch(setExecutionProgressAll(progressMap));

    dispatch(setSelectedUlids(Object.keys(progressMap)));
    dispatch(setIsFetchingTransactions(false));
  };
}
