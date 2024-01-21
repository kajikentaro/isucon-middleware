"use client";
import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useFetchRecordedTransactions } from "@/hooks/use-fetch-recorded-transactions";
import { useFetchTotalTransactions } from "@/hooks/use-total-transactions";
import Link from "next/link";
import { useEffect } from "react";
import ExecuteCheckedButton from "./execute-checked-button";
import Pagination from "./pagination";
import RemoveAllButton from "./remove-all-button";
import RemoveSelectedButton from "./remove-selected-button";
import StartRecordingButton from "./start-recording-button";
import Table from "./table";

export default function Main() {
  const { fetchTransactions } = useFetchRecordedTransactions();
  const { fetchTotalTransactions } = useFetchTotalTransactions();
  const currentPageNum = useCurrentPageNum();

  useEffect(() => {
    fetchTransactions(currentPageNum);
    fetchTotalTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentPageNum]);

  return (
    <div className="flex flex-col justify-center items-center">
      <div className="flex w-full justify-between px-4 py-3 mb-2">
        <Link href="/">
          <h1 className="text-3xl font-bold">Isucon Middleware</h1>
        </Link>
        <div className="flex gap-x-5">
          <RemoveSelectedButton />
          <RemoveAllButton />
          <StartRecordingButton />
          <ExecuteCheckedButton />
        </div>
      </div>
      <Table />

      <Pagination />
    </div>
  );
}
