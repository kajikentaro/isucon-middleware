import { ExecutionProgress } from "@/store/execution-progress";
import { shouldBeNever } from "@/utils/assert-never";
import React from "react";
import {
  AiOutlineLoading,
  AiOutlineCloseCircle,
  AiOutlineCheckCircle,
} from "react-icons/ai";

type Props = {
  status: ExecutionProgress;
};

const StatusIcon: React.FC<Props> = ({ status }) => {
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
};

export default StatusIcon;
