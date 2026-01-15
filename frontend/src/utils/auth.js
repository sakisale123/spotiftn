export const getUserRole = () => {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(
            window
                .atob(base64)
                .split('')
                .map(function (c) {
                    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
                })
                .join('')
        );

        const payload = JSON.parse(jsonPayload);
        return payload.role;
    } catch (error) {
        console.error('Failed to decode token:', error);
        return null;
    }
};

export const isAdmin = () => {
    return getUserRole() === 'admin';
};

export const getUserId = () => {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(
            window
                .atob(base64)
                .split('')
                .map(function (c) {
                    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
                })
                .join('')
        );

        const payload = JSON.parse(jsonPayload);
        return payload.userId;
    } catch (error) {
        console.error('Failed to decode token:', error);
        return null;
    }
};
