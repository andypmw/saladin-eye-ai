import Link from 'next/link';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';

export default function ListDevice() {
  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Cameras
        <div className="inline-block float-right">
          <Link href="/dashboard/settings/devices/add" className="btn px-6 btn-sm normal-case btn-primary">
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
                <th>Serial Number</th>
                <th>Name</th>
                <th>Type</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: 10 }).map((_, index) => (
                <tr key={index}>
                  <td>SALADIN-0001</td>
                  <td>Gerbang Depan</td>
                  <td>ESP32-CAM-OV2640</td>
                  <td>
                    <div className="badge badge-primary">Online</div>
                  </td>
                  <td>
                    <Link href="/dashboard/settings/devices/configure" className="btn btn-square btn-ghost">
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