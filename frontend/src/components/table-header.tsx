"use client";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectRecordedTransactionUlids } from "@/store/recorded-transaction";
import {
  selectSelectedUlids,
  setSelectedUlids,
} from "@/store/ui/selected-ulids";

export default function TableHeader() {
  const dispatch = useAppDispatch();
  const selectedUlids = useAppSelector(selectSelectedUlids);
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  const isAllSelected =
    selectedUlids.length === recordedTransactionUlids.length;

  const handleClick = () => {
    if (isAllSelected) {
      dispatch(setSelectedUlids([]));
    } else {
      dispatch(setSelectedUlids(recordedTransactionUlids));
    }
  };

  return (
    <thead>
      <tr className="border-b bg-gray-100 text-gray-600">
        <th className="px-4 py-2 whitespace-nowrap" onClick={handleClick}>
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
  );
}
