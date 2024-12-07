"use client";
import { useExecute } from "@/hooks/use-execute";
import { useOpenPopup } from "@/hooks/use-open-popup";
import Code from "@/parts/code";
import { ExecuteButton } from "@/parts/execute-button";
import { TagBinary } from "@/parts/tag-binary";
import { useAppSelector } from "@/store";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { joinClassName } from "@/utils/join-class-name";
import { MouseEvent, ReactNode } from "react";
import ProgressIcon from "../parts/progress-icon";

interface Props {
  ulid: string;
  onCheckboxClick: (event: MouseEvent) => void;
  isSelected: boolean;
}

const MAX_BODY_CHARACTER = 150;

export default function TableRow(props: Props) {
  const { onCheckboxClick, isSelected, ulid } = props;

  const item = useAppSelector(selectRecordedTransaction(ulid));
  const onExecute = useExecute(item.ulid);
  const openPopup = useOpenPopup(item.ulid);

  const classNameTd = "px-2 py-1 whitespace-nowrap max-w-0 overflow-hidden";

  return (
    <tr className="border-b hover:bg-gray-100" onClick={openPopup}>
      <td className={classNameTd} onClick={onCheckboxClick}>
        <Center>
          <div
            className={`w-3 h-3 border  rounded m-auto block ${
              isSelected ? "bg-blue-500 border-blue-500" : " border-gray-500"
            }`}
          />
        </Center>
      </td>
      <td className={classNameTd}>{item.method}</td>
      <td className={classNameTd}>{item.url}</td>
      <td className={classNameTd}>
        {item.isReqText ? (
          <Code isInline>{item.reqBody.slice(0, MAX_BODY_CHARACTER)}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="req-body"
            contentLength={item.reqLength}
          />
        )}
      </td>
      <td className="px-4 py-2 whitespace-nowrap">
        <Center>{item.statusCode.toString()}</Center>
      </td>
      <td className={classNameTd}>
        {item.isResText ? (
          <Code isInline>{item.resBody.slice(0, MAX_BODY_CHARACTER)}</Code>
        ) : (
          <TagBinary
            ulid={ulid}
            type="res-body"
            contentLength={item.resLength}
          />
        )}
      </td>
      <td className={joinClassName(classNameTd, "overflow-visible")}>
        <Center>
          <ProgressIcon ulid={item.ulid} />
        </Center>
      </td>
      <td className={classNameTd}>
        <Center>
          <ExecuteButton
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
          </ExecuteButton>
        </Center>
      </td>
    </tr>
  );
}

function Center({ children }: { children: ReactNode }) {
  return <div className="flex justify-center items-center">{children}</div>;
}
