"use client";
import { fetchTransactions } from "@/actions/fetch-transactions";
import { ENV } from "@/constants";
import { useAppDispatch } from "@/store";
import { getRemoveAllURL } from "@/utils/get-url";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { CiTrash } from "react-icons/ci";

export default function RemoveAllButton() {
  const [isRemoving, setIsRemoving] = useState(false);
  const router = useRouter();
  const dispatch = useAppDispatch();

  const className =
    "border p-2 rounded flex items-center w-36 justify-center gap-1 text-red-500 border-red-500";

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
      router.replace(ENV.TOP_PAGE_PATH);
      dispatch(fetchTransactions(1, ""));
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
