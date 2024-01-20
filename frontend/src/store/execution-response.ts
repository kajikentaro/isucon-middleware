import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from ".";
import { ExecutionResponse } from "../types";

const initialState: Record<string, ExecutionResponse> = {};

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
