import { BodyType, getBodyUrl } from "@/utils/get-url";
import Link from "next/link";

interface Props {
  ulid: string;
  type: BodyType;
}

export function TagThisIsBinary({ type, ulid }: Props) {
  return (
    <Link
      href={getBodyUrl(type, ulid)}
      onClick={(e) => {
        e.stopPropagation();
      }}
      className="inline-flex"
    >
      <span className="bg-green-500 text-white p-2 text-xs rounded-full block w-fit">
        This is binary data
      </span>
    </Link>
  );
}
