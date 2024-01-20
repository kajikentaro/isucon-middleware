import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface Ui {
  isFetchingTransactions: boolean;
}

const initialState: Ui = {
  isFetchingTransactions: true,
};

const slice = createSlice({
  name: "ui",
  initialState,
  reducers: {
    setIsFetchingTransactions: (state, action: PayloadAction<boolean>) => {
      return { ...state, isFetchingTransactions: action.payload };
    },
  },
  selectors: {
    selectIsFetchingTransactions: (state) => state.isFetchingTransactions,
  },
});

export const ui = slice.reducer;

export const { selectIsFetchingTransactions } = slice.selectors;

export const { setIsFetchingTransactions } = slice.actions;
