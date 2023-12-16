import { BodyType, getBodyUrl } from "@/utils/get-url";
import Link from "next/link";

interface Props {
  ulid: string;
  type: BodyType;
  className?: string;
}

export function TagBinary({ type, ulid, className }: Props) {
  return (
    <Link
      href={getBodyUrl(type, ulid)}
      onClick={(e) => {
        e.stopPropagation();
      }}
      className={`inline-flex ${className ? className : ""}`}
    >
      <span className="bg-green-500 text-white p-2 text-xs rounded-full block w-fit">
        binary data
      </span>
    </Link>
  );
}
