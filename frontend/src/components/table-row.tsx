"use client";
import { RecordedTransaction } from "@/types";
import { MouseEvent } from "react";
import { useOnExecute } from "@/hooks/use-queue";
import ProgressIcon from "./progress-icon";

interface Props {
  item: RecordedTransaction;
  handleCheckboxClick: (event: MouseEvent) => void;
  isSelected: boolean;
}

export default function TableRow(props: Props) {
  const { handleCheckboxClick, isSelected, item } = props;

  const onExecute = useOnExecute(item.Ulid);

  return (
    <tr className="border-b hover:bg-gray-100">
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
      <td className="px-4 py-2 whitespace-nowrap">{item.ReqOthers.Method}</td>
      <td className="px-4 py-2 whitespace-nowrap">{item.ReqOthers.Url}</td>
      <td className="px-4 py-2 whitespace-nowrap">{item.ReqBody}</td>
      <td className="px-4 py-2 whitespace-nowrap">
        {item.StatusCode.toString()}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">{item.ResBody}</td>
      <td className="px-4 py-2 whitespace-nowrap">
        <ProgressIcon ulid={item.Ulid} />
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded-full flex items-center m-auto"
          onClick={onExecute}
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
