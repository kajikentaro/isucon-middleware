import { MAX_ROW_LENGTH } from "@/constants";
import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { useAppSelector } from "@/store";
import { selectRecordedTransactionUlids } from "@/store/recorded-transaction";
import Link from "next/link";

export default function Pagination() {
  const { fetchTransactions } = useFetchTransactions();
  const currentPageNum = useCurrentPageNum();
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  return (
    <div className="flex flex-row">
      {currentPageNum !== 1 && (
        <Link
          href={{
            query: { page: currentPageNum - 1 },
          }}
          onClick={() => fetchTransactions(currentPageNum + 1)}
          className="border-blue-500 border-2 text-blue-500 font-bold m-3 py-2 px-3 rounded-lg"
          prefetch={false}
        >
          Previous
        </Link>
      )}
      {recordedTransactionUlids.length === MAX_ROW_LENGTH && (
        <Link
          href={{
            query: { page: currentPageNum + 1 },
          }}
          onClick={() => fetchTransactions(currentPageNum + 1)}
          className="border-blue-500 border-2 text-blue-500 font-bold m-3 py-2 px-3 rounded-lg"
          prefetch={false}
        >
          Next
        </Link>
      )}
    </div>
  );
}
