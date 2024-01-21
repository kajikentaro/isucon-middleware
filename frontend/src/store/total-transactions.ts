import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export type totalTransactions = number;
const initialState: totalTransactions = -1;

const slice = createSlice({
  name: "totalTransactions",
  initialState,
  reducers: {
    setTotalTransactions: (state, action: PayloadAction<number>) => {
      return action.payload;
    },
  },
  selectors: {
    selectTotalTransactions: (state) => state,
  },
});

export const totalTransactions = slice.reducer;

export const { selectTotalTransactions } = slice.selectors;
export const { setTotalTransactions } = slice.actions;
