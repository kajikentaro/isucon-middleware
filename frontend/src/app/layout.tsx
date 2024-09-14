import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Isucon Middleware",
  description: "Record request and response. Then, reproduce same request",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-w-[960px]">{children}</body>
    </html>
  );
}
