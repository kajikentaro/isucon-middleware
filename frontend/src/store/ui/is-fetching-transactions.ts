import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from "..";

export type IsFetchingTransactions = boolean;
const initialState: IsFetchingTransactions = true;

const slice = createSlice({
  name: "isFetchingTransactions",
  initialState,
  reducers: {
    setIsFetchingTransactions: (state, action: PayloadAction<boolean>) => {
      return action.payload;
    },
  },
  selectors: {
    selectIsFetchingTransactions: (state) => state,
  },
});

export const isFetchingTransactions = slice;

export const { selectIsFetchingTransactions } = slice.getSelectors(
  (state: State) => state.ui.isFetchingTransactions
);

export const { setIsFetchingTransactions } = slice.actions;
