import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from "jwt-decode";

const withAuth = (WrappedComponent) => {
  const WrapperComponent = (props) => {
    const navigate = useNavigate();

    useEffect(() => {
      const isAuthenticated = checkAuthenticationStatus();

      if (!isAuthenticated) {
        navigate('/login');
      }
    }, []);

    return <WrappedComponent {...props} />;
  };

  return WrapperComponent;
};

export default withAuth;

const checkAuthenticationStatus = () => {
  const token = localStorage.getItem('token');

  if (token) {
    try {
      const decodedToken = jwtDecode(token);
      const currentTime = Date.now() / 1000;

      if (decodedToken.exp < currentTime) {
        return false;
      }

      return true;
    } catch (error) {
      return false;
    }
  }

  return false;
};
