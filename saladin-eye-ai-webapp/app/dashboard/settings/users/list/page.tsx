import Link from 'next/link';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';

export default function ListUser() {
  const dummyStorages = [
    ['1', 'andyprimawan', 'Admin', '1,424', 'Active'],
    ['2', 'satpam', 'Staff', '341', 'Active'],
    ['3', 'g4s', 'Staff', '1,201', 'Active'],
  ];

  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Users
        <div className="inline-block float-right">
          <Link href="/dashboard/settings/users/add" className="btn px-6 btn-sm normal-case btn-primary">
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
                <th>Username</th>
                <th>Role</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {dummyStorages.map((user, index) => (
                <tr key={index}>
                  <td>{user[0]}</td>
                  <td>{user[1]}</td>
                  <td>{user[2]}</td>
                  <td>{user[3]}</td>
                  <td>
                    <div className="badge badge-primary">{user[4]}</div>
                  </td>
                  <td>
                    <Link href="/dashboard/settings/users/configure" className="btn btn-square btn-ghost">
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