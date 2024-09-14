import { fetchTransactions } from "@/actions/fetch-transactions";
import { MAX_ROW_LENGTH } from "@/constants";
import { usePageParams } from "@/hooks/use-page-params";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectTotalTransactions } from "@/store/total-transactions";
import { selectIsFetchingTransactions } from "@/store/ui/is-fetching-transactions";
import Link from "next/link";

// the number of active buttons before and after the current page
const ACTIVE_BUTTON_LENGTH = 2;

export default function Pagination() {
  const isFetchingTransactions = useAppSelector(selectIsFetchingTransactions);
  const totalTransactions = useAppSelector(selectTotalTransactions);
  const maxPageNum = Math.ceil(totalTransactions / MAX_ROW_LENGTH);
  const { page: currentPageNum, query } = usePageParams();
  const dispatch = useAppDispatch();

  const getLinkProps = (pageNum: number) => ({
    href: {
      query: { page: pageNum, query },
    },
    onClick: () => dispatch(fetchTransactions(pageNum)),
  });

  const shouldShow = (pageNum: number) => {
    const difference = Math.abs(pageNum - currentPageNum);
    return (
      difference <= ACTIVE_BUTTON_LENGTH ||
      pageNum === 1 ||
      pageNum === maxPageNum
    );
  };

  if (isFetchingTransactions) {
    return null;
  }

  if (maxPageNum <= 1) {
    return null;
  }

  return (
    <nav className="my-4">
      <ul className="flex items-center gap-3 bg-gray-200 rounded-full px-3">
        {currentPageNum > 1 && (
          <li className="flex">
            <Link
              {...getLinkProps(currentPageNum - 1)}
              prefetch={false}
              className="font-bold py-2 px-3 rounded-lg"
            >
              &lt;
            </Link>
          </li>
        )}

        {Array.from({ length: maxPageNum }, (_, i) => i + 1).map((pageNum) => {
          if (!shouldShow(pageNum)) {
            if (pageNum === 2 || pageNum === maxPageNum - 1) {
              return (
                <li key={pageNum} className="flex">
                  <span>...</span>
                </li>
              );
            }
            return null;
          }

          if (pageNum === currentPageNum) {
            return (
              <li
                key={pageNum}
                className="text-white w-8 h-7 rounded-full bg-slate-500 flex items-center justify-center"
              >
                <span>{pageNum}</span>
              </li>
            );
          }

          return (
            <li key={pageNum} className="flex">
              <Link
                {...getLinkProps(pageNum)}
                prefetch={false}
                className="font-bold py-2 px-3 rounded-lg"
              >
                {pageNum}
              </Link>
            </li>
          );
        })}

        {currentPageNum < maxPageNum && (
          <li className="flex">
            <Link
              {...getLinkProps(currentPageNum + 1)}
              prefetch={false}
              className="font-bold py-2 px-3 rounded-lg"
            >
              &gt;
            </Link>
          </li>
        )}
      </ul>
    </nav>
  );
}
