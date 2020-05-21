import React, { Component, Fragment, useState } from 'react';
import { withRouter } from 'react-router-dom';

import { useStoreAction, useStoreActions } from "easy-peasy";

// Externals
import classNames from 'classnames';
import compose from 'recompose/compose';
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import {
  Badge,
  IconButton,
  Popover,
  Toolbar,
  Typography
} from '@material-ui/core';

// Material icons
import {
  Menu as MenuIcon,
  Close as CloseIcon,
  NotificationsOutlined as NotificationsIcon,
  Input as InputIcon
} from '@material-ui/icons';

// Shared services
import { getNotifications } from 'services/notification';

// Custom components
import { NotificationList } from './components';

// Component styles
import styles from './styles';
import { useEffect } from 'react';

const Topbar = (props) => {

  const {
    classes,
    className,
    title,
    isSidebarOpen,
    onToggleSidebar
  } = props;

  const rootClassName = classNames(classes.root, className);

  const [ signal, setSignal] = useState(false);
  const [ notifications, setNotifications] = useState([]);
  const [ notificationsLimit, setNotificationsLimit ] = useState(4);
  const [ notificationsCount, setNotificationsCount ] = useState(0);
  const [ notificationsEl, setNotificationsEl ] = useState(null);

  const logout = useStoreActions(actions => actions.auth.logout);


  const showNotifications = Boolean(notificationsEl);


  const loadNotifications = async () => {
    try {

      const { notifications, notificationsCount } = await getNotifications(
        notificationsLimit
      );

      if (signal) {
          setNotifications(notifications);
          setNotificationsCount(notificationsCount);
      }
    } catch (error) {
      return;
    }
  }

  const handleSignOut = () => {
    const { history } = props;
    logout();
    history.push('/prihlasenie');
  };

  const handleShowNotifications = event => {
    setNotificationsEl(event.currentTarget);
  };

  const handleCloseNotifications = () => {
    setNotificationsEl(null);
  };

  useEffect(() => {
    setSignal(false);
  });

  useEffect(() => {
    if(!signal) {
      setSignal(true);
      loadNotifications();
    }
  }, []);

  return (
    <Fragment>
      <div className={rootClassName}>
        <Toolbar className={classes.toolbar}>
          <IconButton
            className={classes.menuButton}
            onClick={onToggleSidebar}
            variant="text"
          >
            {isSidebarOpen ? <CloseIcon /> : <MenuIcon />}
          </IconButton>
          <Typography
            className={classes.title}
            variant="h4"
          >
            {title}
          </Typography>
          <IconButton
            className={classes.notificationsButton}
            onClick={handleShowNotifications}
          >
            <Badge
              badgeContent={notificationsCount}
              color="primary"
              variant="dot"
            >
              <NotificationsIcon />
            </Badge>
          </IconButton>
          <IconButton
            className={classes.signOutButton}
            onClick={handleSignOut}
          >
            <InputIcon />
          </IconButton>
        </Toolbar>
      </div>
      <Popover
        anchorEl={notificationsEl}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center'
        }}
        onClose={handleCloseNotifications}
        open={showNotifications}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center'
        }}
      >
        <NotificationList
          notifications={notifications}
          onSelect={handleCloseNotifications}
        />
      </Popover>
    </Fragment>
  );
}

Topbar.propTypes = {
  className: PropTypes.string,
  classes: PropTypes.object.isRequired,
  history: PropTypes.object.isRequired,
  isSidebarOpen: PropTypes.bool,
  onToggleSidebar: PropTypes.func,
  title: PropTypes.string
};

Topbar.defaultProps = {
  onToggleSidebar: () => {}
};

export default compose(
  withRouter,
  withStyles(styles)
)(Topbar);
