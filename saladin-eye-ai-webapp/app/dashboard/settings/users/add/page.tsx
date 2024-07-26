import Link from 'next/link';

export default function AddUser() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Add New User
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <form>
            <label className="form-control w-full max-w-xs">
              <div className="label">
                <span className="label-text">Username</span>
              </div>
              <input type="text" value="andyprimawan" className="input input-bordered w-full max-w-xs" />
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Password</span>
              </div>
              <input type="password" value="topsecret" className="input input-bordered w-full max-w-xs" />
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Role</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>Admin</option>
                <option>Staff</option>
              </select>
            </label>

            <label className="form-control w-full max-w-xs mt-4">
              <div className="label">
                <span className="label-text">Status</span>
              </div>
              <select className="select select-bordered">
                <option disabled selected>Pick one</option>
                <option>Active</option>
                <option>Inactive</option>
              </select>
            </label>

            <div className="modal-action">
              <Link href="/dashboard/settings/users/list" className="btn btn-ghost">Cancel</Link>
              <Link href="/dashboard/settings/users/list" className="btn btn-primary px-6">Save</Link>
            </div>
          </form>
        </div>
      </div>

    </div>
  );
}