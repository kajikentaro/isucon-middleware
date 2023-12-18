export default function ThreeDotsAnimation() {
  return (
    <div className="flex items-center justify-center space-x-1">
      <div className="w-1 h-1 bg-black rounded-full animate-bounce"></div>
      <div className="w-1 h-1 bg-black rounded-full animate-bounce"></div>
      <div className="w-1 h-1 bg-black rounded-full animate-bounce"></div>
    </div>
  );
}
