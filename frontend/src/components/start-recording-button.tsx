"use client";
import { useIsRecording } from "@/hooks/use-is-recording";

export default function StartRecordingButton() {
  const { recordingStatus, startRecording, stopRecording } = useIsRecording();
  const className =
    "border-2 py-1 px-2 rounded-full flex items-center w-36 justify-center";

  if (recordingStatus === "checking-status") {
    return (
      <button className={className} disabled>
        Now Loading ...
      </button>
    );
  }

  if (recordingStatus === "recording") {
    return (
      <button
        className={className + " border-red-500 text-red-500"}
        onClick={stopRecording}
      >
        Stop Recording
      </button>
    );
  }
  if (recordingStatus === "not-recording") {
    return (
      <button
        className={className + " border-blue-500 text-blue-500 font-bold"}
        onClick={startRecording}
      >
        Start Recording
      </button>
    );
  }
}
