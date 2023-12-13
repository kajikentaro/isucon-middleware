import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RecordedTransaction } from "../types";

export type Status = "loading" | "fail" | "success" | "init";

export type RecordedTransactions = RecordedTransaction[];

const initialState: RecordedTransactions = [];

const recordedTransactionsSlice = createSlice({
  name: "recordedTransactions",
  initialState,
  reducers: {
    setRecordedTransactions: (
      state,
      action: PayloadAction<RecordedTransaction[]>
    ) => {
      return action.payload;
    },
  },
  selectors: {
    selectRecordedTransactions: (state) => state,
  },
});

export const recordedTransactions = recordedTransactionsSlice.reducer;

export const { selectRecordedTransactions } =
  recordedTransactionsSlice.selectors;

export const { setRecordedTransactions } = recordedTransactionsSlice.actions;
