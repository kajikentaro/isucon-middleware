"use client";
import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { useEffect } from "react";
import ExecuteCheckedButton from "./execute-checked-button";
import Pagination from "./pagination";
import RemoveAllButton from "./remove-all-button";
import RemoveSelectedButton from "./remove-selected-button";
import StartRecordingButton from "./start-recording-button";
import Table from "./table";

export default function Main() {
  const { fetchTransactions } = useFetchTransactions();
  const currentPageNum = useCurrentPageNum();

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
          <ExecuteCheckedButton />
        </div>
      </div>
      <Table />

      <Pagination />
    </div>
  );
}
