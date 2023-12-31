import { createSelector, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RecordedTransaction } from "../types";
import { State } from "./main";

const initialState: Record<string, RecordedTransaction> = {};

const slice = createSlice({
  name: "recordedTransaction",
  initialState,
  reducers: {
    setRecordedTransactionList: (
      state,
      action: PayloadAction<RecordedTransaction[]>
    ) => {
      const newState: typeof state = {};
      for (const v of action.payload) {
        newState[v.Ulid] = v;
      }
      return newState;
    },
  },
  selectors: {
    selectRecordedTransactionMap: (state) => state,
  },
});

export const recordedTransaction = slice.reducer;

export const { selectRecordedTransactionMap } = slice.selectors;
export const selectRecordedTransaction = (ulid: string) => (state: State) =>
  selectRecordedTransactionMap(state)[ulid];
export const selectRecordedTransactionUlids = createSelector(
  [selectRecordedTransactionMap],
  (state) => Object.keys(state)
);

export const { setRecordedTransactionList } = slice.actions;
