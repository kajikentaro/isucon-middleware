"use client";
import { fetchTransactions } from "@/actions/fetch-transactions";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectSelectedUlids } from "@/store/ui/selected-ulids";
import { getRemoveURL } from "@/utils/get-url";
import { useState } from "react";
import { CiTrash } from "react-icons/ci";

export default function RemoveSelectedButton() {
  const [isRemoving, setIsRemoving] = useState(false);
  const selectedUlids = useAppSelector(selectSelectedUlids);
  const dispatch = useAppDispatch();

  const className =
    "border p-2 rounded flex items-center w-44 justify-center gap-1 text-red-500 border-red-500";

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
    dispatch(fetchTransactions());
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
