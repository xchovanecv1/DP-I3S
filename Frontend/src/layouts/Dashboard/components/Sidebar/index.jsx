import React, { Component } from 'react';
import { Link, NavLink } from 'react-router-dom';
import { getInitials } from 'helpers';


import { useStoreState } from 'easy-peasy';
// Externals
import classNames from 'classnames';
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

import { connect } from 'react-redux';

// Material components
import {
  Avatar,
  CircularProgress,
  Divider,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListSubheader,
  Typography
} from '@material-ui/core';

// Material icons
import {
  DashboardOutlined as DashboardIcon,
  PeopleOutlined as PeopleIcon,
  ShoppingBasketOutlined as ShoppingBasketIcon,
  LockOpenOutlined as LockOpenIcon,
  TextFields as TextFieldsIcon,
  ImageOutlined as ImageIcon,
  InfoOutlined as InfoIcon,
  AccountBoxOutlined as AccountBoxIcon,
  SettingsOutlined as SettingsIcon,
  DevicesOutlined,
  MeetingRoomOutlined,
} from '@material-ui/icons';

// Component styles
import styles from './styles';

const Sidebar = (props) => {

  const { classes, className } = props;

  const token = useStoreState(state => state.auth.token);

  const rootClassName = classNames(classes.root, className);

  return (
    <nav className={rootClassName}>
        <div className={classes.logoWrapper}>
          <Link
            className={classes.logoLink}
            to="/"
          >
            <img
              alt="Brainalytica logo"
              className={classes.logoImage}
              src="/images/logos/android-icon-192x192.png"
            />
          </Link>
        </div>
        <Divider className={classes.logoDivider} />
        {token === null ? (
         <div className={classes.profile}><CircularProgress className={classes.progress} /></div>
        ) : (
          <div className={classes.profile}>
            <Link to="/account">
              {/*<Avatar
                alt={token.sub}
                className={classes.avatar}
                src="/images/avatars/avatar_1.png"
              />*/}
              <Avatar
                className={classes.avatar}
              >
                {getInitials(token.sub)}
              </Avatar>
            </Link>
            <Typography
              className={classes.nameText}
              variant="h6"
            >
              {token.sub}
            </Typography>
            <Typography
              className={classes.bioText}
              variant="caption"
            >
              {token.roles}
            </Typography>
          </div>) 
        }
        <Divider className={classes.profileDivider} />
        <List
          component="div"
          disablePadding
        >
          <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/dashboard"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Domov"
            />
          </ListItem>
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/users"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <PeopleIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Users"
            />
          </ListItem>}
          <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/miestnosti"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <MeetingRoomOutlined />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Miestnosti"
            />
          </ListItem>

          <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/zariadenia"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <DevicesOutlined />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Zariadenia"
            />
          </ListItem>
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/products"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <ShoppingBasketIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Products"
            />
          </ListItem>}
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/sign-in"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <LockOpenIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Authentication"
            />
          </ListItem>}
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/typography"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <TextFieldsIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Typography"
            />
          </ListItem>}
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/icons"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <ImageIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Icons and Images"
            />
          </ListItem>}
          <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/account"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <AccountBoxIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Účet"
            />
          </ListItem>
          {false && <ListItem
            activeClassName={classes.activeListItem}
            className={classes.listItem}
            component={NavLink}
            to="/settings"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <SettingsIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Settings"
            />
          </ListItem>}
        </List>
        <Divider className={classes.listDivider} />
        {false && <List
          component="div"
          disablePadding
          subheader={
            <ListSubheader className={classes.listSubheader}>
              Support
            </ListSubheader>
          }
        >
          <ListItem
            className={classes.listItem}
            component="a"
            href="https://devias.io/contact-us"
            target="_blank"
          >
            <ListItemIcon className={classes.listItemIcon}>
              <InfoIcon />
            </ListItemIcon>
            <ListItemText
              classes={{ primary: classes.listItemText }}
              primary="Customer support"
            />
          </ListItem>
        </List>}
      </nav>
  );
}

export default withStyles(styles)(Sidebar);
