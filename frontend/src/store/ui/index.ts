import { combineSlices } from "@reduxjs/toolkit";
import { isFetchingTransactions } from "./is-fetching-transactions";
import { selectedUlids } from "./selectedRowIdx";

const ui = combineSlices(selectedUlids, isFetchingTransactions);

export default ui;
