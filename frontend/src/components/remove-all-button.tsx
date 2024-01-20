"use client";

import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { getRemoveAllURL } from "@/utils/get-url";
import { useState } from "react";
import { CiTrash } from "react-icons/ci";

export default function RemoveAllButton() {
  const [isRemoving, setIsRemoving] = useState(false);
  const { fetchTransactions } = useFetchTransactions();
  const currentPageNum = useCurrentPageNum();

  const className =
    "border py-1 px-2 rounded flex items-center w-36 justify-center gap-1 text-red-500 border-red-500";

  const handleClick = async () => {
    const shouldProceed = confirm(
      "Are you sure you want to proceed with removal?"
    );
    if (!shouldProceed) {
      return;
    }

    setIsRemoving(true);
    const res = await fetch(getRemoveAllURL(), { method: "POST" });
    setIsRemoving(false);
    if (res.status === 200) {
      await fetchTransactions(currentPageNum);
    }
  };

  if (isRemoving) {
    return (
      <button className={className} disabled>
        Removing. . .
      </button>
    );
  }

  return (
    <button className={className} onClick={handleClick}>
      <CiTrash />
      Remove All
    </button>
  );
}
