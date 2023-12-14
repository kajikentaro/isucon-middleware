import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { executionProgress } from "./execution-progress";
import { recordedTransaction } from "./recorded-transaction";
import { executionResponse } from "./execution-response";

const store = configureStore({
  reducer: {
    executionProgress,
    recordedTransaction,
    executionResponse,
  },
});

export default store;

export type State = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<State> = useSelector;
