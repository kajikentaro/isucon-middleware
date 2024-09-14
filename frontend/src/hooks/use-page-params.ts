"use client";

import { useSearchParams } from "next/navigation";

interface SearchParams {
  get: (name: string) => string | null;
}

function common(searchParams: SearchParams) {
  const pageStr = searchParams.get("page") || "1";
  const pageInt = Number(pageStr);
  const page = isNaN(pageInt) ? 1 : pageInt;
  const query = searchParams.get("query") || "";

  return { page, query };
}

// this function is for actions since we don't need useEffect in that case.
export function getPageParams() {
  const currentUrl = new URL(window.location.href);
  return common(currentUrl.searchParams);
}

export function usePageParams() {
  const searchParams = useSearchParams();
  return common(searchParams);
}
