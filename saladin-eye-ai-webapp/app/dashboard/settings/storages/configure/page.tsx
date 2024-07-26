import Link from 'next/link';

export default function ConfigureStorage() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Configure Storage
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
                <option>Local Disk</option>
                <option>Amazon S3</option>
                <option>Google Cloud Storage</option>
                <option>CloudFlare R2</option>
                <option>Minio</option>
              </select>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Region</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>ap-southeast-1</option>
                <option>ap-southeast-2</option>
                <option>ap-southeast-3</option>
                <option>us-east-1</option>
              </select>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Credential Method</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>Assume IAM Role</option>
                <option>Amazon Access Key ID and Secret</option>
              </select>
            </label>

            <div className="modal-action">
              <Link href="/dashboard/settings/storages/list" className="btn btn-ghost">Cancel</Link>
              <Link href="/dashboard/settings/storages/list" className="btn btn-primary px-6">Test</Link>
            </div>
          </form>
        </div>
      </div>

    </div>
  );
}