import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from "./main";

export type ExecutionProgress = "loading" | "fail" | "success" | "init";

export type ExecutionProgressMap = Record<string, ExecutionProgress>;

const initialState: ExecutionProgressMap = {};

const slice = createSlice({
  name: "executionProgress",
  initialState,
  reducers: {
    setExecutionProgressAll: (
      state,
      action: PayloadAction<ExecutionProgressMap>
    ) => {
      return action.payload;
    },
    setExecutionProgress: (
      state,
      action: PayloadAction<{
        ulid: string;
        executionProgress: ExecutionProgress;
      }>
    ) => {
      return {
        ...state,
        [action.payload.ulid]: action.payload.executionProgress,
      };
    },
  },
  selectors: {
    selectExecutionProgressMap: (state) => state,
  },
});

export const executionProgress = slice.reducer;

export const { selectExecutionProgressMap } = slice.selectors;
export const selectExecutionProgress = (ulid: string) => (state: State) =>
  selectExecutionProgressMap(state)[ulid];

export const { setExecutionProgress, setExecutionProgressAll } = slice.actions;
