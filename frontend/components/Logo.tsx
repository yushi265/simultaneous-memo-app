export function Logo({ className = "w-8 h-8" }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* メモ帳のアイコン */}
      <rect
        x="6"
        y="4"
        width="20"
        height="24"
        rx="2"
        stroke="currentColor"
        strokeWidth="2"
        fill="white"
      />
      
      {/* リアルタイム同期を表す円 */}
      <circle
        cx="22"
        cy="8"
        r="5"
        fill="#3B82F6"
        stroke="white"
        strokeWidth="1.5"
      />
      <path
        d="M22 6.5V8.5L23.5 10"
        stroke="white"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      
      {/* メモの線 */}
      <line x1="10" y1="12" x2="18" y2="12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="10" y1="16" x2="20" y2="16" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="10" y1="20" x2="16" y2="20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
    </svg>
  )
}