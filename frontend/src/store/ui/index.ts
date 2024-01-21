import { combineSlices } from "@reduxjs/toolkit";
import { comparisonPopup } from "./comparison-popup";
import { isFetchingTransactions } from "./is-fetching-transactions";
import { selectedUlids } from "./selected-ulids";

const ui = combineSlices(
  selectedUlids,
  isFetchingTransactions,
  comparisonPopup
);

export default ui;
