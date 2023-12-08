"use client";
import { useState, useEffect, MouseEvent } from "react";

interface FetchResponse {
  IsReqText: boolean;
  IsResText: boolean;
  StatusCode: number;
  Ulid: string;
  ResBody: string;
  ResHeader: { [key: string]: string[] };
  ReqBody: string;
  ReqOthers: {
    Url: string;
    Header: { [key: string]: string[] };
    Method: string;
  };
}

export default function Home() {
  const [data, setData] = useState<FetchResponse[]>([]);
  const [selected, setSelected] = useState<boolean[]>([]);
  const [lastSelectedIndex, setLastSelectedIndex] = useState(-1);

  const isAllSelected = selected.every((s) => s) && selected.length > 0;

  const handleCheckboxClick = (
    event: MouseEvent<HTMLDivElement, globalThis.MouseEvent>,
    index: number
  ) => {
    if (!selected.length) return;

    const newSelected = [...selected];

    if (event.shiftKey) {
      // Shift-click: select all rows in range
      const nextIsTrue = !selected[index];
      for (
        let i = Math.min(lastSelectedIndex, index);
        i <= Math.max(lastSelectedIndex, index);
        i++
      ) {
        newSelected[i] = nextIsTrue;
      }
    } else {
      // Select only clicked row
      newSelected[index] = !selected[index];
    }

    setSelected(newSelected);
    setLastSelectedIndex(index);

    event.preventDefault();
    event.stopPropagation();
  };

  const fetchData = async () => {
    const response = await fetch("http://localhost:8080/fetch-all", {});
    const json: FetchResponse[] = await response.json();
    setData(json);
    setSelected(Array(json.length).fill(false));
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div className="flex flex-col justify-center items-center">
      <h1 className="text-3xl font-bold mb-4">Fetch All</h1>
      <table className="table-auto border-collapse w-full">
        <thead>
          <tr className="border-b bg-gray-100 text-gray-600">
            <th
              className="px-4 py-2 whitespace-nowrap"
              onClick={() => {
                setSelected(Array(selected.length).fill(!isAllSelected));
              }}
            >
              <div
                className={`w-4 h-4 border border-gray-500 rounded m-auto block ${
                  isAllSelected ? "bg-blue-500" : "bg-white"
                }`}
              />
            </th>
            <th className="px-4 py-2 whitespace-nowrap">Method</th>
            <th className="px-4 py-2 whitespace-nowrap">URL</th>
            <th className="px-4 py-2 w-1/2">ReqBody</th>
            <th className="px-4 py-2">Status Code</th>
            <th className="px-4 py-2 w-1/2">ResBody</th>
            <th className="px-4 py-2 w-1/2">Execute</th>
          </tr>
        </thead>
        <tbody>
          {data.map((item, index) => (
            <tr key={index} className="border-b hover:bg-gray-100">
              <td
                className="px-4 py-2 whitespace-nowrap"
                onClick={(e) => handleCheckboxClick(e, index)}
              >
                <div
                  className={`w-3 h-3 border  rounded m-auto block ${
                    selected[index]
                      ? "bg-blue-500 border-blue-500"
                      : " border-gray-500"
                  }`}
                />
              </td>
              <td className="px-4 py-2 whitespace-nowrap">
                {item.ReqOthers.Method}
              </td>
              <td className="px-4 py-2 whitespace-nowrap">
                {item.ReqOthers.Url}
              </td>
              <td className="px-4 py-2 whitespace-nowrap">{item.ReqBody}</td>
              <td className="px-4 py-2 whitespace-nowrap">
                {item.StatusCode.toString()}
              </td>
              <td className="px-4 py-2 whitespace-nowrap">{item.ResBody}</td>
              <td className="px-4 py-2 whitespace-nowrap">
                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded-full flex items-center m-auto">
                  <svg
                    className="h-4 w-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M9 5l8 8-8 8"
                    />
                  </svg>
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
