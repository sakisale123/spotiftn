import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { isAdmin } from '../../utils/auth';

const AdminRoute = () => {
    if (!isAdmin()) {
        return <Navigate to="/artists" replace />;
    }

    return <Outlet />;
};

export default AdminRoute;
