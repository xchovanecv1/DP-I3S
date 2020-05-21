import React, { Component } from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';

// Views
import Dashboard from './views/Dashboard';
import ProductList from './views/ProductList';

import UserList from './views/UserList';

import RoomsList from './views/RoomList'
import Room from './views/Room'

import PropList from './views/PropertyList'
import Prop from './views/Property'

import Typography from './views/Typography';
import Icons from './views/Icons';
import Account from './views/Account';
import Settings from './views/Settings';
import SignUp from './views/SignUp';
import SignIn from './views/SignIn';
import UnderDevelopment from './views/UnderDevelopment';
import NotFound from './views/NotFound';

const Router = (props) => {
  const { isLogged } = props;
  return (
    <Switch>
      <Redirect
        exact
        from="/"
        to="/dashboard"
      />
      {isLogged && <Redirect
        from="/prihlasenie"
        to="/dashboard"
      />
      }
      <Route
        component={SignIn}
        exact
        path="/prihlasenie"
      />
      // Iba pre autentifikovanych uzivatelov
      {!isLogged && <Redirect
        exact
        to="/prihlasenie"
      />
      }
      <Route
        component={Dashboard}
        exact
        path="/dashboard"
      />
      <Route
        component={UserList}
        exact
        path="/users"
      />
      <Route
        component={RoomsList}
        exact
        path="/miestnosti"
      />
      <Route
        component={Room}
        exact
        path="/miestnosti/:id"
      />
      <Route
        component={PropList}
        exact
        path="/zariadenia"
      />
      <Route
        component={Prop}
        exact
        path="/zariadenia/:id"
      />
      <Route
        component={ProductList}
        exact
        path="/products"
      />
      <Route
        component={Typography}
        exact
        path="/typography"
      />
      <Route
        component={Icons}
        exact
        path="/icons"
      />
      <Route
        component={Account}
        exact
        path="/account"
      />
      <Route
        component={Settings}
        exact
        path="/settings"
      />
      <Route
        component={SignUp}
        exact
        path="/sign-up"
      />
      <Route
        component={UnderDevelopment}
        exact
        path="/under-development"
      />
      <Route
        component={NotFound}
        exact
        path="/not-found"
      />
    </Switch>
  );
}

export default Router;
