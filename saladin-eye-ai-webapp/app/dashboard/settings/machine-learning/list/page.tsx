import Link from 'next/link';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';

export default function ListMachineLearning() {
  const dummyMachineLearning = [
    ['Car Detection', 'COCO_based', 'ESP32-S3 on board', 'Active'],
    ['Human Detection', 'ResNet_based', 'Centralized on Server', 'Active'],
  ];

  return (
    <div className="card w-full p-6 bg-base-100 shadow-xl mt-2">

      <div className="text-xl font-semibold inline-block">
        Machine Learning
        <div className="inline-block float-right">
          <Link href="/dashboard/settings/machine-learning/add" className="btn px-6 btn-sm normal-case btn-primary">
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
                <th>Model</th>
                <th>Note</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {dummyMachineLearning.map((machineLearning, index) => (
                <tr key={index}>
                  <td>{machineLearning[0]}</td>
                  <td>{machineLearning[1]}</td>
                  <td>{machineLearning[2]}</td>
                  <td>
                    <div className="badge badge-primary">{machineLearning[3]}</div>
                  </td>
                  <td>
                    <Link href="/dashboard/settings/machine-learning/configure" className="btn btn-square btn-ghost">
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