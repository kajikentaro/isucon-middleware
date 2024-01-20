"use client";
import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useExecuteChecked } from "@/hooks/use-execute-checked";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { useAppSelector } from "@/store";
import { selectRecordedTransactionUlids } from "@/store/recorded-transaction";
import Link from "next/link";
import { useEffect } from "react";
import RemoveAllButton from "./remove-all-button";
import RemoveSelectedButton from "./remove-selected-button";
import StartRecordingButton from "./start-recording-button";
import Table from "./table";

const MAX_ROW_LENGTH = 100;

export default function Main() {
  const { fetchTransactions } = useFetchTransactions();
  const executeChecked = useExecuteChecked();
  const currentPageNum = useCurrentPageNum();
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  useEffect(() => {
    fetchTransactions(currentPageNum);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentPageNum]);

  return (
    <div className="flex flex-col justify-center items-center">
      <div className="flex w-full justify-between px-4 py-3 mb-2">
        <h1 className="text-3xl font-bold">Isucon Middleware</h1>
        <div className="flex gap-x-5">
          <RemoveSelectedButton />
          <RemoveAllButton />
          <StartRecordingButton />
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded-full flex items-center"
            onClick={(e) => {
              executeChecked();
              e.stopPropagation();
            }}
          >
            Execute Checked
          </button>
        </div>
      </div>
      <Table />

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
    </div>
  );
}
