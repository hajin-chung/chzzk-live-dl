import React from "react";

interface LayoutProps {
  title?: string;
  children: React.ReactNode;
}

export function Layout({ title, children }: LayoutProps) {
  document.title = title ?? "simedia";
  return (
    <div className="flex flex-col items-center">
      <div className="w-full px-4 md:w-2xl">
        <div className="h-4" />
        {children}
      </div>
    </div>
  )
}
