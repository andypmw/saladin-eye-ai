import Image from "next/image";

export default function Dashboard() {
  return (
    <div className="flex flex-col md:flex-row h-screen">
      <div className="flex-1 flex items-center justify-center">
        <Image src="/_assets/saladin-eye-ai-esp32-banner-1.png" alt="SaladinEye.ai project title banner" width={600} height={600} />
      </div>

      <div className="flex-1 flex items-center justify-center">
        <Image src="/_assets/saladin-eye-ai-esp32-banner-2.png" alt="SaladinEye.ai project tech stacks" width={600} height={600} />
      </div>
    </div>
  )
}