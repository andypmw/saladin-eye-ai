import Link from 'next/link';

export default function ConfigureMachineLearning() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Configure Machine Learning
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <form>
            <label className="form-control w-full max-w-xs">
              <div className="label">
                <span className="label-text">Type</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>Human Detection</option>
                <option>Car Detection</option>
                <option>Motorcycle Detection</option>
              </select>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Model</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>COCO_based</option>
                <option>ResNet_based</option>
                <option>MiniNet_based</option>
              </select>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Deployment</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>ESP32-S3 onboard</option>
                <option>Centralized on Server</option>
              </select>
            </label>

            <div className="modal-action">
              <Link href="/dashboard/settings/machine-learning/list" className="btn btn-ghost">Cancel</Link>
              <Link href="/dashboard/settings/machine-learning/list" className="btn btn-primary px-6">Test</Link>
            </div>
          </form>
        </div>
      </div>

    </div>
  );
}