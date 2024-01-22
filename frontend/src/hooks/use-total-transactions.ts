import { useAppDispatch, useAppSelector } from "@/store";
import {
  selectTotalTransactions,
  setTotalTransactions,
} from "@/store/total-transactions";
import { TotalTransactions } from "@/types";
import { getTotalTransactionsURL } from "@/utils/get-url";

export function useFetchTotalTransactions() {
  const dispatch = useAppDispatch();
  const totalTransactions = useAppSelector(selectTotalTransactions);

  const fetchTotalTransactions = async () => {
    dispatch(setTotalTransactions(-1));
    const res = await fetch(getTotalTransactionsURL());
    const json = (await res.json()) as TotalTransactions;
    const count = json.count;
    dispatch(setTotalTransactions(count));
  };

  return { fetchTotalTransactions, totalTransactions };
}
