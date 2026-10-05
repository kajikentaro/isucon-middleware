import ThreeDotsAnimation from "@/parts/3-dots-animation";
import Modal from "@/parts/modal";
import { ExportSize } from "@/types";
import { getExportSizeURL, getExportURL } from "@/utils/get-url";
import { useEffect, useState } from "react";
import { CiExport } from "react-icons/ci";

const DEFAULT_LIMIT_MB = "100";
const MB = 1024 * 1024;

export default function ExportButton() {
  const [isVisible, setIsVisible] = useState(false);

  return (
    <>
      <button
        className="border p-2 rounded flex items-center w-36 justify-center gap-1"
        onClick={() => setIsVisible(true)}
      >
        <CiExport />
        Export
      </button>
      <Modal
        isVisible={isVisible}
        closePopup={() => setIsVisible(false)}
        title="Export as SQLite"
      >
        <ModalContents />
      </Modal>
    </>
  );
}

function ModalContents() {
  const [size, setSize] = useState<ExportSize | null>(null);
  const [sizeError, setSizeError] = useState(false);
  const [metaMB, setMetaMB] = useState(DEFAULT_LIMIT_MB);
  const [bodyMB, setBodyMB] = useState(DEFAULT_LIMIT_MB);

  useEffect(() => {
    fetch(getExportSizeURL())
      .then((res) => {
        if (!res.ok) {
          throw new Error(`status ${res.status}`);
        }
        return res.json();
      })
      .then((json: ExportSize) => setSize(json))
      .catch(() => setSizeError(true));
  }, []);

  const isValid = isLimit(metaMB) && isLimit(bodyMB);

  const onDownload = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    if (!isValid) {
      return;
    }
    // let the browser download it directly, not to hold a large file in memory
    window.location.href = getExportURL(Number(metaMB), Number(bodyMB));
  };

  return (
    <form className="mt-6 flex flex-col gap-4 max-w-lg" onSubmit={onDownload}>
      <RecordedSize size={size} isError={sizeError} />

      <LimitInput
        label="Metadata limit (MB)"
        value={metaMB}
        onChange={setMetaMB}
        totalBytes={size?.metaBytes}
      />
      <LimitInput
        label="Body limit (MB)"
        value={bodyMB}
        onChange={setBodyMB}
        totalBytes={size?.textBodyBytes}
      />

      <p className="text-sm text-gray-600">
        Both are filled from the oldest request. Rows and bodies after the
        limit are omitted. Only text bodies are embedded, and gzip-compressed
        ones are decompressed. The total of text bodies above counts them by
        the compressed size, so the actual size is larger.
      </p>

      <button
        type="submit"
        disabled={!isValid}
        className="self-end py-2 px-4 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition duration-300 disabled:opacity-50"
      >
        Download
      </button>
    </form>
  );
}

function RecordedSize({
  size,
  isError,
}: {
  size: ExportSize | null;
  isError: boolean;
}) {
  if (isError) {
    return <p className="text-red-500">Failed to fetch the recorded size.</p>;
  }
  if (size === null) {
    return (
      <div className="flex items-center gap-2">
        Recorded: <ThreeDotsAnimation />
      </div>
    );
  }
  return (
    <div>
      <p>Recorded: {size.count.toLocaleString()} requests</p>
      <p className="text-sm text-gray-600">
        metadata ≈ {formatMB(size.metaBytes)} MB, text bodies ≈{" "}
        {formatMB(size.textBodyBytes)} MB
      </p>
    </div>
  );
}

function LimitInput({
  label,
  value,
  onChange,
  totalBytes,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  totalBytes: number | undefined;
}) {
  return (
    <label className="flex items-center gap-3">
      <span className="w-44">{label}</span>
      <input
        type="number"
        min={0}
        step={1}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-28 px-3 py-2 border border-gray-300 rounded-md shadow-sm"
      />
      <span className="text-sm text-gray-600">
        {coverage(value, totalBytes)}
      </span>
    </label>
  );
}

function isLimit(value: string) {
  return /^\d+$/.test(value) && Number(value) <= 1024 * 1024;
}

function formatMB(bytes: number) {
  return Math.ceil(bytes / MB).toLocaleString();
}

function coverage(value: string, totalBytes: number | undefined) {
  if (!isLimit(value)) {
    return <span className="text-red-500">Enter an integer (0 or more)</span>;
  }
  if (totalBytes === undefined) {
    return null;
  }
  const limitBytes = Number(value) * MB;
  if (limitBytes >= totalBytes) {
    return "all included";
  }
  return `about ${Math.floor((limitBytes / totalBytes) * 100)}% of the total`;
}
