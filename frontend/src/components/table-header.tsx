"use client";
import { usePageParams } from "@/hooks/use-page-params";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectRecordedTransactionUlids } from "@/store/recorded-transaction";
import { showFilterPopup } from "@/store/ui/filter-popup";
import {
  selectSelectedUlids,
  setSelectedUlids,
} from "@/store/ui/selected-ulids";
import { AiOutlineFilter } from "react-icons/ai";

export default function TableHeader() {
  const dispatch = useAppDispatch();
  const selectedUlids = useAppSelector(selectSelectedUlids);
  const recordedTransactionUlids = useAppSelector(
    selectRecordedTransactionUlids
  );

  const { query } = usePageParams();
  const isFilterEnabled = !!query;
  const isAllSelected =
    selectedUlids.length === recordedTransactionUlids.length;

  const onClickRow = () => {
    if (isAllSelected) {
      dispatch(setSelectedUlids([]));
    } else {
      dispatch(setSelectedUlids(recordedTransactionUlids));
    }
  };

  const onClickFilter = () => {
    dispatch(showFilterPopup());
  };

  const classNameTh = "px-2 w-0";
  const classNameThUrl = "px-4 py-2 w-5/12";
  const classNameThBody = "px-4 py-2 w-3/12";

  return (
    <thead>
      <tr className="border-b bg-gray-100 text-gray-600">
        <th className={classNameTh} onClick={onClickRow}>
          <div
            className={`w-4 h-4 border border-gray-500 rounded m-auto block ${
              isAllSelected ? "bg-blue-500" : "bg-white"
            }`}
          />
        </th>
        <th className={classNameTh}>Method</th>
        <th className={classNameThUrl}>
          <div className="flex justify-center items-center gap-2">
            URL
            <button
              className="rounded"
              aria-label="filter"
              onClick={onClickFilter}
            >
              <AiOutlineFilter
                size="20"
                className={isFilterEnabled ? "text-red-600" : "text-gray-600"}
              />
            </button>
          </div>
        </th>
        <th className={classNameThBody}>ReqBody</th>
        <th className={classNameTh}>Status Code</th>
        <th className={classNameThBody}>ResBody</th>
        <th className={classNameTh}>Execution Result</th>
        <th className={classNameTh}>Execute</th>
      </tr>
    </thead>
  );
}
