import { ReactNode } from "react";
import TagEmpty from "./tag-empty";

interface Props {
  children: ReactNode;
  isInline?: boolean;
}
export default function Code({ children, isInline }: Props) {
  if (!children) {
    return <TagEmpty />;
  }

  if (isInline) {
    return (
      <code className="bg-gray-700 text-white p-1 text-xs">{children}</code>
    );
  }

  return (
    <code className="block bg-black text-white p-2 rounded-md my-2 whitespace-pre-line break-all">
      {children}
    </code>
  );
}
