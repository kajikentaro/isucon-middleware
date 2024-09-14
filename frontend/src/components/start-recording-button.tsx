"use client";
import { useIsRecording } from "@/hooks/use-is-recording";
import { FaPlay, FaStop } from "react-icons/fa";

export default function StartRecordingButton() {
  const { recordingStatus, startRecording, stopRecording } = useIsRecording();
  const className =
    "border-2 p-2 rounded-full flex items-center w-44 justify-center gap-1";

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
        <FaStop />
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
        <FaPlay />
        Start Recording
      </button>
    );
  }
}
