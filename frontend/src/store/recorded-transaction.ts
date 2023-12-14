import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RecordedTransaction } from "../types";

export type Status = "loading" | "fail" | "success" | "init";

export type RecordedTransactionList = RecordedTransaction[];

const initialState: RecordedTransactionList = [];

const slice = createSlice({
  name: "recordedTransaction",
  initialState,
  reducers: {
    setRecordedTransactionList: (
      state,
      action: PayloadAction<RecordedTransactionList>
    ) => {
      return action.payload;
    },
  },
  selectors: {
    selectRecordedTransactionList: (state) => state,
  },
});

export const recordedTransaction = slice.reducer;

export const { selectRecordedTransactionList } = slice.selectors;

export const { setRecordedTransactionList } = slice.actions;
