import { combineSlices } from "@reduxjs/toolkit";
import { isFetchingTransactions } from "./is-fetching-transactions";
import { selectedUlids } from "./selected-ulids";

const ui = combineSlices(selectedUlids, isFetchingTransactions);

export default ui;
