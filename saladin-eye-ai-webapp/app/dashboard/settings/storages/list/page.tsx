import Link from 'next/link';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';

export default function ListStorage() {
  const dummyStorages = [
    ['Local Disk', 'hdd_lokal', 'RAID HDD 2 x 1 TiB', 'Active'],
    ['Amazon S3', 's3_jakarta', 'standard tier', 'Active'],
    ['Google Cloud Storage', 'gcs_jakarta', 'standard tier', 'Active'],
  ];

  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Storages
        <div className="inline-block float-right">
          <Link href="/dashboard/settings/storages/add" className="btn px-6 btn-sm normal-case btn-primary">
            Add New
          </Link>
        </div>
      </div>

      <div className="divider mt-2"></div>

      <div className="h-full w-full pb-6 bg-base-100">

        <div className="overflow-x-auto w-full">
          <table className="table w-full">
            <thead>
              <tr>
                <th>Type</th>
                <th>Name</th>
                <th>Note</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {dummyStorages.map((storage, index) => (
                <tr key={index}>
                  <td>{storage[0]}</td>
                  <td>{storage[1]}</td>
                  <td>{storage[2]}</td>
                  <td>
                    <div className="badge badge-primary">{storage[3]}</div>
                  </td>
                  <td>
                    <Link href="/dashboard/settings/storages/configure" className="btn btn-square btn-ghost">
                      <SettingIcon className="h-6 w-6" />
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

    </div>
  );
}