"use client";

import { useCurrentPageNum } from "@/hooks/use-current-page-num";
import { useFetchTransactions } from "@/hooks/use-fetch-transactions";
import { getRemoveURL } from "@/utils/get-url";
import { useState } from "react";
import { CiTrash } from "react-icons/ci";

interface Props {
  selectedUlids: string[];
}

export default function RemoveSelectedButton({ selectedUlids }: Props) {
  const [isRemoving, setIsRemoving] = useState(false);
  const { fetchTransactions } = useFetchTransactions();
  const currentPageNum = useCurrentPageNum();

  const className =
    "border py-1 px-2 rounded flex items-center w-44 justify-center gap-1 text-red-500 border-red-500";

  const handleClick = async () => {
    const shouldProceed = confirm(
      "Are you sure you want to proceed with removal?"
    );
    if (!shouldProceed) {
      return;
    }

    setIsRemoving(true);
    for (const ulid of selectedUlids) {
      await fetch(getRemoveURL(ulid), { method: "POST" });
    }
    setIsRemoving(false);
    await fetchTransactions(currentPageNum);
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
      Remove Selected
    </button>
  );
}
