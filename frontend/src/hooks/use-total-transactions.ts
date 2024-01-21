import { TotalTransactions } from "@/types";
import { getTotalTransactionsURL } from "@/utils/get-url";
import { useEffect, useState } from "react";

export function useTotalTransactions() {
  const [count, setCount] = useState(0);

  useEffect(() => {
    (async () => {
      const res = await fetch(getTotalTransactionsURL());
      const json = (await res.json()) as TotalTransactions;
      setCount(json.Count);
    })();
  }, []);

  return count;
}
