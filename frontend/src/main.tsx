import ComparisonPopup from "@/components/comparison-popup";
import FilterPopup from "@/components/filter-popup";
import Main from "@/components/main";
import store from "@/store";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import "./globals.css";

createRoot(document.getElementById("root")!).render(
  <Provider store={store}>
    <Main />
    <ComparisonPopup />
    <FilterPopup />
  </Provider>
);
