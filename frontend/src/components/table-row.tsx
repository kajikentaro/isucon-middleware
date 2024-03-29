"use client";
import { useExecute } from "@/hooks/use-execute";
import { useOpenPopup } from "@/hooks/use-open-popup";
import Code from "@/parts/code";
import { TagBinary } from "@/parts/tag-binary";
import { useAppSelector } from "@/store";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { MouseEvent } from "react";
import ProgressIcon from "../parts/progress-icon";

interface Props {
  ulid: string;
  onCheckboxClick: (event: MouseEvent) => void;
  isSelected: boolean;
}

export default function TableRow(props: Props) {
  const { onCheckboxClick, isSelected, ulid } = props;

  const item = useAppSelector(selectRecordedTransaction(ulid));
  const onExecute = useExecute(item.ulid);
  const openPopup = useOpenPopup(item.ulid);

  return (
    <tr className="border-b hover:bg-gray-100" onClick={openPopup}>
      <td className="px-4 py-2 whitespace-nowrap" onClick={onCheckboxClick}>
        <div
          className={`w-3 h-3 border  rounded m-auto block ${
            isSelected ? "bg-blue-500 border-blue-500" : " border-gray-500"
          }`}
        />
      </td>
      <td className="px-4 py-2 whitespace-nowrap">{item.method}</td>
      <td className="px-4 py-2 whitespace-nowrap max-w-lg overflow-hidden">
        {item.url}
      </td>
      <td className="px-4 whitespace-nowrap overflow-hidden max-w-0">
        {item.isReqText ? (
          <Code isInline>{item.reqBody}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="req-body"
            contentLength={item.reqLength}
          />
        )}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        {item.statusCode.toString()}
      </td>
      <td className="px-4 py-2 whitespace-nowrap overflow-hidden max-w-0">
        {item.isResText ? (
          <Code isInline>{item.resBody}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="res-body"
            contentLength={item.resLength}
          />
        )}
      </td>
      <td className="px-4 whitespace-nowrap text-center">
        <ProgressIcon ulid={item.ulid} />
      </td>
      <td className="px-4 whitespace-nowrap">
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
