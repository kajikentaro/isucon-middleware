"use client";
import { useExecute } from "@/hooks/use-execute";
import { useOpenPopup } from "@/hooks/use-open-popup";
import Code from "@/parts/code";
import { TagBinary } from "@/parts/tag-binary";
import { useAppSelector } from "@/store/main";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { MouseEvent } from "react";
import ProgressIcon from "../parts/progress-icon";

interface Props {
  ulid: string;
  handleCheckboxClick: (event: MouseEvent) => void;
  isSelected: boolean;
}

export default function TableRow(props: Props) {
  const { handleCheckboxClick, isSelected, ulid } = props;

  const item = useAppSelector(selectRecordedTransaction(ulid));
  const onExecute = useExecute(item.Ulid);
  const openPopup = useOpenPopup(item.Ulid);

  return (
    <tr className="border-b hover:bg-gray-100" onClick={openPopup}>
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
      <td className="px-4 py-2 whitespace-nowrap max-w-lg overflow-hidden">
        {item.Url}
      </td>
      <td className="px-4 whitespace-nowrap overflow-hidden max-w-0">
        {item.IsReqText ? (
          <Code isInline>{item.ReqBody}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="req-body"
            contentLength={item.ReqLength}
          />
        )}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        {item.StatusCode.toString()}
      </td>
      <td className="px-4 py-2 whitespace-nowrap overflow-hidden max-w-0">
        {item.IsResText ? (
          <Code isInline>{item.ResBody}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="res-body"
            contentLength={item.ResLength}
          />
        )}
      </td>
      <td className="px-4 whitespace-nowrap text-center">
        <ProgressIcon ulid={item.Ulid} />
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
