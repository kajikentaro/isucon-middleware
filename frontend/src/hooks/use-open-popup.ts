import { showComparisonPopup } from "@/store/comparison-popup";
import { useAppDispatch } from "@/store/main";

export const POPUP_SEARCH_PARAM = "popup";

export function useOpenPopup(ulid: string) {
  const dispatch = useAppDispatch();

  return () => {
    dispatch(showComparisonPopup(ulid));
  };
}
