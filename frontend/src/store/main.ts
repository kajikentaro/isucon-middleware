import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { executionProgressMap } from "./executionProgressMap";
import { recordedTransactions } from "./recordedTransactions";

const store = configureStore({
  reducer: {
    executionProgressMap,
    recordedTransactions,
  },
});

export default store;

export type State = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<State> = useSelector;
