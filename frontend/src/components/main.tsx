"use client";
import TableRow from "@/components/table-row";
import { useExecuteChecked } from "@/hooks/use-execute-checked";
import {
  ExecutionProgressMap,
  setExecutionProgressAll,
} from "@/store/execution-progress";
import { useAppDispatch, useAppSelector } from "@/store/main";
import {
  selectRecordedTransactionUlids,
  setRecordedTransactionList,
} from "@/store/recorded-transaction";
import { RecordedTransaction } from "@/types";
import { getFetchListUrl } from "@/utils/get-url";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { MouseEvent, useEffect, useMemo, useState } from "react";
import StartRecordingButton from "./start-recording-button";

const MAX_ROW_LENGTH = 20;

export default function Main() {
  const dispatch = useAppDispatch();
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  const [selected, setSelected] = useState<boolean[]>([]);
  const [lastSelectedIndex, setLastSelectedIndex] = useState(-1);
  const [isLoading, setIsLoading] = useState(true);

  const onExecuteChecked = useExecuteChecked();

  const searchParams = useSearchParams();
  const currentPageNum = useMemo(() => {
    const page = searchParams.get("page") || "1";

    const pageInt = Number(page);
    if (isNaN(pageInt)) {
      return 1;
    }
    return pageInt;
  }, [searchParams]);

  const isAllSelected = selected.every((s) => s) && selected.length > 0;

  const handleCheckboxClick = (event: MouseEvent, index: number) => {
    if (!selected.length) return;

    const newSelected = [...selected];

    if (event.shiftKey) {
      // Shift-click: select all rows in range
      const nextIsTrue = !selected[index];
      for (
        let i = Math.min(lastSelectedIndex, index);
        i <= Math.max(lastSelectedIndex, index);
        i++
      ) {
        newSelected[i] = nextIsTrue;
      }
    } else {
      // Select only clicked row
      newSelected[index] = !selected[index];
    }

    setSelected(newSelected);
    setLastSelectedIndex(index);

    event.preventDefault();
    event.stopPropagation();
  };

  const fetchData = async (page: number) => {
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
    setIsLoading(false);
  };

  useEffect(() => {
    fetchData(currentPageNum);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  return (
    <div className="flex flex-col justify-center items-center">
      <div className="flex w-full justify-between px-4 py-3 mb-2">
        <h1 className="text-3xl font-bold">Isucon Middleware</h1>
        <div className="flex gap-x-5">
          <StartRecordingButton />
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded-full flex items-center"
            onClick={(e) => {
              const checkedUlids = recordedTransactionUlids.filter(
                (_, idx) => selected[idx]
              );
              onExecuteChecked(checkedUlids);
              e.stopPropagation();
            }}
          >
            Execute Checked
          </button>
        </div>
      </div>
      <table className="table-auto border-collapse w-full">
        <thead>
          <tr className="border-b bg-gray-100 text-gray-600">
            <th
              className="px-4 py-2 whitespace-nowrap"
              onClick={() => {
                setSelected(Array(selected.length).fill(!isAllSelected));
              }}
            >
              <div
                className={`w-4 h-4 border border-gray-500 rounded m-auto block ${
                  isAllSelected ? "bg-blue-500" : "bg-white"
                }`}
              />
            </th>
            <th className="px-4 py-2 whitespace-nowrap">Method</th>
            <th className="px-4 py-2 whitespace-nowrap">URL</th>
            <th className="px-4 py-2 w-1/2">ReqBody</th>
            <th className="px-4 py-2">Status Code</th>
            <th className="px-4 py-2 w-1/2">ResBody</th>
            <th>Execution Result</th>
            <th className="px-4 py-2 w-1/2">Execute</th>
          </tr>
        </thead>
        <tbody>
          {recordedTransactionUlids.map((ulid, index) => (
            <TableRow
              isSelected={selected[index]}
              handleCheckboxClick={(e: MouseEvent) =>
                handleCheckboxClick(e, index)
              }
              ulid={ulid}
              key={ulid}
            />
          ))}
        </tbody>
      </table>
      {isLoading && <p>now loading</p>}
      {!isLoading && recordedTransactionUlids.length === 0 && (
        <p>result was not found</p>
      )}

      <div className="flex flex-row">
        {currentPageNum !== 1 && (
          <Link
            href={{
              query: { page: currentPageNum - 1 },
            }}
            onClick={() => fetchData(currentPageNum + 1)}
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
            onClick={() => fetchData(currentPageNum + 1)}
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
