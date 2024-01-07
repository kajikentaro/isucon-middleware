import { showComparisonPopup } from "@/store/comparison-popup";
import { useAppDispatch } from "@/store/main";
import { useRouter } from "next/navigation";

export const POPUP_SEARCH_PARAM = "popup";

export function useOpenPopup(ulid: string) {
  const dispatch = useAppDispatch();
  const router = useRouter();

  return () => {
    dispatch(showComparisonPopup(ulid));

    // add a query parameter
    const queryParams = new URLSearchParams(window.location.search);
    queryParams.set(POPUP_SEARCH_PARAM, "true");
    router.push(`/?${queryParams.toString()}`);
  };
}

export function useClosePopup() {
  return () => {
    // remove a query parameter
  };
}
