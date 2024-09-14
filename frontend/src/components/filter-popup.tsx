import { fetchTransactions } from "@/actions/fetch-transactions";
import { usePageParams } from "@/hooks/use-page-params";
import Modal from "@/parts/modal";
import { useAppDispatch, useAppSelector } from "@/store";
import { closeFilterPopup, selectFilterPopup } from "@/store/ui/filter-popup";
import { useRouter } from "next/navigation";
import { ChangeEventHandler, useState } from "react";

export default function FilterPopup() {
  const dispatch = useAppDispatch();

  const closePopup = () => {
    dispatch(closeFilterPopup());
  };

  const { isVisible } = useAppSelector(selectFilterPopup);

  return (
    <Modal isVisible={isVisible} closePopup={closePopup} title={"Filter URL"}>
      <ModalContents />
    </Modal>
  );
}

function ModalContents() {
  const dispatch = useAppDispatch();
  const router = useRouter();

  const { query: initialFilterText } = usePageParams();
  const [filterText, setFilterText] = useState(initialFilterText);

  const onChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    setFilterText(e.target.value);
  };

  const onApply = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    dispatch(closeFilterPopup());

    const newUrl = new URL(window.location.href);
    newUrl.searchParams.set("query", filterText);
    router.push(newUrl.toString());

    dispatch(fetchTransactions(undefined, filterText));
  };

  return (
    <form
      className="flex flex-center mt-6 mb-5 mx-auto max-w-screen-sm"
      onSubmit={onApply}
    >
      <input
        type="text"
        value={filterText}
        onChange={onChange}
        placeholder="Enter URL to filter"
        className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm"
        autoFocus
      />
      <button
        onClick={onApply}
        className="flex py-2 px-4 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition duration-300"
      >
        Apply
      </button>
    </form>
  );
}
