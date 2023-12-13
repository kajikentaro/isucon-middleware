"use client";
import { shouldBeNever } from "@/utils/assert-never";
import {
  AiOutlineLoading,
  AiOutlineCloseCircle,
  AiOutlineCheckCircle,
} from "react-icons/ai";
import { useAppSelector } from "@/store/main";
import { selectExecutionProgressMap } from "@/store/executionProgressMap";

interface Props {
  ulid: string;
}

export default function ProgressIcon(props: Props) {
  const status =
    useAppSelector(selectExecutionProgressMap)[props.ulid] || "fail";

  switch (status) {
    case "loading":
      return <AiOutlineLoading className="animate-spin" />;
    case "fail":
      return <AiOutlineCloseCircle />;
    case "success":
      return <AiOutlineCheckCircle />;
    case "init":
      return <></>;
    default:
      shouldBeNever(status);
  }
}
