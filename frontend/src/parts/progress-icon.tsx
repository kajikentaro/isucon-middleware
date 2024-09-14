"use client";
import ThreeDotsAnimation from "@/parts/3-dots-animation";
import { useAppSelector } from "@/store";
import { selectExecutionProgressMap } from "@/store/execution-progress";
import { shouldBeNever } from "@/utils/assert-never";
import {
  AiOutlineCheckCircle,
  AiOutlineCloseCircle,
  AiOutlineLoading,
  AiOutlineWarning,
} from "react-icons/ai";

interface Props {
  ulid: string;
}

export default function ProgressIcon(props: Props) {
  const status =
    useAppSelector(selectExecutionProgressMap)[props.ulid] || "fail";

  const commonProps: (typeof AiOutlineLoading)["Props"] = {
    size: "25",
  };
  let icon = <></>;
  let tooltipText = "";

  switch (status) {
    case "statusCodeNotSame":
      icon = <AiOutlineWarning {...commonProps} color="red" />;
      tooltipText = "Status code was not same";
      break;
    case "bodyNotSame":
      icon = <AiOutlineCheckCircle {...commonProps} color="red" />;
      tooltipText = "Response body was not same";
      break;
    case "headerNotSame":
      icon = <AiOutlineCheckCircle {...commonProps} color="orange" />;
      tooltipText = "Headers were not same";
      break;

    case "waitingResponse":
      icon = <AiOutlineLoading {...commonProps} className="animate-spin" />;
      break;
    case "fail":
      icon = <AiOutlineCloseCircle {...commonProps} />;
      tooltipText = "Request failed";
      break;
    case "success":
      icon = <AiOutlineCheckCircle {...commonProps} color="green" />;
      tooltipText = "Response body and headers were same";
      break;
    case "waitingQueue":
      icon = <ThreeDotsAnimation {...commonProps} />;
      break;
    case "init":
      return <></>;
    default:
      shouldBeNever(status);
  }

  return (
    <div className="relative inline-flex align-middle text-center">
      {icon}

      {tooltipText && (
        <>
          <div className="opacity-0 pointer-events-none absolute z-10 bg-black text-white text-xs py-1 px-2 rounded whitespace-nowrap top-full -left-1/2 transform -translate-x-1/2 transition-opacity duration-300">
            {tooltipText}
          </div>
          <style>
            {`
          .relative:hover .absolute {
            opacity: 1;
          }
          `}
          </style>
        </>
      )}
    </div>
  );
}
