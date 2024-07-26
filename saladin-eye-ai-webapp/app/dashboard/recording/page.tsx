import Link from 'next/link';
import Image from "next/image";

export default function RecordingPage() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Recording
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 p-8">
            {Array.from({ length: 12 }, (_, index) => (
              <div className="card bg-base-100 shadow-xl">
                <figure>
                  <Image src="/_assets/saladin-eye-sample-esp32-photo.jpg" alt="SaladinEye.ai ESP32 sample photo" width={1600} height={1200} />
                </figure>
                <div className="card-body">
                  <h2 className="card-title">SALADIN-0001</h2>
                  <p>Gerbang Depan</p>
                  <div className="card-actions justify-end">
                    <button className="btn btn-primary">View</button>
                  </div>
                </div>
              </div>

            ))}
          </div>
        </div>
      </div>

    </div>
  );
}