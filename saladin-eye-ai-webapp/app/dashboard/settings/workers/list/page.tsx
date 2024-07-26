import Link from 'next/link';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';

export default function ListWorker() {
  const dummyStorages = [
    ['98239854932457', 'digitalocean_kube', '3 days 12 hours', '1,424', 'Active'],
    ['34853894939859', 'tencent_kube', '5 hours', '341', 'Active'],
    ['68384384384384', 'alibaba_kube', '1 day 8 hours', '1,201', 'Active'],
  ];

  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Workers
        <div className="inline-block float-right">
          <Link href="/dashboard/settings/workers/add" className="btn px-6 btn-sm normal-case btn-primary">
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
                <th>ID</th>
                <th>Name</th>
                <th>Uptime</th>
                <th>Completed Tasks</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {dummyStorages.map((worker, index) => (
                <tr key={index}>
                  <td>{worker[0]}</td>
                  <td>{worker[1]}</td>
                  <td>{worker[2]}</td>
                  <td>{worker[3]}</td>
                  <td>
                    <div className="badge badge-primary">{worker[4]}</div>
                  </td>
                  <td>
                    <Link href="/dashboard/settings/workers/configure" className="btn btn-square btn-ghost">
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