import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from "..";

export interface FilterPopup {
  isVisible: boolean;
}

const initialState: FilterPopup = {
  isVisible: false,
};

const slice = createSlice({
  name: "filterPopup",
  initialState,
  reducers: {
    showFilterPopup: (state, action: PayloadAction) => {
      return { ...state, isVisible: true };
    },
    closeFilterPopup: (state, action: PayloadAction) => {
      return { ...state, isVisible: false };
    },
  },
  selectors: {
    selectFilterPopup: (state) => state,
  },
});

export const filterPopup = slice;

export const { selectFilterPopup: selectFilterPopup } = slice.getSelectors(
  (state: State) => state.ui.filterPopup
);

export const { showFilterPopup, closeFilterPopup } = slice.actions;
