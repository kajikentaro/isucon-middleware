import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ExecutionResponse } from "../types";
import { State } from "./main";

export type Status = "loading" | "fail" | "success" | "init";

export type ExecutionResponseMap = Record<string, ExecutionResponse>;

const initialState: ExecutionResponseMap = {};

const slice = createSlice({
  name: "executionResponse",
  initialState,
  reducers: {
    setExecutionResponse: (
      state,
      action: PayloadAction<{
        ulid: string;
        executeResponse: ExecutionResponse;
      }>
    ) => {
      return {
        ...state,
        [action.payload.ulid]: action.payload.executeResponse,
      };
    },
  },
  selectors: {
    selectExecutionResponseMap: (state) => state,
  },
});

export const executionResponse = slice.reducer;

export const { selectExecutionResponseMap } = slice.selectors;
export const selectExecutionResponse = (ulid: string) => (state: State) =>
  selectExecutionResponseMap(state)[ulid];

export const { setExecutionResponse } = slice.actions;
