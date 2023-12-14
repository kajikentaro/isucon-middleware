import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { executionProgress } from "./execution-progress";
import { recordedTransaction } from "./recorded-transaction";
import { executionResponse } from "./execution-response";
import { comparisonPopup } from "./comparison-popup";

const store = configureStore({
  reducer: {
    executionProgress,
    recordedTransaction,
    executionResponse,
    comparisonPopup,
  },
});

export default store;

export type State = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<State> = useSelector;
