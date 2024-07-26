import Link from 'next/link';
import HomeIcon from '@heroicons/react/24/outline/HomeIcon';
import PhotoIcon from '@heroicons/react/24/outline/PhotoIcon';
import ShieldExclamationIcon from '@heroicons/react/24/outline/ShieldExclamationIcon';
import SettingIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';
import CameraIcon from '@heroicons/react/24/outline/CameraIcon';
import CircleStackIcon from '@heroicons/react/24/outline/CircleStackIcon';
import EyeIcon from '@heroicons/react/24/outline/EyeIcon';
import BriefcaseIcon from '@heroicons/react/24/outline/BriefcaseIcon';
import UsersIcon from '@heroicons/react/24/outline/UsersIcon';
import ArrowRightStartOnRectangleIcon from '@heroicons/react/24/outline/ArrowRightStartOnRectangleIcon';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-500">
      <div className="drawer lg:drawer-open">
        <input id="my-drawer-2" type="checkbox" className="drawer-toggle" />
        <div className="drawer-content flex flex-col">
          {/* Page content here */}
          {/* 
          <label htmlFor="my-drawer-2" className="btn btn-primary drawer-button lg:hidden">
            Open drawer
          </label>
          */}

          <div className="navbar sticky top-0 bg-base-100 z-10 shadow-md">
            <div className="flex-1">
              <h1 className="text-2xl font-semibold ml-2">SaladinEye.ai</h1>
            </div>
          </div>

          <main className="flex-1 overflow-y-auto md:pt-4 pt-4 pb-4 px-6 bg-base-200">
            {children}
          </main>
        </div>
        <div className="drawer-side z-30">
          <label htmlFor="my-drawer-2" aria-label="close sidebar" className="drawer-overlay"></label>

          <ul className="menu bg-base-100 text-base-content min-h-full w-80 p-4">
            <li>
              <Link href="/dashboard">
                <HomeIcon className="h-6 w-6" /> Home
              </Link>
            </li>
            <li>
              <Link href="/dashboard/recording">
                <PhotoIcon className="h-6 w-6" /> Recording
              </Link>
            </li>
            <li>
              <Link href="/dashboard/detection-alert">
                <ShieldExclamationIcon className="h-6 w-6" /> Detection & Alert
              </Link>
            </li>
            <li>
              <details open>
                <summary><SettingIcon className="h-6 w-6" /> Settings</summary>
                <ul>
                  <li>
                    <Link href="/dashboard/settings/devices/list">
                      <CameraIcon className="h-6 w-6" /> Camera
                    </Link>
                  </li>
                  <li>
                    <Link href="/dashboard/settings/storages/list">
                      <CircleStackIcon className="h-6 w-6" /> Storage
                    </Link>
                  </li>
                  <li>
                    <Link href="/dashboard/settings/machine-learning/list">
                      <EyeIcon className="h-6 w-6" /> Machine Learning
                    </Link>
                  </li>
                  <li>
                    <Link href="/dashboard/settings/workers/list">
                      <BriefcaseIcon className="h-6 w-6" /> Worker
                    </Link>
                  </li>
                  <li>
                    <Link href="/dashboard/settings/users/list">
                      <UsersIcon className="h-6 w-6" /> User
                    </Link>
                  </li>
                </ul>
              </details>
            </li>
            <li>
              <Link href="/">
                <ArrowRightStartOnRectangleIcon className="h-6 w-6" /> Exit
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </div>
  )
}