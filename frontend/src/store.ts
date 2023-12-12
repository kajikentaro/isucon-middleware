import { configureStore, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { RecordedTransaction } from "./types";

export type State = {
  status: "fetching" | "fetched";
  recordedTransactions: RecordedTransaction[];
};

const initialState: State = {
  status: "fetching",
  recordedTransactions: [],
};

const recordedTransactionsSlice = createSlice({
  name: "recordedTransactions",
  initialState,
  reducers: {
    setRecordedTransactions: (
      state,
      action: PayloadAction<RecordedTransaction[]>
    ) => ({
      ...state,
      status: "fetched",
      recordedTransactions: action.payload,
    }),
  },
  selectors: {
    selectRecordedTransactions: (state) => state.recordedTransactions,
  },
});

export const { setRecordedTransactions } = recordedTransactionsSlice.actions;
export const { selectRecordedTransactions } =
  recordedTransactionsSlice.getSelectors();

const store = configureStore({
  reducer: recordedTransactionsSlice.reducer,
});

export default store;

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<State> = useSelector;
