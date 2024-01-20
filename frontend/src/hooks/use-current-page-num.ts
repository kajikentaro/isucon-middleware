"use client";
import { useSearchParams } from "next/navigation";
import { useMemo } from "react";

export function useCurrentPageNum() {
  const searchParams = useSearchParams();

  const currentPageNum = useMemo(() => {
    const page = searchParams.get("page") || "1";

    const pageInt = Number(page);
    if (isNaN(pageInt)) {
      return 1;
    }
    return pageInt;
  }, [searchParams]);

  return currentPageNum;
}
