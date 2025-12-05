import { Link, useLocation } from 'react-router-dom';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { LayoutDashboard, LogOut, LayoutGrid } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';

interface SidebarProps extends React.HTMLAttributes<HTMLDivElement> {}

export function Sidebar({ className, ...props }: SidebarProps) {
  const { pathname } = useLocation();
  const { logout } = useAuth();

  const links = [
    { href: '/', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/multiview', label: 'MultiView', icon: LayoutGrid },
  ];

  return (
    <div className={cn("pb-12 w-64 border-r min-h-screen bg-background", className)} {...props}>
      <div className="space-y-4 py-4">
        <div className="px-3 py-2">
          <h2 className="mb-2 px-4 text-lg font-semibold tracking-tight">
            RTSPtoWeb
          </h2>
          <div className="space-y-1">
            {links.map((link) => (
              <Button
                key={link.href}
                variant={pathname === link.href ? "secondary" : "ghost"}
                className="w-full justify-start"
                asChild
              >
                <Link to={link.href}>
                  <link.icon className="mr-2 h-4 w-4" />
                  {link.label}
                </Link>
              </Button>
            ))}
          </div>
        </div>
      </div>
      <div className="px-3 py-2 mt-auto">
          <Button variant="ghost" className="w-full justify-start text-red-500 hover:text-red-600 hover:bg-red-100 dark:hover:bg-red-900/20" onClick={logout}>
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
      </div>
    </div>
  );
}
