import { ReactNode } from "react";

interface Props {
  children: ReactNode;
}
export default function Code({ children }: Props) {
  return (
    <code className="block bg-black text-white p-2 rounded-md my-2 whitespace-pre-line break-all">
      {children}
    </code>
  );
}
