import { ReactNode } from "react";
import TagEmpty from "./tag-empty";

interface Props {
  children: ReactNode;
  isInline?: boolean;
  className?: string;
}

export default function Code({ children, isInline, className = "" }: Props) {
  if (!children) {
    return (
      <span className={`${className} inline-block`}>
        <TagEmpty />
      </span>
    );
  }

  if (isInline) {
    return (
      <code className={`${className} bg-gray-700 text-white p-1 text-xs`}>
        {children}
      </code>
    );
  }

  return (
    <code
      className={`${className} block bg-black text-white p-2 rounded-md whitespace-pre-line break-all`}
    >
      {children}
    </code>
  );
}
