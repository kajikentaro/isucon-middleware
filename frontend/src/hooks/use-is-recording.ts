import { IsRecording } from "@/types";
import {
  getIsRecordingURL,
  getStartRecordingURL,
  getStopRecordingURL,
} from "@/utils/get-url";
import { useEffect, useState } from "react";

type RecordingStatus = "recording" | "not-recording" | "checking-status";

export function useIsRecording() {
  const [recordingStatus, setRecordingStatus] =
    useState<RecordingStatus>("checking-status");

  const updateStatus = async () => {
    const res = await fetch(getIsRecordingURL());
    const json = (await res.json()) as IsRecording;
    setRecordingStatus(json.isRecording ? "recording" : "not-recording");
  };

  const startRecording = async () => {
    setRecordingStatus("checking-status");
    await fetch(getStartRecordingURL(), { method: "POST" });
    setTimeout(updateStatus, 1000);
  };

  const stopRecording = async () => {
    setRecordingStatus("checking-status");
    await fetch(getStopRecordingURL(), { method: "POST" });
    setTimeout(updateStatus, 1000);
  };

  useEffect(() => {
    updateStatus();
  }, []);

  return { recordingStatus, startRecording, stopRecording };
}
