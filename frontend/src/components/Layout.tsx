import { Sidebar } from '@/components/Sidebar';
import { Outlet } from 'react-router-dom';

export default function Layout() {
  return (
    <div className="flex min-h-screen bg-background">
      <Sidebar className="hidden md:block" />
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
