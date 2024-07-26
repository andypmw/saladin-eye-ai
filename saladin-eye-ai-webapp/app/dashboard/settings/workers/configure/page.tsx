import Link from 'next/link';

export default function ConfigureWorker() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Configure Worker
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <form>
            <label className="form-control w-full max-w-xs">
              <div className="label">
                <span className="label-text">Server Public Address</span>
              </div>
              <input type="text" value="https://10.10.10.10:8888" disabled={true} placeholder="SALADIN-EYE-0039" className="input input-bordered w-full max-w-xs" />
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">New Worker Ticket ID</span>
              </div>
              <input type="text" value="84962048203493842945939432" disabled={true} placeholder="Gerbang Depan" className="input input-bordered w-full max-w-xs" />
              <div className="label">
                <span className="label-text-alt">Ramdomly generated</span>
              </div>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Worker Name</span>
              </div>
              <input type="text" value="on_kubenertes_worker" placeholder="on_kubenertes_worker" className="input input-bordered w-full max-w-xs" />
            </label>

            <div className="modal-action">
              <Link href="/dashboard/settings/workers/list" className="btn btn-ghost">Cancel</Link>
              <Link href="/dashboard/settings/workers/list" className="btn btn-primary px-6">Test</Link>
            </div>
          </form>
        </div>
      </div>

    </div>
  );
}