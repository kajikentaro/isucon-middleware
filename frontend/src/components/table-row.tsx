"use client";
import { useOnExecute } from "@/hooks/use-execute";
import { TagBinary } from "@/parts/tag-binary";
import { showComparisonPopup } from "@/store/comparison-popup";
import { useAppDispatch, useAppSelector } from "@/store/main";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { MouseEvent } from "react";
import ProgressIcon from "./progress-icon";

interface Props {
  ulid: string;
  handleCheckboxClick: (event: MouseEvent) => void;
  isSelected: boolean;
}

export default function TableRow(props: Props) {
  const { handleCheckboxClick, isSelected, ulid } = props;

  const item = useAppSelector(selectRecordedTransaction(ulid));
  const onExecute = useOnExecute(item.Ulid);
  const dispatch = useAppDispatch();

  const onClickRow = () => {
    dispatch(showComparisonPopup(ulid));
  };

  return (
    <tr className="border-b hover:bg-gray-100" onClick={onClickRow}>
      <td
        className="px-4 py-2 whitespace-nowrap"
        onClick={(e) => handleCheckboxClick(e)}
      >
        <div
          className={`w-3 h-3 border  rounded m-auto block ${
            isSelected ? "bg-blue-500 border-blue-500" : " border-gray-500"
          }`}
        />
      </td>
      <td className="px-4 py-2 whitespace-nowrap">{item.Method}</td>
      <td className="px-4 py-2 whitespace-nowrap">{item.Url}</td>
      <td className="px-4 whitespace-nowrap overflow-hidden max-w-0">
        {item.IsReqText ? (
          <code className="bg-gray-700 text-white p-1 text-xs">
            {item.ReqBody}
          </code>
        ) : (
          <TagBinary ulid={ulid} type="req-body" />
        )}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        {item.StatusCode.toString()}
      </td>
      <td className="px-4 py-2 whitespace-nowrap overflow-hidden max-w-0">
        {item.IsResText ? (
          <code className="bg-gray-700 text-white p-1 text-xs">
            {item.ResBody}
          </code>
        ) : (
          <TagBinary ulid={ulid} type="res-body" />
        )}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        <ProgressIcon ulid={item.Ulid} />
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded-full flex items-center m-auto"
          onClick={(e) => {
            e.stopPropagation();
            onExecute();
          }}
        >
          <svg
            className="h-4 w-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M9 5l8 8-8 8"
            />
          </svg>
        </button>
      </td>
    </tr>
  );
}
