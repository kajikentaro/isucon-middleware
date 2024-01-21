import { useAppDispatch } from "@/store";
import { showComparisonPopup } from "@/store/ui/comparison-popup";

export const POPUP_SEARCH_PARAM = "popup";

export function useOpenPopup(ulid: string) {
  const dispatch = useAppDispatch();

  return () => {
    dispatch(showComparisonPopup(ulid));
  };
}
