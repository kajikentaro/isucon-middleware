import { BodyType, getBodyPath } from "@/utils/get-url";
import Link from "next/link";
import TagEmpty from "./tag-empty";

interface Props {
  ulid: string;
  type: BodyType;
  className?: string;
  contentLength: number;
}

function normalizeContentLength(contentLength: number): string {
  const kb = Number(contentLength) / 1024;
  const mb = kb / 1024;
  const gb = mb / 1024;
  const tb = gb / 1024;

  if (tb >= 1) {
    return tb.toFixed(2) + " TB";
  } else if (gb >= 1) {
    return gb.toFixed(2) + " GB";
  } else if (mb >= 1) {
    return mb.toFixed(2) + " MB";
  } else if (kb >= 1) {
    return kb.toFixed(2) + " KB";
  } else {
    return contentLength + " Bytes";
  }
}

export function TagBinary({
  type,
  ulid,
  className = "",
  contentLength,
}: Props) {
  if (contentLength === 0) {
    return (
      <span className={`${className} inline-block`}>
        <TagEmpty />
      </span>
    );
  }

  return (
    <Link
      href={getBodyPath(type, ulid)}
      onClick={(e) => {
        e.stopPropagation();
      }}
      className={`inline-flex ${className}`}
      prefetch={false}
    >
      <span className="bg-green-500 text-white py-1 px-2 text-xs rounded-full block w-fit">
        binary data {contentLength ? normalizeContentLength(contentLength) : ""}
      </span>
    </Link>
  );
}
