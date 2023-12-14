import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface ComparisonPopup {
  isVisible: boolean;
  ulid: string;
}

const initialState: ComparisonPopup = {
  isVisible: false,
  ulid: "",
};

const slice = createSlice({
  name: "comparisonPopup",
  initialState,
  reducers: {
    showComparisonPopup: (state, action: PayloadAction<string>) => {
      return { isVisible: true, ulid: action.payload };
    },
    closeComparisonPopup: (state, action: PayloadAction) => {
      return initialState;
    },
  },
  selectors: {
    selectComparisonPopup: (state) => state,
  },
});

export const comparisonPopup = slice.reducer;

export const { selectComparisonPopup } = slice.selectors;

export const { showComparisonPopup, closeComparisonPopup } = slice.actions;
