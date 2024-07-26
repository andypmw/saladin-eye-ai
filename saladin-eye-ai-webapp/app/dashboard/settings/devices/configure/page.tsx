import Link from 'next/link';

export default function ConfigureDevice() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Configure Camera
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <form>
            <label className="form-control w-full max-w-xs">
              <div className="label">
                <span className="label-text">Serial Number</span>
              </div>
              <input type="text" value="SALADIN-EYE-0039" disabled={true} placeholder="SALADIN-EYE-0039" className="input input-bordered w-full max-w-xs" />
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Type</span>
              </div>
              <input type="text" value="ESP32-CAM-0V2640" disabled={true} placeholder="ESP32-CAM-OV2640" className="input input-bordered w-full max-w-xs" />
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Name</span>
              </div>
              <input type="text" value="Gerbang Depan" placeholder="Gerbang Depan" className="input input-bordered w-full max-w-xs" />
              <div className="label">
                <span className="label-text-alt">Example: Gerbang Depan</span>
              </div>
            </label>

            <div className="modal-action">
              <Link href="/dashboard/settings/devices/list" className="btn btn-ghost">Cancel</Link>
              <Link href="/dashboard/settings/devices/list" className="btn btn-primary px-6">Save</Link>
            </div>
          </form>
        </div>
      </div>

    </div>
  );
}