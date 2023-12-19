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
import { getFetchAllUrl } from "@/utils/get-url";
import { MouseEvent, useEffect, useState } from "react";

export default function Main() {
  const dispatch = useAppDispatch();
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  const [selected, setSelected] = useState<boolean[]>([]);
  const [lastSelectedIndex, setLastSelectedIndex] = useState(-1);

  const onExecuteChecked = useExecuteChecked();

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

  const fetchData = async () => {
    const response = await fetch(getFetchAllUrl(), {});
    const json: RecordedTransaction[] = await response.json();
    dispatch(setRecordedTransactionList(json));

    const progressMap: ExecutionProgressMap = {};
    for (const progress of json) {
      progressMap[progress.Ulid] = "init";
    }
    dispatch(setExecutionProgressAll(progressMap));

    setSelected(Array(json.length).fill(true));
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="flex flex-col justify-center items-center">
      <div className="flex w-full justify-between px-4 py-3 mb-2">
        <h1 className="text-3xl font-bold">Isucon Middleware</h1>
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
    </div>
  );
}
