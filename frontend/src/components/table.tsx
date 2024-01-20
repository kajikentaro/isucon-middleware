"use client";
import TableRow from "@/components/table-row";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectRecordedTransactionUlids } from "@/store/recorded-transaction";
import {
  selectSelectedUlids,
  setSelectedUlids,
} from "@/store/ui/selectedRowIdx";
import { MouseEvent, useState } from "react";
import TableHeader from "./table-header";

export default function Table() {
  const { isFetchingTransactions } = useFetchTransactions();
  const selectedUlids = useAppSelector(selectSelectedUlids);
  const dispatch = useAppDispatch();
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  // this is used for selecting the range where user click transactions with Shift key
  const [lastSelectedIndex, setLastSelectedIndex] = useState(-1);

  const onCheckboxClick = (event: MouseEvent, index: number, ulid: string) => {
    const isSelectedAsTrue = !selectedUlids.includes(ulid);
    const set = new Set(selectedUlids);

    if (event.shiftKey) {
      // Shift-click: select all rows in range
      for (
        let i = Math.min(lastSelectedIndex, index);
        i <= Math.max(lastSelectedIndex, index);
        i++
      ) {
        if (isSelectedAsTrue) {
          set.add(recordedTransactionUlids[i]);
        } else {
          set.delete(recordedTransactionUlids[i]);
        }
      }
    } else {
      // Select only clicked row
      if (isSelectedAsTrue) {
        set.add(ulid);
      } else {
        set.delete(ulid);
      }
    }

    dispatch(setSelectedUlids(Array.from(set)));
    setLastSelectedIndex(index);

    event.preventDefault();
    event.stopPropagation();
  };

  return (
    <>
      <table className="table-auto border-collapse w-full">
        <TableHeader />
        <tbody>
          {recordedTransactionUlids.map((ulid, index) => (
            <TableRow
              isSelected={selectedUlids.includes(ulid)}
              onCheckboxClick={(e: MouseEvent) =>
                onCheckboxClick(e, index, ulid)
              }
              ulid={ulid}
              key={ulid}
            />
          ))}
        </tbody>
      </table>

      {isFetchingTransactions && <p>now loading</p>}
      {!isFetchingTransactions && recordedTransactionUlids.length === 0 && (
        <p>result was not found</p>
      )}
    </>
  );
}
