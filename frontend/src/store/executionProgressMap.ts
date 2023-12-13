import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export type ExecutionProgress = "loading" | "fail" | "success" | "init";

export type ExecutionProgressMap = Record<string, ExecutionProgress>;

const initialState: ExecutionProgressMap = {};

const executionProgressMapSlice = createSlice({
  name: "executionProgressMap",
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

export const executionProgressMap = executionProgressMapSlice.reducer;

export const { selectExecutionProgressMap } =
  executionProgressMapSlice.selectors;
export const { setExecutionProgress, setExecutionProgressAll } =
  executionProgressMapSlice.actions;
