"use client";
import ComparisonPopup from "@/components/comparison-popup";
import Main from "@/components/main";
import store from "@/store";
import { Provider } from "react-redux";

export default function Home() {
  return (
    <Provider store={store}>
      <Main />
      <ComparisonPopup />
    </Provider>
  );
}
