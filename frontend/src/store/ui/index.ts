import { combineSlices } from "@reduxjs/toolkit";
import { comparisonPopup } from "./comparison-popup";
import { filterPopup } from "./filter-popup";
import { isFetchingTransactions } from "./is-fetching-transactions";
import { selectedUlids } from "./selected-ulids";

const ui = combineSlices(
  selectedUlids,
  isFetchingTransactions,
  comparisonPopup,
  filterPopup
);

export default ui;
