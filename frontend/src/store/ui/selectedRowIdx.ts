import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { State } from "..";

export type SelectedUlids = string[];
const initialState: SelectedUlids = [];

const slice = createSlice({
  name: "selectedUlids",
  initialState,
  reducers: {
    setSelectedUlids: (state, action: PayloadAction<string[]>) => {
      return action.payload;
    },
  },
  selectors: {
    selectSelectedUlids: (state) => state,
  },
});

export const selectedUlids = slice;

export const { selectSelectedUlids } = slice.getSelectors(
  (state: State) => state.ui.selectedUlids
);

export const { setSelectedUlids } = slice.actions;
