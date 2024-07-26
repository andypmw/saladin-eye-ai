import Image from "next/image";
import Link from 'next/link';

export default function Home() {
  return (
    <div className="flex flex-col md:flex-row h-screen">
      <div className="flex-1 flex items-center justify-center">
        <Image src="/_assets/saladin-eye-ai-esp32-banner-1.png" alt="SaladinEye.ai project title banner" width={600} height={600} />
      </div>

      <div className="flex-1 flex items-center justify-center bg-white">
        <form className="w-full max-w-sm pt-20 pr-8 pb-20 pl-8 ml-4 mr-4 bg-gray-100 rounded-2xl shadow-md bg-gradient-to-r from-black to-blue-800">
          <h2 className="text-2xl font-bold mb-6 text-center text-white">Login</h2>
          
          <label className="input input-bordered flex items-center gap-2 mt-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="h-4 w-4 opacity-70">
              <path
                d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
            </svg>
            <input type="text" className="grow" value="AndyPrimawan" placeholder="Username" />
          </label>
          <label className="input input-bordered flex items-center gap-2 mt-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              className="h-4 w-4 opacity-70">
              <path
                fillRule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clipRule="evenodd" />
            </svg>
            <input type="password" className="grow" value="password" />
          </label>

          <div className="flex items-center justify-between mt-8">
            <Link href="/dashboard" className="btn btn-secondary btn-block">Sign In</Link>
          </div>
        </form>
      </div>

      <div className="flex-1 flex items-center justify-center">
        <Image src="/_assets/saladin-eye-ai-esp32-banner-2.png" alt="SaladinEye.ai project tech stacks" width={600} height={600} />
      </div>
    </div>
  );
}
