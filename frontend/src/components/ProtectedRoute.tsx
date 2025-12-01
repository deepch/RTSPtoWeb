import { useAuth } from '@/context/AuthContext';
import { Navigate, Outlet } from 'react-router-dom';

export default function ProtectedRoute({ role }: { role?: 'admin' | 'user' }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (role && user.role !== role && user.role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
}
